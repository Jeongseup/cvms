package indexer

import (
	"fmt"

	"github.com/cosmostation/cvms/internal/common/types"
	"github.com/cosmostation/cvms/internal/helper"
)

type Poll struct {
	ContractAddress      string   `json:"_contract_address"`
	ConfirmationHeight   string   `json:"confirmation_height"`
	ExpiresAt            string   `json:"expires_at"`
	Messsage             string   `json:"message"`
	Participants         []string `json:"participants"`
	PollID               string   `json:"poll_id"`
	SourceChain          string   `json:"source_chain"`
	SourceGatewayAddress string   `json:"source_gateway_address"`
}

type PollVote struct {
	ContractAddress string `json:"_contract_address"`
	PollID          string `json:"poll_id"`
	Voter           string `json:"voter"`
}

type EndPoll struct {
	ContractAddress string   `json:"_contract_address"`
	PollID          string   `json:"poll_id"`
	SourceChain     string   `json:"source_chain"`
	Results         []string `json:"results"`
}

const pollStartType = "wasm-verifier_set_poll_started"
const pollVoteType = "wasm-voted"
const pollEndedType = "wasm-poll_ended"

// db update logic for axelar amplifier
// 1. extract poll & verifier set
// 2. insert the poll into db
// 3. if any votes found in the results, update the status on the db
// 4. skip expire logic, currently each verifier doesn't vote again by themselves
// 4. extract end poll, update poll status finished?
func AmplifierPollStartFillter(events []types.Event) []Poll {
	polls := make([]Poll, 0)
	for _, event := range events {
		switch event.TypeName {
		case pollStartType:
			poll := new(Poll)
			for _, attr := range event.Attributes {
				helper.SetFieldByTag(poll, attr.Key, attr.Value)
			}
			polls = append(polls, *poll)
		}
	}
	return polls
}

func amplifierEventFilter2(events []types.Event) {
	for _, event := range events {
		switch event.TypeName {
		case pollStartType:
			poll := new(Poll)
			for _, attr := range event.Attributes {
				helper.SetFieldByTag(poll, attr.Key, attr.Value)
			}

			fmt.Printf("%v", poll)
		case pollVoteType:
			pollVote := new(PollVote)
			for _, attr := range event.Attributes {
				helper.SetFieldByTag(pollVote, attr.Key, attr.Value)
			}

		case pollEndedType:
			endPoll := new(EndPoll)
			for _, attr := range event.Attributes {
				helper.SetFieldByTag(endPoll, attr.Key, attr.Value)
			}
		}
	}

}
