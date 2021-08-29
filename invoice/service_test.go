package invoice_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/antklim/go-invoice/invoice"
	"github.com/antklim/go-invoice/storage"
	invapi "github.com/antklim/go-invoice/test/api"
	"github.com/google/uuid"
)

func TestCreateInvoice(t *testing.T) {
	strg, _ := storage.Factory("memory")
	srv := invoice.New(strg)

	t.Run("creates valid invoice", func(t *testing.T) {
		customer := "John Doe"
		inv, err := srv.CreateInvoice(customer)
		if err != nil {
			t.Fatalf("CreateInvoice(%q) failed: %v", customer, err)
		}

		if inv.ID == "" {
			t.Errorf("invoice.ID should not be empty")
		}

		if inv.Date != nil {
			t.Errorf("invalid invoice.Date, want nil")
		}

		if inv.CustomerName != customer {
			t.Errorf("invalid invoice.CustomerName %q, want %q", inv.CustomerName, customer)
		}

		if inv.Status != invoice.Open {
			t.Errorf("invalid invoice.Status %d, want %d", inv.Status, invoice.Open)
		}

		if !inv.CreatedAt.Equal(inv.UpdatedAt) {
			t.Errorf("invoice.CreatedAt = %s is not equalt to invoice.UpdatedAt = %s",
				inv.CreatedAt.Format(time.RFC3339),
				inv.UpdatedAt.Format(time.RFC3339))
		}
	})

	t.Run("propagates data storage errors", func(t *testing.T) {})
}

func TestViewInvoice(t *testing.T) {
	strg, _ := storage.Factory("memory")
	srv := invoice.New(strg)
	invoiceAPI := invapi.NewIvoiceAPI(strg)

	t.Run("returns nil when no invoice is found in data storage", func(t *testing.T) {
		invID := uuid.Nil.String()
		inv, err := srv.ViewInvoice(invID)
		if err != nil {
			t.Fatalf("ViewInvoice(%q) failed: %v", invID, err)
		}
		if inv != nil {
			t.Errorf("invalid invoice %v, want nil", inv)
		}
	})

	t.Run("returns invoice", func(t *testing.T) {
		invID, err := invoiceAPI.CreateInvoice()
		if err != nil {
			t.Fatalf("invoiceAPI.CreateInvoice() failed: %v", err)
		}

		vinv, err := srv.ViewInvoice(invID)
		if err != nil {
			t.Fatalf("ViewInvoice(%q) failed: %v", invID, err)
		}
		if vinv.ID != invID {
			t.Errorf("invalid invoice.ID %s, want %s", vinv.ID, invID)
		}
	})

	t.Run("propagates data storage errors", func(t *testing.T) {})
}

func TestUpdateInvoiceCustomer(t *testing.T) {
	strg, _ := storage.Factory("memory")
	srv := invoice.New(strg)
	invoiceAPI := invapi.NewIvoiceAPI(strg)

	t.Run("fails when no invoice found", func(t *testing.T) {
		invID := uuid.Nil.String()
		customer := "John Doe"
		err := srv.UpdateInvoiceCustomer(invID, customer)
		if err == nil {
			t.Fatalf("expected UpdateInvoiceCustomer(%q, %q) to fail when invoice does not exist", invID, customer)
		}
		if got, want := err.Error(), fmt.Sprintf("invoice %q not found", invID); got != want {
			t.Errorf("UpdateInvoiceCustomer(%q, %q) failed with: %s, want %s", invID, customer, got, want)
		}
	})

	t.Run("fails when invoice is in the status other than open", func(t *testing.T) {})

	t.Run("fails when data storage error occurred", func(t *testing.T) {
		// search failed
		// update failed
	})

	t.Run("successfully updates customer name of open invoice", func(t *testing.T) {
		// place open invoice
		invID, err := invoiceAPI.CreateInvoice(
			invapi.WithCustomerName("John Doe"),
			invapi.WithStatus(invoice.Open),
		)
		if err != nil {
			t.Fatalf("invoiceAPI.CreateInvoice() failed: %v", err)
		}

		// update customer name
		customer := "John Wick"
		if err := srv.UpdateInvoiceCustomer(invID, customer); err != nil {
			t.Fatalf("UpdateCustomer(%q) failed: %v", customer, err)
		}

		// validate that customer name updated
		inv, err := srv.ViewInvoice(invID)
		if err != nil {
			t.Fatalf("ViewInvoice(%q) failed: %v", invID, err)
		}
		if inv.CustomerName != customer {
			t.Errorf("invalid invoice.CustomerName %q, want %q", inv.CustomerName, customer)
		}
	})
}

func TestAddInvoiceItemFails(t *testing.T) {
	t.Run("when no invoice found", func(t *testing.T) {})
	t.Run("when data storage error occurred", func(t *testing.T) {})
}

func TestDeleteInvoiceItemFails(t *testing.T) {
	t.Run("when no invoice found", func(t *testing.T) {})
	t.Run("when data storage error occurred", func(t *testing.T) {})
}

func TestPayInvoiceFails(t *testing.T) {
	t.Run("when no invoice found", func(t *testing.T) {})
	t.Run("when data storage error occurred", func(t *testing.T) {})
}

func TestCancelInvoiceFails(t *testing.T) {
	t.Run("when no invoice found", func(t *testing.T) {})
	t.Run("when data storage error occurred", func(t *testing.T) {})
}

// Following are the business rules tests
func TestOpenInvoice(t *testing.T) {
	t.Run("can be updated", func(t *testing.T) {
		t.Run("items can be added", func(t *testing.T) {})

		t.Run("items can be deleted", func(t *testing.T) {
			// when deleting non existent item it does not return error
			// error returned only in case of data access layer
		})
	})

	t.Run("can be issued", func(t *testing.T) {
		// verify issue date is set
	})

	t.Run("can be canceled", func(t *testing.T) {})

	t.Run("cannot be paid", func(t *testing.T) {})
}

func TestIssuedInvoice(t *testing.T) {
	t.Run("can be viewed", func(t *testing.T) {})

	t.Run("cannot be updated", func(t *testing.T) {
		t.Run("issue date cannot be updated", func(t *testing.T) {})

		t.Run("items cannot be added", func(t *testing.T) {})

		t.Run("items cannot be deleted", func(t *testing.T) {})
	})

	t.Run("cannot be issued", func(t *testing.T) {})

	t.Run("can be canceled", func(t *testing.T) {})

	t.Run("can be paid", func(t *testing.T) {})
}

func TestClosedInvoice(t *testing.T) {
	t.Run("can be viewed", func(t *testing.T) {})

	t.Run("cannot be updated", func(t *testing.T) {
		t.Run("issue date cannot be updated", func(t *testing.T) {})

		t.Run("items cannot be added", func(t *testing.T) {})

		t.Run("items cannot be deleted", func(t *testing.T) {})
	})

	t.Run("cannot be issued", func(t *testing.T) {})

	t.Run("cannot be canceled", func(t *testing.T) {})

	t.Run("cannot be paid", func(t *testing.T) {})
}
