package dynamo

import (
	"errors"
	"fmt"
	"time"

	"github.com/antklim/go-invoice/invoice"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

var errUnknownUnmarshalSource = errors.New("unknown unmarshal source")

const (
	dKeyDelim        = "#"
	dInvoicePKPrefix = "INVOICE"
)

type dInvoice struct {
	PK           string     `dynamodbav:"pk"`
	ID           string     `dynamodbav:"id"`
	CustomerName string     `dynamodbav:"customerName"`
	Date         *time.Time `dynamodbav:"issueDate"`
	Status       int        `dynamodbav:"status"`
	Items        []dItem    `dynamodbav:"items"`
	CreatedAt    time.Time  `dynamodbav:"createdAt"`
	UpdatedAt    time.Time  `dynamodbav:"updatedAt"`
}

func (dInv *dInvoice) InvoiceMarshal() invoice.Invoice {
	items := make([]invoice.Item, 0, len(dInv.Items))
	for _, dItem := range dInv.Items {
		item := dItem.InvoiceItemMarshal()
		items = append(items, item)
	}

	return invoice.Invoice{
		ID:           dInv.ID,
		CustomerName: dInv.CustomerName,
		Date:         dInv.Date,
		Status:       invoice.Status(dInv.Status),
		Items:        items,
		CreatedAt:    dInv.CreatedAt,
		UpdatedAt:    dInv.UpdatedAt,
	}
}

func invoiceUnmarshal(inv invoice.Invoice) *dInvoice {
	dItems := make([]dItem, 0, len(inv.Items))
	for _, invItem := range inv.Items {
		dItem := invoiceItemUnmarshal(invItem)
		dItems = append(dItems, dItem)
	}

	pk := dInvoicePartitionKey(inv.ID)
	return &dInvoice{
		PK:           pk,
		ID:           inv.ID,
		CustomerName: inv.CustomerName,
		Date:         inv.Date,
		Status:       int(inv.Status),
		Items:        dItems,
		CreatedAt:    inv.CreatedAt,
		UpdatedAt:    inv.UpdatedAt,
	}
}

func getItemOutputUnmarshal(output *dynamodb.GetItemOutput) (*dInvoice, error) {
	if output == nil || len(output.Item) == 0 {
		return nil, nil
	}

	var dInv dInvoice
	if err := dynamodbattribute.UnmarshalMap(output.Item, &dInv); err != nil {
		return nil, err
	}

	return &dInv, nil
}

// unmarshalDinvoice unmarshals value v to an instance of dInvoice.
func unmarshalDinvoice(v interface{}) (*dInvoice, error) {
	if inv, ok := v.(invoice.Invoice); ok {
		return invoiceUnmarshal(inv), nil
	}

	if output, ok := v.(*dynamodb.GetItemOutput); ok {
		return getItemOutputUnmarshal(output)
	}

	return nil, errUnknownUnmarshalSource
}

// dInvoicePartitionKey builds invoice partition key based on invoice id.
func dInvoicePartitionKey(id string) string {
	return fmt.Sprintf("%s%s%s", dInvoicePKPrefix, dKeyDelim, id)
}

type dItem struct {
	ID          string    `dynamodbav:"id"`
	ProductName string    `dynamodbav:"productName"`
	Price       int       `dynamodbav:"price"`
	Qty         int       `dynamodbav:"qty"`
	CreatedAt   time.Time `dynamodbav:"createdAt"`
}

func (di *dItem) InvoiceItemMarshal() invoice.Item {
	return invoice.Item{
		ID:          di.ID,
		ProductName: di.ProductName,
		Price:       di.Price,
		Qty:         di.Qty,
		CreatedAt:   di.CreatedAt,
	}
}

func invoiceItemUnmarshal(item invoice.Item) dItem {
	return dItem{
		ID:          item.ID,
		ProductName: item.ProductName,
		Price:       item.Price,
		Qty:         item.Qty,
		CreatedAt:   item.CreatedAt,
	}
}

type API interface {
	GetItem(*dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
	PutItem(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
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
	cond := expression.Name("id").NotEqual(expression.Value(inv.ID))
	expr, err := expression.NewBuilder().
		WithCondition(cond).
		Build()
	if err != nil {
		return err
	}

	err = d.upsertInvoice(inv, expr)
	if isConditionalCheckError(err) {
		return fmt.Errorf("invoice %q exists", inv.ID)
	}

	return err
}

func (d *Dynamo) FindInvoice(id string) (*invoice.Invoice, error) {
	primaryKey := map[string]string{"pk": dInvoicePartitionKey(id)}
	pk, err := dynamodbattribute.MarshalMap(primaryKey)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(d.table),
		Key:       pk,
	}

	result, err := d.client.GetItem(input)
	if err != nil {
		return nil, err
	}

	dInv, err := unmarshalDinvoice(result)
	if err != nil {
		return nil, err
	}
	if dInv == nil {
		return nil, nil
	}

	inv := dInv.InvoiceMarshal()
	return &inv, nil
}

func (d *Dynamo) UpdateInvoice(inv invoice.Invoice) error {
	cond := expression.Name("id").Equal(expression.Value(inv.ID))
	expr, err := expression.NewBuilder().
		WithCondition(cond).
		Build()
	if err != nil {
		return err
	}

	inv.UpdatedAt = time.Now()
	err = d.upsertInvoice(inv, expr)
	if isConditionalCheckError(err) {
		return fmt.Errorf("invoice %q not found", inv.ID)
	}

	return err
}

// upsertInvoice inserts or updates an invoice depending on provided expression.
func (d *Dynamo) upsertInvoice(inv invoice.Invoice, expr expression.Expression) error {
	dinv, err := unmarshalDinvoice(inv)
	if err != nil {
		return err
	}

	item, err := dynamodbattribute.MarshalMap(dinv)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName:                 aws.String(d.table),
		Item:                      item,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ConditionExpression:       expr.Condition(),
	}

	_, err = d.client.PutItem(input)
	return err
}

func isConditionalCheckError(err error) bool {
	var aerr awserr.Error
	return errors.As(err, &aerr) &&
		aerr.Code() == dynamodb.ErrCodeConditionalCheckFailedException
}
