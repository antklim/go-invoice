package dynamo_test

import (
	"encoding/json"
	"errors"
	"os"
	"path"
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
	t.Run("dInvoice - invoice unmarshal/marshal", func(t *testing.T) {
		inv := invoice.NewInvoice("123", "customer")
		if err := inv.AddItem(invoice.NewItem("456", "pen", 1000, 3)); err != nil {
			t.Errorf("inv.AddItem() failed: %v", err)
		}

		dInv, err := dynamo.UnmarshalDinvoice(inv)
		if err != nil {
			t.Errorf("UnmarshalDinvoice(%v) failed: %v", inv, err)
		}

		if got := dInv.InvoiceMarshal(); !inv.Equal(&got) {
			t.Errorf("invalid invoice %v, want %v", got, inv)
		}
	})

	t.Run("dInvoice - get item output unmarshal", func(t *testing.T) {
		{
			output := (*dynamodb.GetItemOutput)(nil)
			dInv, err := dynamo.UnmarshalDinvoice(output)
			if err != nil {
				t.Errorf("UnmarshalDinvoice(%v) failed: %v", output, err)
			}
			if dInv != nil {
				t.Errorf("UnmarshalDinvoice(%v): %v, want nil", output, dInv)
			}
		}

		{
			testDataDir := "../../test/fixtures"
			testCases := []struct {
				desc         string
				file         string
				id           string
				hasIssueDate bool
				nItems       int
				status       int
			}{
				{
					desc:   "open invoice",
					file:   "get-item-open-invoice.json",
					id:     "170bf55e-ca81-4a17-99ad-54f6411d610b",
					nItems: 0,
					status: int(invoice.Open),
				},
				{
					desc:         "issued invoice",
					file:         "get-item-issued-invoice.json",
					id:           "170bf55e-ca81-4a17-99ad-54f6411d610c",
					hasIssueDate: true,
					nItems:       1,
					status:       int(invoice.Issued),
				},
			}
			for _, tC := range testCases {
				t.Run(tC.desc, func(t *testing.T) {
					testData, err := os.Open(path.Join(testDataDir, tC.file))
					if err != nil {
						t.Fatalf("failed to open test data: %v", err)
					}
					defer testData.Close()

					var output dynamodb.GetItemOutput
					if err := json.NewDecoder(testData).Decode(&output); err != nil {
						t.Fatalf("failed to Decode test data: %v", err)
					}

					dInv, err := dynamo.UnmarshalDinvoice(&output)
					if err != nil {
						t.Errorf("UnmarshalDinvoice(%v) failed: %v", output, err)
					}

					if dInv.ID != tC.id {
						t.Errorf("invalid dInv.ID %q, want %q", dInv.ID, tC.id)
					}

					if dInv.Status != tC.status {
						t.Errorf("invalid dInv.Status %d, want %d", dInv.Status, tC.status)
					}

					if tC.hasIssueDate && dInv.Date == nil {
						t.Error("dInv.Date should be defined")
					}

					if !tC.hasIssueDate && dInv.Date != nil {
						t.Errorf("invalid dInv.Date %v, want nil", dInv.Date)
					}

					if len(dInv.Items) != tC.nItems {
						t.Errorf("invalid dInv.Items number %d, want %d", len(dInv.Items), tC.nItems)
					}
				})
			}
		}
	})
}

func TestAddInvoice(t *testing.T) {
	t.Run("builds correct DynamoDB input", func(t *testing.T) {
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

	t.Run("handles DynamoDB errors", func(t *testing.T) {
		client := mocks.NewDynamoAPI(mocks.WithPutItemError(errors.New("DynamoDB PutItem failed")))
		strg := dynamo.New(client, "invoices")
		inv := invoice.NewInvoice("123", "customer")

		err := strg.AddInvoice(inv)
		if err == nil {
			t.Errorf("expected AddInvoice(%v) to fail", inv)
		} else if got, want := err.Error(), `DynamoDB PutItem failed`; got != want {
			t.Errorf("AddInvoice(%v) = %v, want %v", inv, got, want)
		}
	})
}

func TestFindInvoice(t *testing.T) {
	t.Run("builds correct DynamoDB input", func(t *testing.T) {
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
	})

	t.Run("handles DynamoDB errors", func(t *testing.T) {
		client := mocks.NewDynamoAPI(mocks.WithGetItemError(errors.New("DynamoDB GetItem failed")))
		strg := dynamo.New(client, "invoices")
		invID := "123"

		_, err := strg.FindInvoice(invID)
		if err == nil {
			t.Errorf("expected FindInvoice(%q) to fail", invID)
		} else if got, want := err.Error(), `DynamoDB GetItem failed`; got != want {
			t.Errorf("FindInvoice(%q) = %v, want %v", invID, got, want)
		}
	})
}

func TestUpdateInvoice(t *testing.T) {
	t.Run("builds correct DynamoDB input", func(t *testing.T) {
		client := mocks.NewDynamoAPI()
		strg := dynamo.New(client, "invoices")
		inv := invoice.NewInvoice("123", "customer")

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
	})

	t.Run("handles DynamoDB errors", func(t *testing.T) {
		client := mocks.NewDynamoAPI(mocks.WithPutItemError(errors.New("DynamoDB PutItem failed")))
		strg := dynamo.New(client, "invoices")
		inv := invoice.NewInvoice("123", "customer")

		err := strg.UpdateInvoice(inv)
		if err == nil {
			t.Errorf("expected UpdateInvoice(%v) to fail", inv)
		} else if got, want := err.Error(), `DynamoDB PutItem failed`; got != want {
			t.Errorf("UpdateInvoice(%v) = %v, want %v", inv, got, want)
		}
	})
}
