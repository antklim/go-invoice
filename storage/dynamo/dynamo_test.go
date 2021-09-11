package dynamo_test

import (
	"testing"

	"github.com/antklim/go-invoice/invoice"
	"github.com/antklim/go-invoice/storage/dynamo"
	"github.com/antklim/go-invoice/test/mocks"
)

func TestAddInvoice(t *testing.T) {
	t.Skip("WIP")
	client := mocks.DynamoAPI{}
	strg := dynamo.New(&client)
	inv := invoice.NewInvoice("123", "customer")

	if err := strg.AddInvoice(inv); err != nil {
		t.Errorf("AddInvoice(%v) failed: %v", inv, err)
	}

	err := strg.AddInvoice(inv)
	if err == nil {
		t.Errorf("expected second call AddInvoice(%v) to fail", inv)
	} else if got, want := err.Error(), `ID "123" exists`; got != want {
		t.Errorf("second call AddInvoice(%v) = %v, want %v", inv, got, want)
	}
}

func TestFindInvoice(t *testing.T) {
	t.Skip("not implemented")
}

func TestUpdateInvoice(t *testing.T) {
	t.Skip("not implemented")
}
