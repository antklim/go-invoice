package invoice_test

import (
	"testing"
	"time"

	"github.com/antklim/go-invoice/invoice"
	"github.com/antklim/go-invoice/storage"
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

	t.Run("returns a zero invoice when no invoice is found in data storage", func(t *testing.T) {
		invID := uuid.Nil.String()
		inv, err := srv.ViewInvoice(invID)
		if err != nil {
			t.Fatalf("ViewInvoice(%q) failed: %v", invID, err)
		}
		if inv != nil {
			t.Errorf("invalid invoice %v, want nil", inv)
		}
	})

	t.Run("propagates data storage errors", func(t *testing.T) {})
}

func TestUpdateInvoice(t *testing.T) {
	t.Run("propagates data storage errors", func(t *testing.T) {})
}

func TestCloseInvoice(t *testing.T) {
	t.Run("propagates data storage errors when canceling invoice", func(t *testing.T) {})

	t.Run("propagates data storage errors when paying invoice", func(t *testing.T) {})
}

// Following are the business rules tests
func TestOpenInvoice(t *testing.T) {
	strg, _ := storage.Factory("memory")
	srv := invoice.New(strg)
	inv, _ := srv.CreateInvoice("John Doe")

	t.Run("can be viewed", func(t *testing.T) {
		vinv, err := srv.ViewInvoice(inv.ID)
		if err != nil {
			t.Fatalf("ViewInvoice(%q) failed: %v", inv.ID, err)
		}
		if !vinv.Equal(inv) {
			t.Errorf("invalid invoice %v, want %v", vinv, inv)
		}
	})

	t.Run("can be updated", func(t *testing.T) {
		t.Run("customer name can be updated", func(t *testing.T) {
			newCustomer := "John Wick"
			if err := srv.UpdateInvoiceCustomer(inv.ID, newCustomer); err != nil {
				t.Fatalf("UpdateCustomer(%q) failed: %v", newCustomer, err)
			}
			vinv, err := srv.ViewInvoice(inv.ID)
			if err != nil {
				t.Fatalf("ViewInvoice(%q) failed: %v", inv.ID, err)
			}
			if vinv.CustomerName != newCustomer {
				t.Errorf("invalid invoice.CustomerName %q, want %q", vinv.CustomerName, newCustomer)
			}
		})

		t.Run("items can be added", func(t *testing.T) {})

		t.Run("items can be deleted", func(t *testing.T) {
			// when deleting non existent item it does not return error
			// error returned only in case of data access layer
		})
	})

	t.Run("cannot be updated", func(t *testing.T) {
		t.Run("when invoice not found", func(t *testing.T) {})
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
		t.Run("customer name cannot updated", func(t *testing.T) {})

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
		t.Run("customer name cannot updated", func(t *testing.T) {})

		t.Run("issue date cannot be updated", func(t *testing.T) {})

		t.Run("items cannot be added", func(t *testing.T) {})

		t.Run("items cannot be deleted", func(t *testing.T) {})
	})

	t.Run("cannot be issued", func(t *testing.T) {})

	t.Run("cannot be canceled", func(t *testing.T) {})

	t.Run("cannot be paid", func(t *testing.T) {})
}
