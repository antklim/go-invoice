package dynamo_test

import (
	"testing"

	"github.com/antklim/go-invoice/invoice"
	"github.com/antklim/go-invoice/storage/dynamo"
	"github.com/antklim/go-invoice/test/mocks"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func testPutItemInput(t *testing.T, inv invoice.Invoice, input *dynamodb.PutItemInput) {
	if got, want := aws.StringValue(input.TableName), "invoices"; got != want {
		t.Errorf("invalid PutItem input table %q, want %q", got, want)
	}

	// TODO: unmarshal dynamo attributes to dInvoice
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
