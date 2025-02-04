package indexer

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/cosmostation/cvms/internal/common/api"
	"github.com/cosmostation/cvms/internal/testutil"
	"github.com/stretchr/testify/assert"
)

// voteExist := false
// voteCounter := 0
// for _, event := range summary.TxsEvents {
// 	if event.TypeName == pollVoteType {
// 		voteExist = true
// 		voteCounter++
// 	}
// }

// t.Logf("total vote: %d", voteCounter)

// vmList := make([]VoteMessage, 0)
// requester := exporter.APIClient.R()
// ch := make(chan helper.Result)
// var wg sync.WaitGroup
// wg.Add(voteCounter)

// if voteExist {
// 	for _, tx := range summary.TxHashes {
// 		queryPath := fmt.Sprintf("/cosmos/tx/v1beta1/txs/%s", tx)
// 		go func(ch chan helper.Result) {
// 			defer wg.Done()

// 			resp, err := requester.Get(queryPath)
// 			if err != nil || resp == nil {
// 				t.Log("Unexpeced status code")
// 				ch <- helper.Result{Item: nil, Success: false}
// 				return
// 			}
// 			if resp.StatusCode() != http.StatusOK {
// 				t.Log("Unexpeced status code")
// 				ch <- helper.Result{Item: nil, Success: false}
// 				return
// 			}

// 			vmList, err := sdkhelper.TxParser(resp.Body(), func() *VoteMessage {
// 				return &VoteMessage{}
// 			})

// 			if err != nil {
// 				t.Error(err)
// 				ch <- helper.Result{Item: nil, Success: false}
// 				return
// 			}

// 			for _, vm := range vmList {
// 				ch <- helper.Result{Item: vm, Success: true}
// 			}
// 		}(ch)
// 		time.Sleep(10 * time.Millisecond)
// 	}
// }

// go func() {
// 	wg.Wait()
// 	close(ch)
// }()

// errorCount := 0
// for r := range ch {
// 	if r.Success {
// 		vmList = append(vmList, r.Item.(VoteMessage))
// 		continue
// 	}
// 	errorCount++
// }

// if errorCount > 0 {
// 	t.Errorf("current errors count: %d", errorCount)
// }

// t.Logf("got total vm: %d", len(vmList))

// 만약 wasm 보트가 있으면, DB 업데이트
// for _, vm := range vmList {
// 	t.Logf("%+v", vm)
// }

// 근데 존재하지 않는 poll이면? poll쿼리
// 아래 쿼리를 이용해서 poll를 채워넣는다.

// end poll이 뜨면?
// 위와 같은 방법을 재구성?
// 존재하지 않는 폴이면? 스킵?
// 존해자는 폴이면 업데이트?

func TestLogic(t *testing.T) {
	_ = testutil.SetupForTest()
	exporter := testutil.GetTestExporter()
	exporter.SetAPIEndPoint(os.Getenv("TEST_AXELAR_API_ENDPOINT"))
	exporter.SetRPCEndPoint(os.Getenv("TEST_AXELAR_RPC_ENDPOINT"))

	height := int64(17293863)
	txsEvents, _, err := api.GetBlockResults(exporter.CommonClient, height)
	assert.NoError(t, err)

	pollExist := false

	// modify poll to polls
	polls := AmplifierPollStartFillter(txsEvents)

	if len(polls) > 0 {
		pollExist = true

		// call crawling poll voting function
	}

	for idx, poll := range polls {
		verifiers := poll.Participants
		pollStartHeight := height
		expireHeightStr := poll.ExpiresAt
		pollID := poll.PollID
		pollSourceChain := poll.PollID

		t.Logf("\nPoll Details:\n Verifiers: %d\n Start Height: %d\n Expire Height: %s\n Poll ID: %s\n Poll Source Chain: %s",
			len(verifiers[0]), pollStartHeight, expireHeightStr, pollID, pollSourceChain)

		expireHeight, err := strconv.ParseInt(expireHeightStr, 10, 64)
		assert.NoError(t, err)

		// start to crawl txs by height to expire height
		if pollExist {
			for h := (pollStartHeight + 1); h <= expireHeight; h++ {
				log.Printf("crawling height to  check voting: %d", h)

				// txsEvents, _, err := api.GetBlockTxs(exporter.CommonClient, height)
				// assert.NoError(t, err)
			}
		}

		if idx == 0 {
			break
		}
	}

	// 1. get current testnet contracts
	// by using -> https://github.com/axelarnetwork/axelar-contract-deployments/blob/main/axelar-chains-config/info/testnet.json#L2699

	// 2. when got a contract wasm
	// "body": {
	// "messages": [
	// 	{
	// 	"@type": "/cosmwasm.wasm.v1.MsgExecuteContract",
	// 	"sender": "axelar1r25hycaye0uz3k554mdu4a7dvc82uelj7y6ddn",
	// 	"contract": "axelar1sykyha8kzf35kc5hplqk76kdufntjn6w45ntwlevwxp74dqr3rvsq7fazh",

	// "message_id": "5AoX7DUKR1dQSAyxXv5kjvpLmXjLW2qqJ25E35mk6WLh-0",
	// "new_verifier_set": {

	// 3. make a new voting set with default unvoted status

	// 4. got some blocks and then update the status by txs

	// 5. when I find end_poll

}

func TestAmplifierPollQuery(t *testing.T) {
	_ = testutil.SetupForTest()
	exporter := testutil.GetTestExporter()
	exporter.SetAPIEndPoint(os.Getenv("TEST_AXELAR_API_ENDPOINT"))

	/*
		CONTRACT=axelar1ce9rcvw8htpwukc048z9kqmyk5zz52d5a7zqn9xlq2pg0mxul9mqxlx2cq
		QUERY='{"poll":{"poll_id":"853"}}'
		axelard query wasm contract-state smart $CONTRACT $QUERY
	*/
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	requester := exporter.APIClient.R().SetContext(ctx)
	contractAddress := "axelar1ce9rcvw8htpwukc048z9kqmyk5zz52d5a7zqn9xlq2pg0mxul9mqxlx2cq"
	pollID := "853"
	queryMethod := fmt.Sprintf(`{"poll":{"poll_id":"%s"}}`, pollID)
	base64EncodedMethod := base64.StdEncoding.EncodeToString([]byte(queryMethod))
	queryPath := fmt.Sprintf(
		"/cosmwasm/wasm/v1/contract/%s/smart/%s",
		contractAddress, base64EncodedMethod,
	)

	resp, err := requester.Get(queryPath)
	if err != nil || resp.StatusCode() != http.StatusOK {
		t.Logf("unexpeced err: [%d] %s", resp.StatusCode(), err)
	}

	// ref; https://github.com/axelarnetwork/axelar-contract-deployments/blob/main/axelar-chains-config/info/testnet.json#L2699
	// start height = expire height - blockExpiry(10)
	t.Logf("%s", resp)
}

// func TestMakePollDataList(t *testing.T) {
// 	// Get current working directory
// 	pwd, _ := os.Getwd()
// 	dirPath := fmt.Sprintf("%s/tests", pwd)

// 	// Read the directory
// 	files, err := os.ReadDir(dirPath)
// 	if err != nil {
// 		t.Fatalf("Failed to read directory: %s", err)
// 	}

// 	summaryList := make([]types.TxSearchSummary, 0)
// 	// Loop through the files
// 	for _, file := range files {
// 		// Skip directories
// 		if file.IsDir() {
// 			continue
// 		}
// 		t.Logf("Reading file: %s ...", file.Name())

// 		// Construct the full file path
// 		filePath := filepath.Join(dirPath, file.Name())

// 		// Read the file content as bytes
// 		content, err := os.ReadFile(filePath)
// 		if err != nil {
// 			t.Logf("Error reading file %s: %v\n", file.Name(), err)
// 			continue
// 		}

// 		txHashes, txEvents, err := parser.CosmosTxSearchParser(content)
// 		if err != nil {
// 			t.Logf("Error parsing file %s: %v\n", file.Name(), err)
// 			continue
// 		}

// 		summaryList = append(summaryList, types.TxSearchSummary{
// 			TxsEvents: txEvents,
// 			TxHashes:  txHashes,
// 		})
// 	}

// 	// amplifierEventFilter(summaryList)
// }

// // 17283786 poll start?
// // 17283786 poll voting
// func TestTxParser(t *testing.T) {
// 	txJSON := `{
//         "tx": {
//             "body": {
//                 "messages": [
//                     {
//                         "@type": "/cosmwasm.wasm.v1.MsgExecuteContract",
//                         "sender": "axelar1...",
//                         "contract": "contract123",
//                         "msg": {
//                             "vote": {
//                                 "poll_id": "123",
//                                 "votes": ["yes"]
//                             }
//                         }
//                     }
//                 ]
//             }
//         }
//     }`

// 	txBz := []byte(txJSON)

// 	// Pass the parser with a factory function for VoteMessage
// 	results, err := sdkhelper.TxParser(txBz, func() *VoteMessage {
// 		return &VoteMessage{}
// 	})

// 	if err != nil {
// 		t.Fatalf("Failed to parse transaction: %v", err)
// 	}

// 	if len(results) != 1 {
// 		t.Fatalf("Expected 1 message, got %d", len(results))
// 	}

// 	t.Logf("Parsed results: %+v", results[0])
// }

// // Ensures VoteMessage implements MessageParser
// var _ sdkhelper.MessageParser = &VoteMessage{}

// type VoteMessage struct {
// 	Sender   string
// 	Contract string
// 	PollID   string
// 	Status   string
// }

// // processMessage processes a single message and extracts details, returning a generic type T
// func (vm *VoteMessage) ProcessMessage(msg interface{}) error {
// 	const (
// 		AuxiliaryBatchType      = "/axelar.auxiliary.v1beta1.BatchRequest"
// 		WasmExecuteContractType = "/cosmwasm.wasm.v1.MsgExecuteContract"
// 	)

// 	messageMap, ok := msg.(map[string]interface{})
// 	if !ok {
// 		return errors.New("failed to cast message to map")
// 	}

// 	// Extract @type
// 	msgType, ok := messageMap["@type"].(string)
// 	if !ok {
// 		return errors.New("failed to extract '@type' from message")
// 	}

// 	switch msgType {
// 	case AuxiliaryBatchType:
// 		// Recursively process nested messages
// 		nestedMessages, ok := messageMap["messages"].([]interface{})
// 		if !ok {
// 			return errors.New("failed to extract nested 'messages'")
// 		}
// 		return vm.ProcessMessage(nestedMessages)

// 	case WasmExecuteContractType:
// 		// Parse the current message
// 		return vm.Parse(messageMap)

// 	default:
// 		return fmt.Errorf("unexpected message type: %s", msgType)
// 	}
// }

// // Parse implements the MessageParser interface for VoteMessage
// func (vm *VoteMessage) Parse(messageMap map[string]interface{}) error {
// 	sender, ok := messageMap["sender"].(string)
// 	if !ok {
// 		return errors.New("failed to extract 'sender'")
// 	}

// 	contract, ok := messageMap["contract"].(string)
// 	if !ok {
// 		return errors.New("failed to extract 'contract'")
// 	}

// 	voteStatus, pollID, err := extractVote(messageMap)
// 	if err != nil {
// 		return fmt.Errorf("failed to extract vote details: %w", err)
// 	}

// 	// Assign parsed values to the struct fields
// 	vm.Sender = sender
// 	vm.Contract = contract
// 	vm.PollID = pollID
// 	vm.Status = voteStatus
// 	return nil
// }

// // Helper function to extract vote status and poll ID
// func extractVote(messageMap map[string]interface{}) (string, string, error) {
// 	// Extract "msg" map
// 	msg, ok := messageMap["msg"].(map[string]interface{})
// 	if !ok {
// 		return "", "", errors.New("failed to extract 'msg' from message")
// 	}

// 	// Extract "vote" map
// 	vote, ok := msg["vote"].(map[string]interface{})
// 	if !ok {
// 		return "", "", errors.New("failed to extract 'vote' from 'msg'")
// 	}

// 	// Extract poll ID
// 	pollID, ok := vote["poll_id"].(string)
// 	if !ok {
// 		return "", "", errors.New("failed to extract 'poll_id' from 'vote'")
// 	}

// 	// Extract "votes" array and return the first vote status
// 	votes, ok := vote["votes"].([]interface{})
// 	if !ok || len(votes) == 0 {
// 		return "", "", errors.New("failed to extract 'votes' or 'votes' is empty")
// 	}

// 	// Return vote status and poll ID
// 	voteStatus, ok := votes[0].(string)
// 	if !ok {
// 		return "", "", errors.New("failed to extract vote status from 'votes'")
// 	}

// 	return voteStatus, pollID, nil
// }
