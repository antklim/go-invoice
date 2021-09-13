package dynamo_test

import (
	"testing"
)

func TestAddInvoice(t *testing.T) {
	t.Skip("WIP")

	t.Run("called with correct input", func(t *testing.T) {
		// client := mocks.DynamoAPI{}
		// strg := dynamo.New(&client, "invoices")
		// inv := invoice.NewInvoice("123", "customer")
		// inv.AddItem(invoice.NewItem("456", "pen", 1000, 3))

		// if err := strg.AddInvoice(inv); err != nil {
		// 	t.Errorf("AddInvoice(%v) failed: %v", inv, err)
		// }

		// if got, want := client.CalledTimes("PutItem"), 1; got != want {
		// 	t.Errorf("client.PutItem() called %d times, want %d call(s)", got, want)
		// }

		// input := client.NthCall("PutItem", 1)
		// expectedInput := &dynamodb.PutItemInput{
		// 	TableName: aws.String("invoices"),
		// }
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
