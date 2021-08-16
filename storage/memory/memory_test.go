package memory_test

import (
	"testing"

	"github.com/antklim/go-invoice/invoice"
	storage "github.com/antklim/go-invoice/storage/memory"
)

func TestAddInvoice(t *testing.T) {
	inv := invoice.NewInvoice("123", "customer")

	strg := storage.New()
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
