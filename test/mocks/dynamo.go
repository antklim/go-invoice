package mocks

import (
	"errors"

	"github.com/antklim/go-invoice/storage/dynamo"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DynamoAPI struct {
}

var _ dynamo.API = (*DynamoAPI)(nil)

func (api *DynamoAPI) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return nil, errors.New("not implemented")
}

func (api *DynamoAPI) Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	return nil, errors.New("not implemented")
}
