package storage

import (
	"github.com/antklim/go-invoice/invoice"
	"github.com/antklim/go-invoice/storage/dynamo"
	"github.com/antklim/go-invoice/storage/memory"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Memory struct{}

func (Memory) MakeStorage() invoice.Storage {
	return memory.New()
}

var _ invoice.StorageFactory = new(Memory)

type Dynamo struct {
	client dynamo.API
	table  string
}

func NewDynamo(table string, opts ...DynamoOption) *Dynamo {
	dopts := defaultDynamoOptions
	for _, o := range opts {
		o.apply(&dopts)
	}

	cfg := &aws.Config{Region: aws.String(dopts.region)}
	if dopts.endpoint != "" {
		cfg.WithEndpoint(dopts.endpoint)
	}
	sess := session.Must(session.NewSession(cfg))
	client := dynamodb.New(sess)

	return &Dynamo{
		client: client,
		table:  table,
	}
}

func (s *Dynamo) MakeStorage() invoice.Storage {
	return dynamo.New(s.client, s.table)
}

var _ invoice.StorageFactory = (*Dynamo)(nil)

type dynamoOptions struct {
	endpoint string
	region   string
}

var defaultDynamoOptions = dynamoOptions{
	region: "ap-southeast-2",
}

type DynamoOption interface {
	apply(*dynamoOptions)
}

type funcDynamoOption struct {
	f func(*dynamoOptions)
}

func (f *funcDynamoOption) apply(o *dynamoOptions) {
	f.f(o)
}

func newFuncDynamoOption(f func(*dynamoOptions)) DynamoOption {
	return &funcDynamoOption{f: f}
}

func WithEndpoint(v string) DynamoOption {
	return newFuncDynamoOption(func(o *dynamoOptions) {
		o.endpoint = v
	})
}

func WithRegion(v string) DynamoOption {
	return newFuncDynamoOption(func(o *dynamoOptions) {
		o.region = v
	})
}
