package memory_test

import (
	"testing"
	"time"

	"github.com/antklim/go-invoice/invoice"
	"github.com/antklim/go-invoice/storage/memory"
)

func TestFindInvoice(t *testing.T) {
	strg := memory.New()
	inv := invoice.NewInvoice("123", "customer")

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

func TestUpdateInvoice(t *testing.T) {
	strg := memory.New()
	inv := invoice.NewInvoice("123", "customer")

	if err := strg.AddInvoice(inv); err != nil {
		t.Errorf("AddInvoice(%v) failed: %v", inv, err)
	}

	newCustomer := "new customer"
	inv.CustomerName = newCustomer
	if err := strg.UpdateInvoice(inv); err != nil {
		t.Errorf("UpdateInvoice(%v) failed: %v", inv, err)
	}

	vinv, err := strg.FindInvoice(inv.ID)
	if err != nil {
		t.Errorf("FindInvoice(%q) failed: %v", inv.ID, err)
	}
	if vinv.CustomerName != newCustomer {
		t.Errorf("invalid updated invoice.CustomerName %q, want %q", vinv.CustomerName, newCustomer)
	}
	if !vinv.UpdatedAt.After(inv.UpdatedAt) {
		t.Errorf("invalid udated invoice.UpdatedAt %s, want after %s",
			vinv.UpdatedAt.Format(time.RFC3339),
			inv.UpdatedAt.Format(time.RFC3339))
	}
}
