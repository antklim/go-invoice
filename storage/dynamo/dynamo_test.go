package dynamo_test

import (
	"testing"

	"github.com/antklim/go-invoice/invoice"
	"github.com/antklim/go-invoice/storage/dynamo"
	"github.com/antklim/go-invoice/test/mocks"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func TestInvoicePK(t *testing.T) {
	got := dynamo.InvoicePartitionKey("123ABC")
	want := "INVOICE#123ABC"
	if got != want {
		t.Errorf("invalid invoice partition key %q, want %q", got, want)
	}
}

func TestInvoiceMarshalUnmarshal(t *testing.T) {
	// TODO: test invoice -> dInvoice -> invoice
}

func TestAddInvoice(t *testing.T) {
	t.Run("called with correct input", func(t *testing.T) {
		client := mocks.NewDynamoAPI()
		strg := dynamo.New(client, "invoices")
		inv := invoice.NewInvoice("123", "customer")

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
		testAddItemConditionExression(t, inv.ID, dinput)
	})
}

func TestFindInvoice(t *testing.T) {
	client := mocks.NewDynamoAPI()
	strg := dynamo.New(client, "invoices")
	invID := "123"

	if _, err := strg.FindInvoice(invID); err != nil {
		t.Errorf("FindInvoice(%q) failed: %v", invID, err)
	}

	if got, want := client.CalledTimes("GetItem"), 1; got != want {
		t.Errorf("client.GetItem() called %d times, want %d call(s)", got, want)
	}

	ncall := 1
	input := client.NthCall("GetItem", ncall)
	if input == nil {
		t.Fatalf("input of GetItem call #%d is nil", ncall)
	}

	dinput, ok := input.(*dynamodb.GetItemInput)
	if !ok {
		t.Errorf("type of GetItem input is %T, want *dynamodb.GetItemInput", input)
	}

	testGetItemInput(t, invID, dinput)
}

func TestUpdateInvoice(t *testing.T) {
	client := mocks.NewDynamoAPI()
	strg := dynamo.New(client, "invoices")
	inv := invoice.NewInvoice("123", "customer")
	if err := inv.AddItem(invoice.NewItem("456", "pen", 1000, 3)); err != nil {
		t.Errorf("inv.AddItem() failed: %v", err)
	}

	if err := strg.UpdateInvoice(inv); err != nil {
		t.Errorf("UpdateInvoice(%v) failed: %v", inv, err)
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
	testUpdateItemConditionExression(t, inv.ID, dinput)
}
