package types

// MsgInjectedCheckpoint is the structure for the specific message type.
type MsgAxelarBatchRequest struct {
	TypeURL  string
	Sender   string
	Messages []AmplifierVoteMessage `json:"messages"`
}

type AmplifierVoteMessage struct {
	TypeURL  string
	Sender   string
	Contract string
	Message  struct {
		Vote struct {
			PollID string
			Votes  []string
		}
	} `json:"msg"`
	Funds interface{} `json:"-"`
}

type AmplifierVoteStatus int

const (
	Failed  AmplifierVoteStatus = iota // failed
	Success                            // succeeded_on_chain
)

type AmplifierVote struct {
	VerifierAddress string
	Contract        string
	PollID          string
	Status          AmplifierVoteStatus
}

// messages": [
// {
// "@type":
// "sender": "axelar1j3u6kd4027wln9vnvmg449hmc3xj2m2g5uh69q",
// "messages": [
// {
// "@type": "/cosmwasm.wasm.v1.MsgExecuteContract",
// "sender": "axelar1j3u6kd4027wln9vnvmg449hmc3xj2m2g5uh69q",
// "contract": "axelar1sykyha8kzf35kc5hplqk76kdufntjn6w45ntwlevwxp74dqr3rvsq7fazh",
// "msg": {
// "vote": {
// "poll_id": "25",
// "votes": [
// "succeeded_on_chain"
// ]
// }
// },
// "funds": []
// }
// ]
// }
// ]
