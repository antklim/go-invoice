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

func TestViewInvoice(t *testing.T) {
	inv := invoice.NewInvoice("123", "customer")

	strg := storage.New()
	vinv, err := strg.FindInvoice(inv.ID)
	if err != nil {
		t.Errorf("FindInvoice(%q) failed: %v", inv.ID, err)
	}
	if vinv != nil {
		t.Errorf("FindInvoice(%q) no invoice expected, got %v", inv.ID, vinv)
	}

	if err := strg.AddInvoice(inv); err != nil {
		t.Errorf("AddInvoice(%v) failed: %v", inv, err)
	}

	vinv, err = strg.FindInvoice(inv.ID)
	if err != nil {
		t.Errorf("FindInvoice(%q) failed: %v", inv.ID, err)
	}
	if vinv == nil {
		t.Errorf("FindInvoice(%q) invoice expected, got nil", inv.ID)
	}
}
