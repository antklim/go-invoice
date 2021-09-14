package memory_test

import (
	"testing"
	"time"

	"github.com/antklim/go-invoice/invoice"
	"github.com/antklim/go-invoice/storage/memory"
)

func TestAddInvoice(t *testing.T) {
	strg := memory.New()
	inv := invoice.NewInvoice("123", "customer")

	if err := strg.AddInvoice(inv); err != nil {
		t.Errorf("AddInvoice(%v) failed: %v", inv, err)
	}

	// TODO: move second call check to the service level
	err := strg.AddInvoice(inv)
	if err == nil {
		t.Errorf("expected second call AddInvoice(%v) to fail", inv)
	} else if got, want := err.Error(), `ID "123" exists`; got != want {
		t.Errorf("second call AddInvoice(%v) = %v, want %v", inv, got, want)
	}
}

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

func TestUpdateInvoiceError(t *testing.T) {
	strg := memory.New()
	inv := invoice.NewInvoice("123", "customer")

	err := strg.UpdateInvoice(inv)
	if err == nil {
		t.Errorf("expected UpdateInvoice(%v) to fail", inv)
	} else if got, want := err.Error(), `invoice "123" not found`; got != want {
		t.Errorf("UpdateInvoice(%v) = %v, want %v", inv, got, want)
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
