package mocks

import (
	"sync"

	"github.com/antklim/go-invoice/storage/dynamo"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type dynamoOp int

const (
	putItem dynamoOp = iota
	query
)

var dynamoOps = map[string]dynamoOp{
	"PutItem": putItem,
	"Query":   query,
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

func (api *DynamoAPI) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	api.Lock()
	defer api.Unlock()
	api.recordPutItemCall(input)

	return nil, api.errors[putItem]
}

func (api *DynamoAPI) Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	api.Lock()
	defer api.Unlock()
	api.recordQueryCall(input)

	return nil, api.errors[query]
}

func (api *DynamoAPI) CalledTimes(op string) int {
	api.RLock()
	defer api.RUnlock()

	times, ok := api.callsTimes[dynamoOpFrom(op)]
	if !ok {
		return -1
	}
	return times
}

func (api *DynamoAPI) recordPutItemCall(input *dynamodb.PutItemInput) {
	api.callsTimes[putItem]++
	api.callsArgs[putItem] = append(api.callsArgs[putItem], input)
}

func (api *DynamoAPI) recordQueryCall(input *dynamodb.QueryInput) {
	api.callsTimes[query]++
	api.callsArgs[query] = append(api.callsArgs[query], input)
}
