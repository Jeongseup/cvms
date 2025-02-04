package parser

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/cosmostation/cvms/internal/common/types"
	"github.com/pkg/errors"
)

const AxelarBatchRequestMessageType = "/axelar.auxiliary.v1beta1.BatchRequest"

func ExtractAmplifierVotes(resp []byte) (
	/* block height */ int64,
	/* block timestamp */ time.Time,
	types.AmplifierVote,
	error,
) {
	result := types.CosmosBlockTxsResponse{}
	err := json.Unmarshal(resp, &result)
	if err != nil {
		return 0, time.Time{}, types.AmplifierVote{}, err
	}

	for _, tx := range result.Txs {
		for _, message := range tx.Body.Messages {
			var preResult map[string]json.RawMessage
			if err := json.Unmarshal(message, &preResult); err != nil {
				return 0, time.Time{}, types.AmplifierVote{}, err
			}

			if rawType, ok := preResult["@type"]; ok {
				var typeValue string
				if err := json.Unmarshal(rawType, &typeValue); err != nil {
					return 0, time.Time{}, types.AmplifierVote{}, err
				}

				votes, err := parseDynamicMessage(message, typeValue)
				if err != nil {
					return 0, time.Time{}, types.AmplifierVote{}, err
				}

				blockHeight, err := strconv.ParseInt(result.Block.Header.Height, 10, 64)
				if err != nil {
					return 0, time.Time{}, types.AmplifierVote{}, err
				}

				return blockHeight, result.Block.Header.Time, votes, nil
			}
		}
	}

	return 0, time.Time{}, types.AmplifierVote{}, errors.New("unexpected errors")
}

// parseDynamicMessage dynamically parses the message based on its type.
func parseDynamicMessage(message json.RawMessage, typeURL string) (types.AmplifierVote, error) {
	switch typeURL {
	case AxelarBatchRequestMessageType:
		var msg types.MsgAxelarBatchRequest
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Failed to parse MsgInjectedCheckpoint: %v", err)
			return types.AmplifierVote{}, err
		}

		for _, m := range msg.Messages {
			if len(m.Message.Vote.Votes) > 0 {
				vote := types.AmplifierVote{
					VerifierAddress: m.Sender,
					Contract:        m.Contract,
					PollID:          m.Message.Vote.PollID,
					Status:          switchVoteStatus(m.Message.Vote.Votes[0]),
				}

				return vote, nil
			}
		}
		return types.AmplifierVote{}, nil
	default:
		return types.AmplifierVote{}, fmt.Errorf("unknown message type: %s", typeURL)
	}
}

const (
	AmplifierSuccessVoteStr = "succeeded_on_chain"
)

func switchVoteStatus(voteStr string) types.AmplifierVoteStatus {
	switch voteStr {
	case AmplifierSuccessVoteStr:
		return types.Success
	default:
		return types.Failed
	}
}
