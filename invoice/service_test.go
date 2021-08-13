package invoice_test

import (
	"testing"
	"time"

	"github.com/antklim/go-invoice/invoice"
	"github.com/google/uuid"
)

func TestCreateInvoice(t *testing.T) {
	srv := invoice.Service{}

	t.Run("creates valid invoice", func(t *testing.T) {
		inv, err := srv.CreateInvoice("John Doe")
		if err != nil {
			t.Fatalf("error creating invoice: %v", err)
		}

		if inv.ID == "" {
			t.Fatal("invoice.ID should not be empty")
		}

		if inv.Date != nil {
			t.Fatalf("invoice.Date should be nil")
		}

		if expected := "John Doe"; inv.CustomerName != expected {
			t.Fatalf("invalid invoice.CustomerName: want=%s, but got=%s", expected, inv.CustomerName)
		}

		if expected := invoice.Open; inv.Status != expected {
			t.Fatalf("invalid invoice.Status: want=%d, but got=%d", expected, inv.Status)
		}

		if !inv.CreatedAt.Equal(inv.UpdatedAt) {
			t.Fatalf("invoice.CreatedAt = %s is not equalt to invoice.UpdatedAt = %s",
				inv.CreatedAt.Format(time.RFC3339),
				inv.UpdatedAt.Format(time.RFC3339))
		}
	})

	t.Run("propagates data storage errors", func(t *testing.T) {})
}

func TestViewInvoice(t *testing.T) {
	srv := invoice.Service{}

	t.Run("returns a zero invoice when no invoice is found in data storage", func(t *testing.T) {
		inv, err := srv.ViewInvoice(uuid.Nil.String())
		if err != nil {
			t.Fatalf("error viewing invoice: %v", err)
		}
		if inv != nil {
			t.Fatalf("invalid invoice: should be empty, but got=%v", err)
		}
	})

	t.Run("returns invoice from data storage", func(t *testing.T) {
		inv, err := srv.CreateInvoice("John Doe")
		if err != nil {
			t.Fatalf("error creating invoice: %v", err)
		}

		vinv, err := srv.ViewInvoice(inv.ID)
		if err != nil {
			t.Fatalf("error viewing invoice: %v", err)
		}
		if !vinv.Equal(inv) {
			t.Fatalf("invalid invoice: want=%v, but got=%v", inv, vinv)
		}
	})

	t.Run("propagates data storage errors", func(t *testing.T) {})
}

func TestUpdateInvoice(t *testing.T) {
	t.Run("propagates data storage errors", func(t *testing.T) {})
}

func TestCloseInvoice(t *testing.T) {
	t.Run("propagates data storage errors when cancelling invoice", func(t *testing.T) {})

	t.Run("propagates data storage errors when paying invoice", func(t *testing.T) {})
}

// Following are the business rules tests
func TestOpenInvoice(t *testing.T) {
	t.Run("can be viewed", func(t *testing.T) {})

	t.Run("can be updated", func(t *testing.T) {
		t.Run("customer name can be updated", func(t *testing.T) {})

		t.Run("items can be added", func(t *testing.T) {})

		t.Run("items can be deleted", func(t *testing.T) {
			// when deleting non existent item it does not return error
			// error returned only in case of data access layer
		})
	})

	t.Run("cannot be updated", func(t *testing.T) {
		t.Run("issue date cannot be updated", func(t *testing.T) {})

		t.Run("items cannot be updated", func(t *testing.T) {})
	})

	t.Run("can be issued", func(t *testing.T) {
		// verify issue date is set
	})

	t.Run("can be cancelled", func(t *testing.T) {})

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

	t.Run("can be cancelled", func(t *testing.T) {})

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

	t.Run("cannot be cancelled", func(t *testing.T) {})

	t.Run("cannot be paid", func(t *testing.T) {})
}
