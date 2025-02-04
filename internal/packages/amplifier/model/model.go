package model

import (
	"fmt"
	"time"

	"github.com/uptrace/bun"
)

/*
	pub enum Vote {
	    SucceededOnChain, // the txn was included on chain, and achieved the intended result
	    FailedOnChain,    // the txn was included on chain, but failed to achieve the intended result
	    NotFound,         // the txn could not be found on chain in any blocks at the time of voting
	}
*/
type AmplifierVoteStatus int64

var (
	PollStart   AmplifierVoteStatus = 1
	Yes         AmplifierVoteStatus = 2
	No          AmplifierVoteStatus = 3
	Unsubmitted AmplifierVoteStatus = 4
)

type AmplifierVote struct {
	bun.BaseModel     `bun:"table:amplifier"`
	ID                int64               `bun:"id,pk,autoincrement"`
	ChainInfoID       int64               `bun:"chain_info_id,pk,notnull"`
	Height            int64               `bun:"height,notnull"`
	VerifierAddressID int64               `bun:"verifier_address_id,notnull"`
	Status            AmplifierVoteStatus `bun:"status,notnull"`
	Timestamp         time.Time           `bun:"timestamp,notnull"`
	PollID            int64
}

func (apv AmplifierVote) String() string {
	return fmt.Sprintf("AmplifierVote<%d %d %d %d %d %d>",
		apv.ID,
		apv.ChainInfoID,
		apv.Height,
		apv.VerifierAddressID,
		apv.Status,
		apv.Timestamp.Unix(),
	)
}
