package dynamo

import (
	"errors"

	"github.com/antklim/go-invoice/invoice"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type API interface {
	PutItem(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	Query(*dynamodb.QueryInput) (*dynamodb.QueryOutput, error)
}

type Dynamo struct {
	client API
}

var _ invoice.Storage = (*Dynamo)(nil)

func New(client API) *Dynamo {
	return &Dynamo{client: client}
}

func (d *Dynamo) AddInvoice(inv invoice.Invoice) error {
	return errors.New("not implemented")
}

func (d *Dynamo) FindInvoice(id string) (*invoice.Invoice, error) {
	return nil, errors.New("not implemented")
}

func (d *Dynamo) UpdateInvoice(inv invoice.Invoice) error {
	return errors.New("not implemented")
}
