package dynamo_test

import (
	"testing"

	"github.com/antklim/go-invoice/invoice"
	"github.com/antklim/go-invoice/storage/dynamo"
	"github.com/antklim/go-invoice/test/mocks"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

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
		testPutItemConditionExression(t, inv, dinput)
	})
}

func TestFindInvoice(t *testing.T) {
	// TODO: implement
	t.Skip("not implemented")
}

func TestUpdateInvoice(t *testing.T) {
	// TODO: implement
	t.Skip("not implemented")
}
