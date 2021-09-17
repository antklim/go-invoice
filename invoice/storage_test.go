package invoice_test

import (
	"testing"

	"github.com/antklim/go-invoice/invoice"
	"github.com/antklim/go-invoice/storage"
)

// TODO: switch between different storages based on env variable

func TestAddInvoiceFails(t *testing.T) {
	t.Run("when repeat adding existing invoice", func(t *testing.T) {
		f := new(storage.Memory)
		strg := f.MakeStorage()

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
	})
}

func TestFindInvoice(t *testing.T) {
	t.Run("returns nil invoice when no invoices found", func(t *testing.T) {
		f := new(storage.Memory)
		strg := f.MakeStorage()

		invID := "123"
		inv, err := strg.FindInvoice(invID)
		if err != nil {
			t.Errorf("FindInvoice(%q) failed: %v", invID, err)
		}
		if inv != nil {
			t.Errorf("FindInvoice(%q) no invoice expected, got %v", invID, inv)
		}
	})
}

func TestUpdateInvoiceFails(t *testing.T) {
	t.Run("when updating non-existing invoice", func(t *testing.T) {
		f := new(storage.Memory)
		strg := f.MakeStorage()
		inv := invoice.NewInvoice("123", "customer")

		err := strg.UpdateInvoice(inv)
		if err == nil {
			t.Errorf("expected UpdateInvoice(%v) to fail", inv)
		} else if got, want := err.Error(), `invoice "123" not found`; got != want {
			t.Errorf("UpdateInvoice(%v) = %v, want %v", inv, got, want)
		}
	})
}
