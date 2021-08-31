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

func testSetup() (*invoice.Service, *invapi.Invoice, error) {
	strg, err := storage.Factory("memory")
	if err != nil {
		return nil, nil, err
	}
	srv := invoice.New(strg)
	api := invapi.NewIvoiceAPI(strg)
	return srv, api, nil
}

func TestCreateInvoice(t *testing.T) {
	srv, _, err := testSetup()
	if err != nil {
		t.Fatalf("testSetup() failed: %v", err)
	}

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

		if len(inv.Items) != 0 {
			t.Errorf("invoice.Items should be empty")
		}

		if inv.Status != invoice.Open {
			t.Errorf("invalid invoice.Status %d, want %d", inv.Status, invoice.Open)
		}

		if !inv.CreatedAt.Equal(inv.UpdatedAt) {
			t.Errorf("invoice.CreatedAt = %s is not equal to invoice.UpdatedAt = %s",
				inv.CreatedAt.Format(time.RFC3339Nano),
				inv.UpdatedAt.Format(time.RFC3339Nano))
		}
	})

	t.Run("successfully stores the invoice", func(t *testing.T) {
		customer := "John Doe"
		inv, err := srv.CreateInvoice(customer)
		if err != nil {
			t.Fatalf("CreateInvoice(%q) failed: %v", customer, err)
		}

		vinv, err := srv.ViewInvoice(inv.ID)
		if err != nil {
			t.Fatalf("ViewInvoice(%q) failed: %v", inv.ID, err)
		}
		if !vinv.Equal(inv) {
			t.Errorf("invalid invoice %v, want %v", vinv, inv)
		}
	})

	t.Run("propagates data storage errors", func(t *testing.T) {})
}

func TestViewInvoice(t *testing.T) {
	srv, invoiceAPI, err := testSetup()
	if err != nil {
		t.Fatalf("testSetup() failed: %v", err)
	}

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
		inv, err := invoiceAPI.CreateInvoice()
		if err != nil {
			t.Fatalf("invoiceAPI.CreateInvoice() failed: %v", err)
		}

		vinv, err := srv.ViewInvoice(inv.ID)
		if err != nil {
			t.Fatalf("ViewInvoice(%q) failed: %v", inv.ID, err)
		}
		if !vinv.Equal(inv) {
			t.Errorf("invalid invoice %v, want %v", vinv, inv)
		}
	})

	t.Run("propagates data storage errors", func(t *testing.T) {})
}

func TestUpdateInvoiceCustomer(t *testing.T) {
	srv, invoiceAPI, err := testSetup()
	if err != nil {
		t.Fatalf("testSetup() failed: %v", err)
	}

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

	t.Run("fails when invoice is in the status other than open", func(t *testing.T) {
		statuses := []invoice.Status{invoice.Issued, invoice.Paid, invoice.Canceled}
		invoices, err := invoiceAPI.CreateInvoicesWithStatuses(statuses...)
		if err != nil {
			t.Fatalf("invoiceAPI.CreateInvoicesWithStatuses() failed: %v", err)
		}

		for _, inv := range invoices {
			customer := "John Doe"
			err := srv.UpdateInvoiceCustomer(inv.ID, customer)
			if err == nil {
				t.Fatalf("expected UpdateInvoiceCustomer(%q, %q) to fail when invoice status is %q",
					inv.ID, customer, inv.FormatStatus())
			}
			if got, want := err.Error(), fmt.Sprintf("%q invoice cannot be updated", inv.FormatStatus()); got != want {
				t.Errorf("UpdateInvoiceCustomer(%q, %q) failed with: %s, want %s", inv.ID, customer, got, want)
			}
		}
	})

	t.Run("fails when data storage error occurred", func(t *testing.T) {
		// search failed
		// update failed
	})

	t.Run("successfully updates customer name of open invoice", func(t *testing.T) {
		// place open invoice
		inv, err := invoiceAPI.CreateInvoice()
		if err != nil {
			t.Fatalf("invoiceAPI.CreateInvoice() failed: %v", err)
		}

		// update customer name
		// adding timestamp to name to avoid potential overlap with default name
		customer := fmt.Sprintf("%s%d", "James Bond", time.Now().Unix())
		if err := srv.UpdateInvoiceCustomer(inv.ID, customer); err != nil {
			t.Fatalf("UpdateCustomer(%q, %q) failed: %v", inv.ID, customer, err)
		}

		// validate that customer name updated
		vinv, err := srv.ViewInvoice(inv.ID)
		if err != nil {
			t.Fatalf("ViewInvoice(%q) failed: %v", inv.ID, err)
		}
		if vinv.CustomerName != customer {
			t.Errorf("invalid invoice.CustomerName %q, want %q", vinv.CustomerName, customer)
		}
		if !vinv.UpdatedAt.After(inv.UpdatedAt) {
			t.Errorf("invalid invoice.UpdatedAt %s, want it to be after %s",
				vinv.UpdatedAt.Format(time.RFC3339Nano), inv.UpdatedAt.Format(time.RFC3339Nano))
		}
	})
}

func TestAddInvoiceItem(t *testing.T) {
	srv, invoiceAPI, err := testSetup()
	if err != nil {
		t.Fatalf("testSetup() failed: %v", err)
	}

	t.Run("fails when no invoice found", func(t *testing.T) {
		invID := uuid.Nil.String()
		item := invoiceAPI.ItemFactory()
		err := srv.AddInvoiceItem(invID, item)
		if err == nil {
			t.Fatalf("expected AddInvoiceItems(%q, %v) to fail when invoice does not exist", invID, item)
		}
		if got, want := err.Error(), fmt.Sprintf("invoice %q not found", invID); got != want {
			t.Errorf("AddInvoiceItems(%q, %v) failed with: %s, want %s", invID, item, got, want)
		}
	})

	t.Run("fails when invoice is in the status other than open", func(t *testing.T) {
		statuses := []invoice.Status{invoice.Issued, invoice.Paid, invoice.Canceled}
		invoices, err := invoiceAPI.CreateInvoicesWithStatuses(statuses...)
		if err != nil {
			t.Fatalf("invoiceAPI.CreateInvoicesWithStatuses() failed: %v", err)
		}

		for _, inv := range invoices {
			item := invoiceAPI.ItemFactory()
			err := srv.AddInvoiceItem(inv.ID, item)
			if err == nil {
				t.Fatalf("expected AddInvoiceItems(%q, %v) to fail when invoice status is %q",
					inv.ID, item, inv.FormatStatus())
			}
			if got, want := err.Error(), fmt.Sprintf("item cannot be added to %q invoice", inv.FormatStatus()); got != want {
				t.Errorf("AddInvoiceItems(%q, %v) failed with: %s, want %s", inv.ID, item, got, want)
			}
		}
	})

	t.Run("fails when data storage error occurred", func(t *testing.T) {
		// search failed
		// update failed
	})

	t.Run("successfully adds invoice item", func(t *testing.T) {
		// place open invoice
		inv, err := invoiceAPI.CreateInvoice()
		if err != nil {
			t.Fatalf("invoiceAPI.CreateInvoice() failed: %v", err)
		}

		items := len(inv.Items)

		// add item
		item := invoiceAPI.ItemFactory()
		if err := srv.AddInvoiceItem(inv.ID, item); err != nil {
			t.Fatalf("AddInvoiceItems(%q, %v) failed: %v", inv.ID, item, err)
		}

		// validate that item added
		vinv, err := srv.ViewInvoice(inv.ID)
		if err != nil {
			t.Fatalf("ViewInvoice(%q) failed: %v", inv.ID, err)
		}
		if len(vinv.Items) != items+1 {
			t.Errorf("invalid invoice.Items number %d, want %d", len(vinv.Items), items+1)
		}
		if !vinv.UpdatedAt.After(inv.UpdatedAt) {
			t.Errorf("invalid invoice.UpdatedAt %s, want it to be after %s",
				vinv.UpdatedAt.Format(time.RFC3339Nano), inv.UpdatedAt.Format(time.RFC3339Nano))
		}
	})
}

func TestDeleteInvoiceItem(t *testing.T) {
	srv, invoiceAPI, err := testSetup()
	if err != nil {
		t.Fatalf("testSetup() failed: %v", err)
	}

	t.Run("fails when no invoice found", func(t *testing.T) {
		invID := uuid.Nil.String()
		item := invoiceAPI.ItemFactory()
		err := srv.DeleteInvoiceItem(invID, item.ID)
		if err == nil {
			t.Fatalf("expected DeleteInvoiceItem(%q, %q) to fail when invoice does not exist", invID, item.ID)
		}
		if got, want := err.Error(), fmt.Sprintf("invoice %q not found", invID); got != want {
			t.Errorf("DeleteInvoiceItem(%q, %q) failed with: %s, want %s", invID, item.ID, got, want)
		}
	})

	t.Run("fails when invoice is in the status other than open", func(t *testing.T) {})
	t.Run("fails when data storage error occurred", func(t *testing.T) {
		// search failed
		// update failed
	})
	t.Run("successfully deletes invoice item", func(t *testing.T) {})
	t.Run("idempotent to repeatable delete", func(t *testing.T) {})
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
	t.Run("can be issued", func(t *testing.T) {
		// verify issue date is set
	})

	t.Run("can be canceled", func(t *testing.T) {})

	t.Run("cannot be paid", func(t *testing.T) {})
}

func TestIssuedInvoice(t *testing.T) {
	t.Run("cannot be issued", func(t *testing.T) {})

	t.Run("can be canceled", func(t *testing.T) {})

	t.Run("can be paid", func(t *testing.T) {})
}

func TestClosedInvoice(t *testing.T) {
	t.Run("cannot be issued", func(t *testing.T) {})

	t.Run("cannot be canceled", func(t *testing.T) {})

	t.Run("cannot be paid", func(t *testing.T) {})
}
