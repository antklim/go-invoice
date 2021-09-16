package mocks

import (
	"sync"

	"github.com/antklim/go-invoice/storage/dynamo"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type dynamoOp int

const (
	getItem dynamoOp = iota
	putItem
)

var dynamoOps = map[string]dynamoOp{
	"GetItem": getItem,
	"PutItem": putItem,
}

func dynamoOpFrom(op string) dynamoOp {
	dop, ok := dynamoOps[op]
	if !ok {
		return -1
	}
	return dop
}

type DynamoAPI struct {
	errors map[dynamoOp]error

	sync.RWMutex // guards calls
	callsTimes   map[dynamoOp]int
	callsArgs    map[dynamoOp][]interface{}
}

func NewDynamoAPI() *DynamoAPI {
	return &DynamoAPI{
		callsTimes: make(map[dynamoOp]int),
		callsArgs:  make(map[dynamoOp][]interface{}),
	}
}

var _ dynamo.API = (*DynamoAPI)(nil)

func (api *DynamoAPI) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	api.Lock()
	defer api.Unlock()
	api.recordGetItemCall(input)

	return nil, api.errors[getItem]
}

func (api *DynamoAPI) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	api.Lock()
	defer api.Unlock()
	api.recordPutItemCall(input)

	return nil, api.errors[putItem]
}

// CalledTimes returns amount of times the DynamoDB operation was called. It
// returns -1 when unknown operation provided.
func (api *DynamoAPI) CalledTimes(op string) int {
	api.RLock()
	times, ok := api.callsTimes[dynamoOpFrom(op)]
	api.RUnlock()

	if !ok {
		return -1
	}
	return times
}

// NthCall returns input of nth operation to DynamoDB. Counter n starts from 1.
// It returns nil in case of unknown operation or when n is greater than amount
// of calls recorded.
func (api *DynamoAPI) NthCall(op string, n int) interface{} {
	api.RLock()
	calls, ok := api.callsArgs[dynamoOpFrom(op)]
	api.RUnlock()

	if !ok {
		return nil
	}

	if n <= 0 || n > len(calls) {
		return nil
	}

	return calls[n-1]
}

func (api *DynamoAPI) recordPutItemCall(input *dynamodb.PutItemInput) {
	api.callsTimes[putItem]++
	api.callsArgs[putItem] = append(api.callsArgs[putItem], input)
}

func (api *DynamoAPI) recordGetItemCall(input *dynamodb.GetItemInput) {
	api.callsTimes[getItem]++
	api.callsArgs[getItem] = append(api.callsArgs[getItem], input)
}
