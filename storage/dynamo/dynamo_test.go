package dynamo_test

import (
	"testing"

	"github.com/antklim/go-invoice/invoice"
	"github.com/antklim/go-invoice/storage/dynamo"
	"github.com/antklim/go-invoice/test/mocks"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func testInvoiceItems(t *testing.T, dItems []dynamo.Item, items []invoice.Item) {
	if len(dItems) != len(items) {
		t.Errorf("invalid dInvoice.items %v, want %v", dItems, items)
	}
}

func testPutItemInput(t *testing.T, inv invoice.Invoice, input *dynamodb.PutItemInput) {
	if got, want := aws.StringValue(input.TableName), "invoices"; got != want {
		t.Errorf("invalid PutItem input table %q, want %q", got, want)
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

	if dinv.Status != inv.Status.String() {
		t.Errorf("invalid dInvoice.Status %q, want %q", dinv.Status, inv.Status.String())
	}

	testInvoiceItems(t, dinv.Items, inv.Items)

	if !dinv.CreatedAt.Equal(inv.CreatedAt) {
		t.Errorf("invalid dInvoice.CreatedAt %v, want %v", dinv.CreatedAt, inv.CreatedAt)
	}

	if !dinv.UpdatedAt.Equal(inv.UpdatedAt) {
		t.Errorf("invalid dInvoice.UpdatedAt %v, want %v", dinv.UpdatedAt, inv.UpdatedAt)
	}
}

func TestInvoicePK(t *testing.T) {
	inv := invoice.NewInvoice("123ABC", "customer")
	got := dynamo.InvoicePartitionKey(inv)
	want := "INVOICE#123ABC"
	if got != want {
		t.Errorf("invalid invoice partition key %q, want %q", got, want)
	}
}

func TestAddInvoice(t *testing.T) {
	t.Run("called with correct input", func(t *testing.T) {
		client := mocks.NewDynamoAPI()
		strg := dynamo.New(client, "invoices")
		inv := invoice.NewInvoice("123", "customer")
		if err := inv.AddItem(invoice.NewItem("456", "pen", 1000, 3)); err != nil {
			t.Errorf("inv.AddItem() failed: %v", err)
		}

		if err := strg.AddInvoice(inv); err != nil {
			t.Errorf("AddInvoice(%v) failed: %v", inv, err)
		}

		if got, want := client.CalledTimes("PutItem"), 1; got != want {
			t.Errorf("client.PutItem() called %d times, want %d call(s)", got, want)
		}

		ncall := 1
		input := client.NthCall("PutItem", ncall)
		if input == nil {
			t.Fatalf("input of PutItem call #%d is nil", ncall)
		}

		dinput, ok := input.(*dynamodb.PutItemInput)
		if !ok {
			t.Errorf("type of PutItem input is %T, want *dynamodb.PutItemInput", input)
		}

		testPutItemInput(t, inv, dinput)
	})

	t.Run("fails when trying to add axisting invoice", func(t *testing.T) {
		// client := mocks.DynamoAPI{}
		// strg := dynamo.New(&client, "invoices")
		// inv := invoice.NewInvoice("123", "customer")
		// inv.AddItem(invoice.NewItem("456", "pen", 1000, 3))

		// if err := strg.AddInvoice(inv); err != nil {
		// 	t.Errorf("AddInvoice(%v) failed: %v", inv, err)
		// }
		// err := strg.AddInvoice(inv)
		// if err == nil {
		// 	t.Errorf("expected second call AddInvoice(%v) to fail", inv)
		// } else if got, want := err.Error(), `ID "123" exists`; got != want {
		// 	t.Errorf("second call AddInvoice(%v) = %v, want %v", inv, got, want)
		// }
	})
}

func TestFindInvoice(t *testing.T) {
	t.Skip("not implemented")
}

func TestUpdateInvoice(t *testing.T) {
	t.Skip("not implemented")
}
