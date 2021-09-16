package invoice_test

import (
	"fmt"
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
	// TODO: implement
	t.Skip("not implemented")

	desc := "returns nil invoice"
	testCases := []struct {
		when string
	}{
		{
			when: "when storage operation fails",
		},
		{
			when: "when no invoices found",
		},
	}
	for _, tC := range testCases {
		t.Run(fmt.Sprintf("%s %s", desc, tC.when), func(t *testing.T) {
		})
	}
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
