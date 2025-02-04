package indexer

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"

	"github.com/cosmostation/cvms/internal/helper"
	"github.com/cosmostation/cvms/internal/helper/healthcheck"

	"github.com/cosmostation/cvms/internal/common"
	indexertypes "github.com/cosmostation/cvms/internal/common/indexer/types"
	"github.com/cosmostation/cvms/internal/packages/amplifier/repository"
)

var Subsystem = "amplifier_verifier"

type AmplifierIndexer struct {
	*common.Indexer
	repo repository.AmplifierIndexerRepository
}

// Compile-time Assertion
var _ common.IIndexer = (*AmplifierIndexer)(nil)

func NewAmplifierIndexer(p common.Packager) (*AmplifierIndexer, error) {
	status := helper.GetOnChainStatus(p.RPCs, p.ProtocolType)
	if status.ChainID == "" {
		return nil, errors.New("failed to create new veindexer")
	}
	indexer := common.NewIndexer(p, p.Package, status.ChainID)
	repo := repository.NewRepository(*p.IndexerDB, indexertypes.SQLQueryMaxDuration)
	indexer.Lh = indexertypes.LatestHeightCache{LatestHeight: status.BlockHeight}
	return &AmplifierIndexer{indexer, repo}, nil
}

func (indexer *AmplifierIndexer) Start() error {
	err := indexer.InitChainInfoID()
	if err != nil {
		return errors.Wrap(err, "failed to init chain_info_id")
	}

	alreadyInit, err := indexer.repo.CheckIndexpoinerAlreadyInitialized(repository.IndexName, indexer.ChainInfoID)
	if err != nil {
		return errors.Wrap(err, "failed to check init tables")
	}
	if !alreadyInit {
		indexer.Warnln("it's not initialized in the database, so that veindexer will init for this package")
		indexer.repo.InitPartitionTablesByChainInfoID(repository.IndexName, indexer.ChainID, indexer.Lh.LatestHeight)
	}

	// get last index pointer, index pointer is always initalize if not exist
	initIndexPointer, err := indexer.repo.GetLastIndexPointerByIndexTableName(repository.IndexName, indexer.ChainInfoID)
	if err != nil {
		return errors.Wrap(err, "failed to get last index pointer")
	}

	err = indexer.FetchValidatorInfoList()
	if err != nil {
		return errors.Wrap(err, "failed to fetch validator_info list")
	}

	indexer.Infof("loaded index pointer(last saved height): %d", initIndexPointer.Pointer)
	indexer.Infof("initial vim length: %d for %s chain", len(indexer.Vim), indexer.ChainID)

	// init indexer metrics
	// indexer.initLabelsAndMetrics()
	// go fetch new height in loop, it must be after init metrics
	go indexer.FetchLatestHeight()
	// loop
	go indexer.Loop(initIndexPointer.Pointer)
	// loop update recent miss counter metrics
	go func() {
		for {
			indexer.Infoln("update recent miss counter metrics and sleep 5s sec...")
			// indexer.updateRecentMissCounterMetric()
			time.Sleep(time.Second * 5)
		}
	}()
	// loop partion table time retention by env parameter
	go func() {
		for {
			indexer.Infof("for time retention, delete old records over %s and sleep %s", indexer.RetentionPeriod, indexertypes.RetentionQuerySleepDuration)
			indexer.repo.DeleteOldValidatorExtensionVoteList(indexer.ChainID, indexer.RetentionPeriod)
			time.Sleep(indexertypes.RetentionQuerySleepDuration)
		}
	}()
	return nil
}

func (veidx *AmplifierIndexer) Loop(indexPoint int64) {
	isUnhealth := false
	for {
		// node health check
		if isUnhealth {
			healthAPIs := healthcheck.FilterHealthEndpoints(veidx.APIs, veidx.ProtocolType)
			for _, api := range healthAPIs {
				veidx.SetAPIEndPoint(api)
				veidx.Warnf("API endpoint will be changed with health endpoint for this package: %s", api)
				isUnhealth = false
				break
			}

			healthRPCs := healthcheck.FilterHealthRPCEndpoints(veidx.RPCs, veidx.ProtocolType)
			for _, rpc := range healthRPCs {
				veidx.SetRPCEndPoint(rpc)
				veidx.Warnf("RPC endpoint will be changed with health endpoint for this package: %s", rpc)
				isUnhealth = false
				break
			}

			if len(healthAPIs) == 0 || len(healthRPCs) == 0 {
				isUnhealth = true
				veidx.Errorln("failed to get any health endpoints from healthcheck filter, retry sleep 10s")
				time.Sleep(indexertypes.UnHealthSleep)
				continue
			}
		}

		// set new index point height
		newIndexPointerHeight := indexPoint + 1

		// trying to sync with new index pointer height
		newIndexPointer, err := veidx.batchSync(indexPoint, newIndexPointerHeight)
		if err != nil {
			common.Health.With(veidx.RootLabels).Set(0)
			common.Ops.With(veidx.RootLabels).Inc()
			isUnhealth = true
			veidx.Errorf("failed to sync validators vote status in %d height: %s\nit will be retried after sleep %s...",
				indexPoint, err, indexertypes.AfterFailedRetryTimeout.String(),
			)
			time.Sleep(indexertypes.AfterFailedRetryTimeout)
			continue
		}

		// update index point
		indexPoint = newIndexPointer

		// update health and ops
		common.Health.With(veidx.RootLabels).Set(1)
		common.Ops.With(veidx.RootLabels).Inc()

		// logging & sleep
		if veidx.Lh.LatestHeight > indexPoint {
			// when node catching_up is true, sleep 100 milli sec
			veidx.WithField("catching_up", true).
				Infof("latest height is %d but updated index pointer is %d ... remaining %d blocks", veidx.Lh.LatestHeight, indexPoint, (veidx.Lh.LatestHeight - indexPoint))
			time.Sleep(indexertypes.CatchingUpSleepDuration)
		} else {
			// when node already catched up, sleep 5 sec
			veidx.WithField("catching_up", false).
				Infof("updated index pointer to %d and sleep %s sec...", indexPoint, indexertypes.DefaultSleepDuration.String())
			time.Sleep(indexertypes.DefaultSleepDuration)
		}
	}
}

// insert chain-info into chain_info table
func (indexer *AmplifierIndexer) InitChainInfoID() error {
	isNewChain := false
	var chainInfoID int64
	chainInfoID, err := indexer.repo.SelectChainInfoIDByChainID(indexer.ChainID)
	if err != nil {
		if err == sql.ErrNoRows {
			indexer.Infof("this is new chain id: %s", indexer.ChainID)
			isNewChain = true
		} else {
			return errors.Wrap(err, "failed to select chain_info_id by chain-id")
		}
	}

	if isNewChain {
		chainInfoID, err = indexer.repo.InsertChainInfo(indexer.ChainName, indexer.ChainID, indexer.Mainnet)
		if err != nil {
			return errors.Wrap(err, "failed to insert new chain_info_id by chain-id")
		}
	}

	indexer.ChainInfoID = chainInfoID
	return nil
}

// TODO: change FetchValidatorInfoList to FetchVerifier...
func (indexer *AmplifierIndexer) FetchValidatorInfoList() error {
	// get already saved validator-set list for mapping validators ids
	verifierInfoList, err := indexer.repo.GetVerifierInfoListByChainInfoID(indexer.ChainInfoID)
	if err != nil {
		return errors.Wrap(err, "failed to get validator info list")
	}

	// when the this pacakge starts, set validator-id map
	for _, verifier := range verifierInfoList {
		indexer.Vim[verifier.VerifierAddress] = int64(verifier.ID)
	}

	return nil
}
