package invoice_test

import (
	"fmt"
	"testing"

	"github.com/antklim/go-invoice/invoice"
	"github.com/antklim/go-invoice/storage"
	"github.com/google/uuid"
)

func storageSetup() invoice.Storage {
	f := new(storage.Memory)
	strg := f.MakeStorage()
	return strg
}

func TestAddInvoiceFails(t *testing.T) {
	t.Run("when repeat adding existing invoice", func(t *testing.T) {
		strg := storageSetup()

		invID := uuid.NewString()
		inv := invoice.NewInvoice(invID, "John Doe")
		if err := strg.AddInvoice(inv); err != nil {
			t.Errorf("AddInvoice(%v) failed: %v", inv, err)
		}

		err := strg.AddInvoice(inv)
		if err == nil {
			t.Errorf("expected second call AddInvoice(%v) to fail", inv)
		} else if got, want := err.Error(), fmt.Sprintf("ID %q exists", invID); got != want {
			t.Errorf("second call AddInvoice(%v) = %v, want %v", inv, got, want)
		}
	})
}

func TestFindInvoice(t *testing.T) {
	t.Run("returns nil invoice when no invoices found", func(t *testing.T) {
		strg := storageSetup()

		invID := uuid.Nil.String()
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
		strg := storageSetup()

		invID := uuid.Nil.String()
		inv := invoice.NewInvoice(invID, "John Doe")

		err := strg.UpdateInvoice(inv)
		if err == nil {
			t.Errorf("expected UpdateInvoice(%v) to fail", inv)
		} else if got, want := err.Error(), fmt.Sprintf("invoice %q not found", invID); got != want {
			t.Errorf("UpdateInvoice(%v) = %v, want %v", inv, got, want)
		}
	})
}
