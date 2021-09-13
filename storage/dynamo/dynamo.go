package dynamo

import (
	"errors"
	"fmt"
	"time"

	"github.com/antklim/go-invoice/invoice"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	dKeyDelim        = "#"
	dInvoicePKPrefix = "INVOICE"
)

type dInvoice struct {
	PK           string     `dynamodbav:"pk"`
	ID           string     `dynamodbav:"id"`
	CustomerName string     `dynamodbav:"customerName"`
	Date         *time.Time `dynamodbav:"issueDate"`
	Status       string     `dynamodbav:"status"`
	Items        []dItem    `dynamodbav:"items"`
	CreatedAt    time.Time  `dynamodbav:"createdAt"`
	UpdatedAt    time.Time  `dynamodbav:"updatedAt"`
}

func newDinvoice(inv invoice.Invoice) dInvoice {
	dItems := make([]dItem, 0, len(inv.Items))
	for _, invItem := range inv.Items {
		dItem := newDitem(invItem)
		dItems = append(dItems, dItem)
	}

	pk := dInvoicePartitionKey(inv)
	return dInvoice{
		PK:           pk,
		ID:           inv.ID,
		CustomerName: inv.CustomerName,
		Date:         inv.Date,
		Status:       inv.Status.String(),
		Items:        dItems,
		CreatedAt:    inv.CreatedAt,
		UpdatedAt:    inv.UpdatedAt,
	}
}

func dInvoicePartitionKey(inv invoice.Invoice) string {
	return fmt.Sprintf("%s%s%s", dInvoicePKPrefix, dKeyDelim, inv.ID)
}

type dItem struct {
	ID          string    `dynamodbav:"id"`
	ProductName string    `dynamodbav:"productName"`
	Price       uint      `dynamodbav:"price"`
	Qty         uint      `dynamodbav:"qty"`
	CreatedAt   time.Time `dynamodbav:"createdAt"`
}

func newDitem(item invoice.Item) dItem {
	return dItem{
		ID:          item.ID,
		ProductName: item.ProductName,
		Price:       item.Price,
		Qty:         item.Qty,
		CreatedAt:   item.CreatedAt,
	}
}

type API interface {
	PutItem(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	Query(*dynamodb.QueryInput) (*dynamodb.QueryOutput, error)
}

type Dynamo struct {
	client API
	table  string
}

var _ invoice.Storage = (*Dynamo)(nil)

func New(client API, table string) *Dynamo {
	return &Dynamo{
		client: client,
		table:  table,
	}
}

func (d *Dynamo) AddInvoice(inv invoice.Invoice) error {
	dinv := newDinvoice(inv)
	item, err := dynamodbattribute.MarshalMap(dinv)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(d.table),
		Item:      item,
	}
	_, err = d.client.PutItem(input)
	return err
}

func (d *Dynamo) FindInvoice(id string) (*invoice.Invoice, error) {
	return nil, errors.New("not implemented")
}

func (d *Dynamo) UpdateInvoice(inv invoice.Invoice) error {
	return errors.New("not implemented")
}
