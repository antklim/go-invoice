package dynamo_test

import (
	"sort"
	"testing"

	"github.com/antklim/go-invoice/invoice"
	"github.com/antklim/go-invoice/storage/dynamo"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type byItemID []invoice.Item

func (x byItemID) Len() int           { return len(x) }
func (x byItemID) Less(i, j int) bool { return x[i].ID < x[j].ID }
func (x byItemID) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

type byDItemID []dynamo.Item

func (x byDItemID) Len() int           { return len(x) }
func (x byDItemID) Less(i, j int) bool { return x[i].ID < x[j].ID }
func (x byDItemID) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

func testInvoiceItem(t *testing.T, idx int, dItem dynamo.Item, item invoice.Item) {
	if dItem.ID != item.ID {
		t.Errorf("invalid dItem[%d].ID %q, want %q", idx, dItem.ID, item.ID)
	}

	if dItem.ProductName != item.ProductName {
		t.Errorf("invalid dItem[%d].ProductName %q, want %q", idx, dItem.ProductName, item.ProductName)
	}

	if dItem.Price != item.Price {
		t.Errorf("invalid dItem[%d].Price %d, want %d", idx, dItem.Price, item.Price)
	}

	if dItem.Qty != item.Qty {
		t.Errorf("invalid dItem[%d].Qty %d, want %d", idx, dItem.Qty, item.Qty)
	}

	if !dItem.CreatedAt.Equal(item.CreatedAt) {
		t.Errorf("invalid dItem[%d].CreatedAt %v, want %v", idx, dItem.CreatedAt, item.CreatedAt)
	}
}

func testInvoiceItems(t *testing.T, dItems []dynamo.Item, items []invoice.Item) {
	if len(dItems) != len(items) {
		t.Errorf("invalid dInvoice.items %v, want %v", dItems, items)
	}

	sort.Sort(byItemID(items))
	sort.Sort(byDItemID(dItems))

	for i := range dItems {
		ditem := dItems[i]
		item := items[i]
		testInvoiceItem(t, i, ditem, item)
	}
}

func testPutItemInput(t *testing.T, inv invoice.Invoice, input *dynamodb.PutItemInput) {
	if got, want := aws.StringValue(input.TableName), "invoices"; got != want {
		t.Errorf("invalid PutItemInput table %q, want %q", got, want)
	}

	var dinv dynamo.Invoice
	if err := dynamodbattribute.UnmarshalMap(input.Item, &dinv); err != nil {
		t.Fatalf("dynamodbattribute.UnmarshalMap() failed: %v", err)
	}

	if dinv.PK == "" {
		t.Error("dInvoice.PK should not be empty")
	}

	if dinv.ID != inv.ID {
		t.Errorf("invalid dInvoice.ID %q, want %q", dinv.ID, inv.ID)
	}

	if dinv.CustomerName != inv.CustomerName {
		t.Errorf("invalid dInvoice.CustomerName %q, want %q", dinv.CustomerName, inv.CustomerName)
	}

	if (dinv.Date == nil && inv.Date != nil) ||
		(dinv.Date != nil && inv.Date == nil) ||
		(dinv.Date != nil && inv.Date != nil && !dinv.Date.Equal(*inv.Date)) {
		t.Errorf("invalid dInvoice.Date %v, want %v", dinv.Date, inv.Date)
	}

	if invoice.Status(dinv.Status) != inv.Status {
		t.Errorf("invalid dInvoice.Status %d, want %d", dinv.Status, inv.Status)
	}

	testInvoiceItems(t, dinv.Items, inv.Items)

	if !dinv.CreatedAt.Equal(inv.CreatedAt) {
		t.Errorf("invalid dInvoice.CreatedAt %v, want %v", dinv.CreatedAt, inv.CreatedAt)
	}

	if !dinv.UpdatedAt.Equal(inv.UpdatedAt) {
		t.Errorf("invalid dInvoice.UpdatedAt %v, want %v", dinv.UpdatedAt, inv.UpdatedAt)
	}
}

func testAddItemConditionExression(t *testing.T, id string, input *dynamodb.PutItemInput) {
	if got, want := aws.StringValue(input.ConditionExpression), "#0 <> :0"; got != want {
		t.Errorf("PutItem condition expression %q, want %q", got, want)
	}

	testPutItemExpressionAttribute(t, "0", "id", id, input)
}

func testUpdateItemConditionExression(t *testing.T, id string, input *dynamodb.PutItemInput) {
	if got, want := aws.StringValue(input.ConditionExpression), "#0 = :0"; got != want {
		t.Errorf("PutItem condition expression %q, want %q", got, want)
	}
	testPutItemExpressionAttribute(t, "0", "id", id, input)
}

// testPutItemExpressionAttribute tests that expression attribute with the index
// idx mapped to the expected field name and value val.
func testPutItemExpressionAttribute(t *testing.T, idx, name, val string, input *dynamodb.PutItemInput) {
	if got := aws.StringValue(input.ExpressionAttributeNames["#"+idx]); got != name {
		t.Errorf("PutItem condition expression: #%s attribute name %q, want %q", idx, got, name)
	}

	var actual string
	if err := dynamodbattribute.Unmarshal(input.ExpressionAttributeValues[":"+idx], &actual); err != nil {
		t.Fatalf("PutItem condition expression: unmarshal :%s attribute value failed: %v", idx, err)
	}
	if actual != val {
		t.Errorf("PutItem condition expression: :%s attribute value %q, want %q", idx, actual, val)
	}
}

func testGetItemInput(t *testing.T, id string, input *dynamodb.GetItemInput) {
	if got, want := aws.StringValue(input.TableName), "invoices"; got != want {
		t.Errorf("invalid GetItemInput table %q, want %q", got, want)
	}

	var pk struct{ PK string }
	if err := dynamodbattribute.UnmarshalMap(input.Key, &pk); err != nil {
		t.Fatalf("GetItemInput key unmarshal failed: %v", err)
	}

	if want := "INVOICE#" + id; pk.PK != want {
		t.Errorf("GetItemInput key value %q, want %q", pk.PK, want)
	}
}
