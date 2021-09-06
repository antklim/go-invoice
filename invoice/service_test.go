package invoice_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/antklim/go-invoice/invoice"
	"github.com/antklim/go-invoice/storage"
	testapi "github.com/antklim/go-invoice/test/api"
	"github.com/antklim/go-invoice/test/mocks"
	"github.com/google/uuid"
)

func testSetup() (*invoice.Service, *testapi.Invoice, error) {
	strg, err := storage.Factory("memory")
	if err != nil {
		return nil, nil, err
	}
	srv := invoice.New(strg)
	api := testapi.NewIvoiceAPI(strg)
	return srv, api, nil
}

func TestCreateInvoice(t *testing.T) {
	srv, _, err := testSetup()
	if err != nil {
		t.Fatalf("testSetup() failed: %v", err)
	}

	t.Run("successfully stores the invoice", func(t *testing.T) {
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

		vinv, err := srv.ViewInvoice(inv.ID)
		if err != nil {
			t.Fatalf("ViewInvoice(%q) failed: %v", inv.ID, err)
		}
		if !inv.Equal(vinv) {
			t.Errorf("invalid invoice %v, want %v", vinv, inv)
		}
	})

	t.Run("propagates data storage errors", func(t *testing.T) {
		e := errors.New("storage failed to add invoice")
		strg := mocks.NewStorage(mocks.WithAddInvoiceError(e))
		srv := invoice.New(strg)

		customer := "John Doe"
		_, err := srv.CreateInvoice(customer)
		if err == nil {
			t.Fatalf("expected CreateInvoice(%q) to fail due to storage error", customer)
		}
		if got, want := err.Error(), fmt.Sprintf("create invoice failed: %s", e.Error()); got != want {
			t.Errorf("CreateInvoice(%q) failed with: %s, want %s", customer, got, want)
		}
	})
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
		if !inv.Equal(vinv) {
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
					inv.ID, customer, inv.Status)
			}
			got := err.Error()
			want := fmt.Sprintf("%q invoice cannot be updated", inv.Status)
			if got != want {
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
		item := testapi.ItemFactory()
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
			item := testapi.ItemFactory()
			err := srv.AddInvoiceItem(inv.ID, item)
			if err == nil {
				t.Fatalf("expected AddInvoiceItems(%q, %v) to fail when invoice status is %q",
					inv.ID, item, inv.Status)
			}
			got := err.Error()
			want := fmt.Sprintf("item cannot be added to %q invoice", inv.Status)
			if got != want {
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

		nitems := len(inv.Items)

		// add item
		item := testapi.ItemFactory()
		if err := srv.AddInvoiceItem(inv.ID, item); err != nil {
			t.Fatalf("AddInvoiceItems(%q, %v) failed: %v", inv.ID, item, err)
		}

		// validate that item added
		vinv, err := srv.ViewInvoice(inv.ID)
		if err != nil {
			t.Fatalf("ViewInvoice(%q) failed: %v", inv.ID, err)
		}
		if len(vinv.Items) != nitems+1 {
			t.Errorf("invalid invoice.Items number %d, want %d", len(vinv.Items), nitems+1)
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
		item := testapi.ItemFactory()
		err := srv.DeleteInvoiceItem(invID, item.ID)
		if err == nil {
			t.Fatalf("expected DeleteInvoiceItem(%q, %q) to fail when invoice does not exist", invID, item.ID)
		}
		if got, want := err.Error(), fmt.Sprintf("invoice %q not found", invID); got != want {
			t.Errorf("DeleteInvoiceItem(%q, %q) failed with: %s, want %s", invID, item.ID, got, want)
		}
	})

	t.Run("fails when invoice is in the status other than open", func(t *testing.T) {
		statuses := []invoice.Status{invoice.Issued, invoice.Paid, invoice.Canceled}
		invoices, err := invoiceAPI.CreateInvoicesWithStatuses(statuses...)
		if err != nil {
			t.Fatalf("invoiceAPI.CreateInvoicesWithStatuses() failed: %v", err)
		}

		itemID := uuid.Nil.String()
		for _, inv := range invoices {
			err := srv.DeleteInvoiceItem(inv.ID, itemID)
			if err == nil {
				t.Fatalf("expected DeleteInvoiceItem(%q, %q) to fail when invoice status is %q",
					inv.ID, itemID, inv.Status)
			}
			got := err.Error()
			want := fmt.Sprintf("item cannot be deleted from %q invoice", inv.Status)
			if got != want {
				t.Errorf("DeleteInvoiceItem(%q, %q) failed with: %s, want %s", inv.ID, itemID, got, want)
			}
		}
	})

	t.Run("fails when data storage error occurred", func(t *testing.T) {
		// search failed
		// update failed
	})

	t.Run("successfully deletes invoice item", func(t *testing.T) {
		nitems := 3
		// place open invoice
		inv, err := invoiceAPI.CreateInvoiceWithNItems(nitems)
		if err != nil {
			t.Fatalf("invoiceAPI.CreateInvoice() failed: %v", err)
		}

		// delete an item
		itemID := inv.Items[0].ID
		if err := srv.DeleteInvoiceItem(inv.ID, itemID); err != nil {
			t.Fatalf("DeleteInvoiceItem(%q, %q) failed: %v", inv.ID, itemID, err)
		}

		// validate that item deleted
		vinv, err := srv.ViewInvoice(inv.ID)
		if err != nil {
			t.Fatalf("ViewInvoice(%q) failed: %v", inv.ID, err)
		}
		if len(vinv.Items) != nitems-1 {
			t.Errorf("invalid invoice.Items number %d, want %d", len(vinv.Items), nitems-1)
		}
		if vinv.ContainsItem(itemID) {
			t.Errorf("invoice %q should not contain item %q", vinv.ID, itemID)
		}
		if !vinv.UpdatedAt.After(inv.UpdatedAt) {
			t.Errorf("invalid invoice.UpdatedAt %s, want it to be after %s",
				vinv.UpdatedAt.Format(time.RFC3339Nano), inv.UpdatedAt.Format(time.RFC3339Nano))
		}
	})

	t.Run("idempotent to repeatable delete", func(t *testing.T) {
		nitems := 3
		// place open invoice
		inv, err := invoiceAPI.CreateInvoiceWithNItems(nitems)
		if err != nil {
			t.Fatalf("invoiceAPI.CreateInvoice() failed: %v", err)
		}

		// delete an item
		itemID := inv.Items[0].ID
		if err := srv.DeleteInvoiceItem(inv.ID, itemID); err != nil {
			t.Fatalf("DeleteInvoiceItem(%q, %q) failed: %v", inv.ID, itemID, err)
		}

		vinv1, err := srv.ViewInvoice(inv.ID)
		if err != nil {
			t.Fatalf("ViewInvoice(%q) failed: %v", inv.ID, err)
		}

		// repeat deletion
		if err := srv.DeleteInvoiceItem(inv.ID, itemID); err != nil {
			t.Fatalf("DeleteInvoiceItem(%q, %q) failed: %v", inv.ID, itemID, err)
		}

		vinv2, err := srv.ViewInvoice(inv.ID)
		if err != nil {
			t.Fatalf("ViewInvoice(%q) failed: %v", inv.ID, err)
		}

		if !vinv1.Equal(vinv2) {
			t.Errorf("invalid invoice %v, want %v", vinv2, vinv1)
		}
	})
}

func TestIssueInvoice(t *testing.T) {
	srv, invoiceAPI, err := testSetup()
	if err != nil {
		t.Fatalf("testSetup() failed: %v", err)
	}

	t.Run("fails when no invoice found", func(t *testing.T) {
		invID := uuid.Nil.String()
		err := srv.IssueInvoice(invID)
		if err == nil {
			t.Fatalf("expected IssueInvoice(%q) to fail when invoice does not exist", invID)
		}
		if got, want := err.Error(), fmt.Sprintf("invoice %q not found", invID); got != want {
			t.Errorf("IssueInvoice(%q) failed with: %s, want %s", invID, got, want)
		}
	})

	t.Run("fails when invoice is in the status other than open", func(t *testing.T) {
		statuses := []invoice.Status{invoice.Issued, invoice.Paid, invoice.Canceled}
		invoices, err := invoiceAPI.CreateInvoicesWithStatuses(statuses...)
		if err != nil {
			t.Fatalf("invoiceAPI.CreateInvoicesWithStatuses() failed: %v", err)
		}

		for _, inv := range invoices {
			err := srv.IssueInvoice(inv.ID)
			if err == nil {
				t.Fatalf("expected IssueInvoice(%q) to fail when invoice status is %q",
					inv.ID, inv.Status)
			}
			got := err.Error()
			want := fmt.Sprintf("%q invoice cannot be issued", inv.Status)
			if got != want {
				t.Errorf("IssueInvoice(%q) failed with: %s, want %s", inv.ID, got, want)
			}
		}
	})

	t.Run("fails when data storage error occurred", func(t *testing.T) {
		// search failed
		// update failed
	})

	t.Run("successfully issues invoice", func(t *testing.T) {
		// place open invoice
		inv, err := invoiceAPI.CreateInvoice()
		if err != nil {
			t.Fatalf("invoiceAPI.CreateInvoice() failed: %v", err)
		}

		// issue invoice
		if err := srv.IssueInvoice(inv.ID); err != nil {
			t.Fatalf("IssueInvoice(%q) failed: %v", inv.ID, err)
		}

		// validate that invoice fields were respectively updated
		vinv, err := srv.ViewInvoice(inv.ID)
		if err != nil {
			t.Fatalf("ViewInvoice(%q) failed: %v", inv.ID, err)
		}
		if vinv.Status != invoice.Issued {
			t.Errorf("invalid invoice.Status %d, want %d", vinv.Status, invoice.Issued)
		}
		if vinv.Date == nil {
			t.Error("invoice issue date should be set")
		}
		if time.Since(*vinv.Date).Milliseconds() > 100 {
			t.Error("invoice issue date should be within the last 100 msec")
		}
		if !vinv.UpdatedAt.After(inv.UpdatedAt) {
			t.Errorf("invalid invoice.UpdatedAt %s, want it to be after %s",
				vinv.UpdatedAt.Format(time.RFC3339Nano), inv.UpdatedAt.Format(time.RFC3339Nano))
		}
	})
}

func TestPayInvoice(t *testing.T) {
	srv, invoiceAPI, err := testSetup()
	if err != nil {
		t.Fatalf("testSetup() failed: %v", err)
	}

	t.Run("fails when no invoice found", func(t *testing.T) {
		invID := uuid.Nil.String()
		err := srv.PayInvoice(invID)
		if err == nil {
			t.Fatalf("expected PayInvoice(%q) to fail when invoice does not exist", invID)
		}
		if got, want := err.Error(), fmt.Sprintf("invoice %q not found", invID); got != want {
			t.Errorf("PayInvoice(%q) failed with: %s, want %s", invID, got, want)
		}
	})

	t.Run("fails when invoice is in the status other than issued", func(t *testing.T) {
		statuses := []invoice.Status{invoice.Open, invoice.Paid, invoice.Canceled}
		invoices, err := invoiceAPI.CreateInvoicesWithStatuses(statuses...)
		if err != nil {
			t.Fatalf("invoiceAPI.CreateInvoicesWithStatuses() failed: %v", err)
		}

		for _, inv := range invoices {
			err := srv.PayInvoice(inv.ID)
			if err == nil {
				t.Fatalf("expected PayInvoice(%q) to fail when invoice status is %q",
					inv.ID, inv.Status)
			}
			got := err.Error()
			want := fmt.Sprintf("%q invoice cannot be paid", inv.Status)
			if got != want {
				t.Errorf("PayInvoice(%q) failed with: %s, want %s", inv.ID, got, want)
			}
		}
	})

	t.Run("fails when data storage error occurred", func(t *testing.T) {
		// search failed
		// update failed
	})

	t.Run("successfully pays invoice", func(t *testing.T) {
		// place issued invoice
		inv, err := invoiceAPI.CreateInvoice(testapi.WithStatus(invoice.Issued))
		if err != nil {
			t.Fatalf("invoiceAPI.CreateInvoice() failed: %v", err)
		}

		// pay invoice
		if err := srv.PayInvoice(inv.ID); err != nil {
			t.Fatalf("PayInvoice(%q) failed: %v", inv.ID, err)
		}

		// validate that invoice fields were respectively updated
		vinv, err := srv.ViewInvoice(inv.ID)
		if err != nil {
			t.Fatalf("ViewInvoice(%q) failed: %v", inv.ID, err)
		}
		if vinv.Status != invoice.Paid {
			t.Errorf("invalid invoice.Status %q, want %q", vinv.Status, invoice.Paid)
		}
		if !vinv.UpdatedAt.After(inv.UpdatedAt) {
			t.Errorf("invalid invoice.UpdatedAt %s, want it to be after %s",
				vinv.UpdatedAt.Format(time.RFC3339Nano), inv.UpdatedAt.Format(time.RFC3339Nano))
		}
	})
}

func TestCancelInvoice(t *testing.T) {
	srv, invoiceAPI, err := testSetup()
	if err != nil {
		t.Fatalf("testSetup() failed: %v", err)
	}

	t.Run("fails when no invoice found", func(t *testing.T) {
		invID := uuid.Nil.String()
		err := srv.CancelInvoice(invID)
		if err == nil {
			t.Fatalf("expected CancelInvoice(%q) to fail when invoice does not exist", invID)
		}
		if got, want := err.Error(), fmt.Sprintf("invoice %q not found", invID); got != want {
			t.Errorf("CancelInvoice(%q) failed with: %s, want %s", invID, got, want)
		}
	})

	t.Run("fails when invoice is in the paid or canceled status", func(t *testing.T) {
		statuses := []invoice.Status{invoice.Paid, invoice.Canceled}
		invoices, err := invoiceAPI.CreateInvoicesWithStatuses(statuses...)
		if err != nil {
			t.Fatalf("invoiceAPI.CreateInvoicesWithStatuses() failed: %v", err)
		}

		for _, inv := range invoices {
			err := srv.CancelInvoice(inv.ID)
			if err == nil {
				t.Fatalf("expected CancelInvoice(%q) to fail when invoice status is %q",
					inv.ID, inv.Status)
			}
			got := err.Error()
			want := fmt.Sprintf("%q invoice cannot be canceled", inv.Status)
			if got != want {
				t.Errorf("CancelInvoice(%q) failed with: %s, want %s", inv.ID, got, want)
			}
		}
	})

	t.Run("fails when data storage error occurred", func(t *testing.T) {
		// search failed
		// update failed
	})

	t.Run("successfully cancels invoice", func(t *testing.T) {
		// place open invoice
		inv, err := invoiceAPI.CreateInvoice()
		if err != nil {
			t.Fatalf("invoiceAPI.CreateInvoice() failed: %v", err)
		}

		// cancel invoice
		if err := srv.CancelInvoice(inv.ID); err != nil {
			t.Fatalf("CancelInvoice(%q) failed: %v", inv.ID, err)
		}

		// validate that invoice fields were respectively updated
		vinv, err := srv.ViewInvoice(inv.ID)
		if err != nil {
			t.Fatalf("ViewInvoice(%q) failed: %v", inv.ID, err)
		}
		if vinv.Status != invoice.Canceled {
			t.Errorf("invalid invoice.Status %q, want %q", vinv.Status, invoice.Canceled)
		}
		if !vinv.UpdatedAt.After(inv.UpdatedAt) {
			t.Errorf("invalid invoice.UpdatedAt %s, want it to be after %s",
				vinv.UpdatedAt.Format(time.RFC3339Nano), inv.UpdatedAt.Format(time.RFC3339Nano))
		}
	})
}
