package invoice_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/antklim/go-invoice/invoice"
	testapi "github.com/antklim/go-invoice/test/api"
	"github.com/antklim/go-invoice/test/mocks"
	"github.com/google/uuid"
)

func TestCreateInvoice(t *testing.T) {
	srv, _ := serviceSetup()

	t.Run("successfully stores the invoice", func(t *testing.T) {
		customer := "John Doe"
		inv, err := srv.CreateInvoice(customer)
		if err != nil {
			t.Fatalf("CreateInvoice(%q) failed: %v", customer, err)
		}

		if inv.ID == "" {
			t.Error("invoice.ID should not be empty")
		}

		if inv.Date != nil {
			t.Error("invalid invoice.Date, want nil")
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
			t.Errorf("invoice.CreatedAt %v is not equal to invoice.UpdatedAt %v", inv.CreatedAt, inv.UpdatedAt)
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
	srv, invoiceAPI := serviceSetup()

	t.Run("returns nil when no invoice is found in data storage", func(t *testing.T) {
		t.Skip()
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

	t.Run("propagates data storage errors", func(t *testing.T) {
		e := errors.New("storage failed to find invoice")
		strg := mocks.NewStorage(mocks.WithFindInvoiceError(e))
		srv := invoice.New(strg)

		invID := uuid.Nil.String()
		_, err := srv.ViewInvoice(invID)
		if err == nil {
			t.Fatalf("expected ViewInvoice(%q) to fail due to storage error", invID)
		}
		if got, want := err.Error(), fmt.Sprintf("find invoice %q failed: %s", invID, e.Error()); got != want {
			t.Errorf("ViewInvoice(%q) failed with: %s, want %s", invID, got, want)
		}
	})
}

func TestUpdateInvoiceCustomer(t *testing.T) {
	srv, invoiceAPI := serviceSetup()

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

	t.Run("fails when data storage error occurred - due to invoice search failure", func(t *testing.T) {
		e := errors.New("storage failed to find invoice")
		strg := mocks.NewStorage(mocks.WithFindInvoiceError(e))
		srv := invoice.New(strg)

		invID := uuid.Nil.String()
		customer := "John Doe"
		err := srv.UpdateInvoiceCustomer(invID, customer)
		if err == nil {
			t.Fatalf("expected UpdateInvoiceCustomer(%q, %q) to fail due to storage error", invID, customer)
		}
		if got, want := err.Error(), fmt.Sprintf("find invoice %q failed: %s", invID, e.Error()); got != want {
			t.Errorf("UpdateInvoiceCustomer(%q, %q) failed with: %s, want %s", invID, customer, got, want)
		}
	})

	t.Run("fails when data storage error occurred - due to invoice update failure", func(t *testing.T) {
		e := errors.New("storage failed to update invoice")
		inv := invoice.NewInvoice(uuid.NewString(), "John Doe")
		strg := mocks.NewStorage(
			mocks.WithFoundInvoice(&inv),
			mocks.WithUpdateInvoiceError(e))
		srv := invoice.New(strg)

		customer := "John Wick"
		err := srv.UpdateInvoiceCustomer(inv.ID, customer)
		if err == nil {
			t.Fatalf("expected UpdateInvoiceCustomer(%q, %q) to fail due to storage error", inv.ID, customer)
		}
		if got, want := err.Error(), fmt.Sprintf("update invoice %q failed: %s", inv.ID, e.Error()); got != want {
			t.Errorf("UpdateInvoiceCustomer(%q, %q) failed with: %s, want %s", inv.ID, customer, got, want)
		}
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
			t.Errorf("invalid invoice.UpdatedAt %v, want it to be after %v", vinv.UpdatedAt, inv.UpdatedAt)
		}
	})
}

func TestAddInvoiceItem(t *testing.T) {
	srv, invoiceAPI := serviceSetup()

	t.Run("fails when no invoice found", func(t *testing.T) {
		invID := uuid.Nil.String()
		_, err := srv.AddInvoiceItem(invID, "Pen", 123, 2)
		if err == nil {
			t.Fatalf("expected AddInvoiceItems(%q) to fail when invoice does not exist", invID)
		}
		if got, want := err.Error(), fmt.Sprintf("invoice %q not found", invID); got != want {
			t.Errorf("AddInvoiceItems(%q) failed with: %s, want %s", invID, got, want)
		}
	})

	t.Run("fails when item details not valid", func(t *testing.T) {
		invID, productName, price, qty := uuid.Nil.String(), "", 0, 0
		_, err := srv.AddInvoiceItem(invID, productName, price, qty)
		if err == nil {
			t.Fatalf("expected AddInvoiceItems(%q) to fail when item details not valid", invID)
		}
		want := "item details not valid: product name cannot be blank, price should be positive, qty should be positive"
		if got := err.Error(); got != want {
			t.Errorf("AddInvoiceItems(%q) failed with: %s, want %s", invID, got, want)
		}
	})

	t.Run("fails when invoice is in the status other than open", func(t *testing.T) {
		statuses := []invoice.Status{invoice.Issued, invoice.Paid, invoice.Canceled}
		invoices, err := invoiceAPI.CreateInvoicesWithStatuses(statuses...)
		if err != nil {
			t.Fatalf("invoiceAPI.CreateInvoicesWithStatuses() failed: %v", err)
		}

		for _, inv := range invoices {
			_, err := srv.AddInvoiceItem(inv.ID, "Pen", 123, 2)
			if err == nil {
				t.Fatalf("expected AddInvoiceItems(%q) to fail when invoice status is %q",
					inv.ID, inv.Status)
			}
			got := err.Error()
			want := fmt.Sprintf("item cannot be added to %q invoice", inv.Status)
			if got != want {
				t.Errorf("AddInvoiceItems(%q) failed with: %s, want %s", inv.ID, got, want)
			}
		}
	})

	t.Run("fails when data storage error occurred - due to invoice search failure", func(t *testing.T) {
		e := errors.New("storage failed to find invoice")
		strg := mocks.NewStorage(mocks.WithFindInvoiceError(e))
		srv := invoice.New(strg)

		invID := uuid.Nil.String()
		_, err := srv.AddInvoiceItem(invID, "Pen", 123, 2)
		if err == nil {
			t.Fatalf("expected AddInvoiceItems(%q) to fail due to storage error", invID)
		}
		if got, want := err.Error(), fmt.Sprintf("find invoice %q failed: %s", invID, e.Error()); got != want {
			t.Errorf("AddInvoiceItems(%q) failed with: %s, want %s", invID, got, want)
		}
	})

	t.Run("fails when data storage error occurred - due to invoice update failure", func(t *testing.T) {
		e := errors.New("storage failed to update invoice")
		inv := invoice.NewInvoice(uuid.NewString(), "John Doe")
		strg := mocks.NewStorage(
			mocks.WithFoundInvoice(&inv),
			mocks.WithUpdateInvoiceError(e))
		srv := invoice.New(strg)

		_, err := srv.AddInvoiceItem(inv.ID, "Pen", 123, 2)
		if err == nil {
			t.Fatalf("expected AddInvoiceItems(%q) to fail due to storage error", inv.ID)
		}
		if got, want := err.Error(), fmt.Sprintf("update invoice %q failed: %s", inv.ID, e.Error()); got != want {
			t.Errorf("AddInvoiceItems(%q) failed with: %s, want %s", inv.ID, got, want)
		}
	})

	t.Run("successfully adds invoice item", func(t *testing.T) {
		// place open invoice
		inv, err := invoiceAPI.CreateInvoice()
		if err != nil {
			t.Fatalf("invoiceAPI.CreateInvoice() failed: %v", err)
		}

		nitems := len(inv.Items)

		// add item
		productName, price, qty := "Pen", 123, 2
		item, err := srv.AddInvoiceItem(inv.ID, productName, price, qty)
		if err != nil {
			t.Fatalf("AddInvoiceItems(%q) failed: %v", inv.ID, err)
		}

		if item.ID == "" {
			t.Error("item.ID should not be empty")
		}

		if item.ProductName != productName {
			t.Errorf("invalid item.ProductName %q, want %q", item.ProductName, productName)
		}

		if item.Price != price {
			t.Errorf("invalid item.Price %d, want %d", item.Price, price)
		}

		if item.Qty != qty {
			t.Errorf("invalid item.Qty %d, want %d", item.Qty, qty)
		}

		if !item.CreatedAt.After(inv.CreatedAt) {
			t.Errorf("item.CreatedAt %v should be after invoice.CreatedAt %v", item.CreatedAt, inv.CreatedAt)
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
			t.Errorf("invalid invoice.UpdatedAt %v, want it to be after %v", vinv.UpdatedAt, inv.UpdatedAt)
		}
	})
}

func TestDeleteInvoiceItem(t *testing.T) {
	srv, invoiceAPI := serviceSetup()

	t.Run("fails when no invoice found", func(t *testing.T) {
		invID, itemID := uuid.Nil.String(), uuid.Nil.String()
		err := srv.DeleteInvoiceItem(invID, itemID)
		if err == nil {
			t.Fatalf("expected DeleteInvoiceItem(%q, %q) to fail when invoice does not exist", invID, itemID)
		}
		if got, want := err.Error(), fmt.Sprintf("invoice %q not found", invID); got != want {
			t.Errorf("DeleteInvoiceItem(%q, %q) failed with: %s, want %s", invID, itemID, got, want)
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

	t.Run("fails when data storage error occurred - due to invoice search failure", func(t *testing.T) {
		e := errors.New("storage failed to find invoice")
		strg := mocks.NewStorage(mocks.WithFindInvoiceError(e))
		srv := invoice.New(strg)

		invID, itemID := uuid.Nil.String(), uuid.Nil.String()
		err := srv.DeleteInvoiceItem(invID, itemID)
		if err == nil {
			t.Fatalf("expected DeleteInvoiceItem(%q, %q) to fail due to storage error", invID, itemID)
		}
		if got, want := err.Error(), fmt.Sprintf("find invoice %q failed: %s", invID, e.Error()); got != want {
			t.Errorf("DeleteInvoiceItem(%q, %q) failed with: %s, want %s", invID, itemID, got, want)
		}
	})

	t.Run("fails when data storage error occurred - due to invoice update failure", func(t *testing.T) {
		e := errors.New("storage failed to update invoice")

		nitems := 2
		inv, _ := invoiceAPI.CreateInvoiceWithNItems(nitems)
		itemID := inv.Items[0].ID

		strg := mocks.NewStorage(
			mocks.WithFoundInvoice(&inv),
			mocks.WithUpdateInvoiceError(e))
		srv := invoice.New(strg)

		err := srv.DeleteInvoiceItem(inv.ID, itemID)
		if err == nil {
			t.Fatalf("expected DeleteInvoiceItem(%q, %q) to fail due to storage error", inv.ID, itemID)
		}
		if got, want := err.Error(), fmt.Sprintf("update invoice %q failed: %s", inv.ID, e.Error()); got != want {
			t.Errorf("DeleteInvoiceItem(%q, %q) failed with: %s, want %s", inv.ID, itemID, got, want)
		}
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
			t.Errorf("invalid invoice.UpdatedAt %v, want it to be after %v", vinv.UpdatedAt, inv.UpdatedAt)
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
	srv, invoiceAPI := serviceSetup()

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

	t.Run("fails when data storage error occurred - due to invoice search failure", func(t *testing.T) {
		e := errors.New("storage failed to find invoice")
		strg := mocks.NewStorage(mocks.WithFindInvoiceError(e))
		srv := invoice.New(strg)

		invID := uuid.Nil.String()
		err := srv.IssueInvoice(invID)
		if err == nil {
			t.Fatalf("expected IssueInvoice(%q) to fail due to storage error", invID)
		}
		if got, want := err.Error(), fmt.Sprintf("find invoice %q failed: %s", invID, e.Error()); got != want {
			t.Errorf("IssueInvoice(%q) failed with: %s, want %s", invID, got, want)
		}
	})

	t.Run("fails when data storage error occurred - due to invoice update failure", func(t *testing.T) {
		e := errors.New("storage failed to update invoice")
		inv := invoice.NewInvoice(uuid.NewString(), "John Doe")
		strg := mocks.NewStorage(
			mocks.WithFoundInvoice(&inv),
			mocks.WithUpdateInvoiceError(e))
		srv := invoice.New(strg)

		err := srv.IssueInvoice(inv.ID)
		if err == nil {
			t.Fatalf("expected IssueInvoice(%q) to fail due to storage error", inv.ID)
		}
		if got, want := err.Error(), fmt.Sprintf("update invoice %q failed: %s", inv.ID, e.Error()); got != want {
			t.Errorf("IssueInvoice(%q) failed with: %s, want %s", inv.ID, got, want)
		}
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
			t.Errorf("invalid invoice.UpdatedAt %v, want it to be after %v", vinv.UpdatedAt, inv.UpdatedAt)
		}
	})
}

func TestPayInvoice(t *testing.T) {
	srv, invoiceAPI := serviceSetup()

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

	t.Run("fails when data storage error occurred - due to invoice search failure", func(t *testing.T) {
		e := errors.New("storage failed to find invoice")
		strg := mocks.NewStorage(mocks.WithFindInvoiceError(e))
		srv := invoice.New(strg)

		invID := uuid.Nil.String()
		err := srv.PayInvoice(invID)
		if err == nil {
			t.Fatalf("expected PayInvoice(%q) to fail due to storage error", invID)
		}
		if got, want := err.Error(), fmt.Sprintf("find invoice %q failed: %s", invID, e.Error()); got != want {
			t.Errorf("PayInvoice(%q) failed with: %s, want %s", invID, got, want)
		}
	})

	t.Run("fails when data storage error occurred - due to invoice update failure", func(t *testing.T) {
		e := errors.New("storage failed to update invoice")

		// using invoiceAPI to generate invoice in required state
		inv, _ := invoiceAPI.CreateInvoice(testapi.WithStatus(invoice.Issued))

		strg := mocks.NewStorage(
			mocks.WithFoundInvoice(&inv),
			mocks.WithUpdateInvoiceError(e))
		srv := invoice.New(strg)

		err := srv.PayInvoice(inv.ID)
		if err == nil {
			t.Fatalf("expected PayInvoice(%q) to fail due to storage error", inv.ID)
		}
		if got, want := err.Error(), fmt.Sprintf("update invoice %q failed: %s", inv.ID, e.Error()); got != want {
			t.Errorf("PayInvoice(%q) failed with: %s, want %s", inv.ID, got, want)
		}
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
			t.Errorf("invalid invoice.UpdatedAt %v, want it to be after %v", vinv.UpdatedAt, inv.UpdatedAt)
		}
	})
}

func TestCancelInvoice(t *testing.T) {
	srv, invoiceAPI := serviceSetup()

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

	t.Run("fails when data storage error occurred - due to invoice search failure", func(t *testing.T) {
		e := errors.New("storage failed to find invoice")
		strg := mocks.NewStorage(mocks.WithFindInvoiceError(e))
		srv := invoice.New(strg)

		invID := uuid.Nil.String()
		err := srv.CancelInvoice(invID)
		if err == nil {
			t.Fatalf("expected CancelInvoice(%q) to fail due to storage error", invID)
		}
		if got, want := err.Error(), fmt.Sprintf("find invoice %q failed: %s", invID, e.Error()); got != want {
			t.Errorf("CancelInvoice(%q) failed with: %s, want %s", invID, got, want)
		}
	})

	t.Run("fails when data storage error occurred - due to invoice update failure", func(t *testing.T) {
		e := errors.New("storage failed to update invoice")
		inv := invoice.NewInvoice(uuid.NewString(), "John Doe")
		strg := mocks.NewStorage(
			mocks.WithFoundInvoice(&inv),
			mocks.WithUpdateInvoiceError(e))
		srv := invoice.New(strg)

		err := srv.CancelInvoice(inv.ID)
		if err == nil {
			t.Fatalf("expected CancelInvoice(%q) to fail due to storage error", inv.ID)
		}
		if got, want := err.Error(), fmt.Sprintf("update invoice %q failed: %s", inv.ID, e.Error()); got != want {
			t.Errorf("CancelInvoice(%q) failed with: %s, want %s", inv.ID, got, want)
		}
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
			t.Errorf("invalid invoice.UpdatedAt %v, want it to be after %v", vinv.UpdatedAt, inv.UpdatedAt)
		}
	})
}
