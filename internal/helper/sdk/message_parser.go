package sdkhelper

import (
	"encoding/json"
	"errors"
	"fmt"
)

// vote tx parser
func TxParser[T any](txBz []byte, typeURL string, processFunc func(msg interface{}) (T, error)) ([]T, error) {
	// Parse the JSON into a generic map
	var response map[string]interface{}
	if err := json.Unmarshal(txBz, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// Extract the "tx" object
	tx, err := extractMap(response, "tx")
	if err != nil {
		return nil, fmt.Errorf("failed to extract 'tx': %w", err)
	}

	// Extract the "body" object
	body, err := extractMap(tx, "body")
	if err != nil {
		return nil, fmt.Errorf("failed to extract 'body': %w", err)
	}

	// Extract the "messages" array
	messages, err := extractArray(body, "messages")
	if err != nil {
		return nil, fmt.Errorf("failed to extract 'messages': %w", err)
	}

	// Traverse messages to find vote information
	resultList := make([]T, 0)
	for _, msg := range messages {
		result, err := processFunc(msg)
		if err != nil {
			fmt.Printf("Error processing message: %v\n", err)
		}
		resultList = append(resultList, result)
	}

	return resultList, nil
}

// Helper function to process a single message and extract details
// func processMessage(msg interface{}) (VoteMessage, error) {
// 	messageMap, ok := msg.(map[string]interface{})
// 	if !ok {
// 		return VoteMessage{}, errors.New("failed to cast message to map")
// 	}

// 	// Extract @type
// 	msgType, ok := messageMap["@type"].(string)
// 	if !ok {
// 		return VoteMessage{}, errors.New("failed to extract '@type' from message")
// 	}

// 	switch msgType {
// 	case AuxiliaryBatchType:
// 		// Recursively process nested messages
// 		return processNestedMessages(messageMap)
// 	case WasmExecuteContractType:
// 		// Parse the vote message
// 		return parseVoteMessage(messageMap)
// 	default:
// 		return VoteMessage{}, fmt.Errorf("unexpected message type: %s", msgType)
// 	}
// }

// Helper function to process nested messages for AuxiliaryBatchType
// func processNestedMessages(messageMap map[string]interface{}) (VoteMessage, error) {
// 	nestedMessages, ok := messageMap["messages"].([]interface{})
// 	if !ok {
// 		return VoteMessage{}, errors.New("failed to extract nested 'messages'")
// 	}

// 	// NOTE: current thought only one message in the nested messages
// 	for _, nestedMsg := range nestedMessages {
// 		vm, err := processMessage(nestedMsg)
// 		if err != nil {
// 			return VoteMessage{}, fmt.Errorf("error processing nested message: %w", err)
// 		}

// 		return vm, nil
// 	}

//		return VoteMessage{}, errors.New("failed to process nested message")
//	}
//
// processNestedMessages processes nested messages and returns a generic type T
func NestedMessagesProcess[T any](
	nestedMessages []interface{},
	processFunc func(interface{}) (T, error),
) (T, error) {
	for _, nestedMsg := range nestedMessages {
		result, err := processFunc(nestedMsg)
		if err != nil {
			return *new(T), fmt.Errorf("error processing nested message: %w", err)
		}
		return result, nil
	}
	return *new(T), errors.New("failed to process nested message")
}

// Helper function to extract a map from a key
func extractMap(data map[string]interface{}, key string) (map[string]interface{}, error) {
	value, ok := data[key]
	if !ok {
		return nil, fmt.Errorf("key '%s' not found", key)
	}
	mapValue, ok := value.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("key '%s' is not a map", key)
	}
	return mapValue, nil
}

// Helper function to extract an array from a key
func extractArray(data map[string]interface{}, key string) ([]interface{}, error) {
	value, ok := data[key]
	if !ok {
		return nil, fmt.Errorf("key '%s' not found", key)
	}
	arrayValue, ok := value.([]interface{})
	if !ok {
		return nil, fmt.Errorf("key '%s' is not an array", key)
	}
	return arrayValue, nil
}
