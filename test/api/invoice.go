package api

import (
	"time"

	"github.com/antklim/go-invoice/invoice"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Invoice struct {
	strg invoice.Storage
}

func NewIvoiceAPI(strg invoice.Storage) *Invoice {
	return &Invoice{strg: strg}
}

// CreateInvoice generates and stores invoice. It returns created invoice or any
// error occurred. In error case an empty invoice returned.
//
// By default it generates invoice with random ID in uuid format, customer name
// is "John Doe", empty issue date, open status, created at and updated at dates
// set to the current time value when the invoice was generated. It is possible
// to provide predefined values for any invoice field.
//
// For example:
//
//	// Create default invoice.
//	invapi.CreateInvoice()
//
//	// Create paid invoice for John Wick.
//	invapi.CreateInvoice(
//		invapi.WithCustomerName("John Wick"),
//		invapi.WithStatus(invoice.Paid))
//
func (api *Invoice) CreateInvoice(opts ...InvoiceOptions) (invoice.Invoice, error) {
	inv := defaultInvoice()

	for _, o := range opts {
		o.apply(&inv)
	}

	if err := api.strg.AddInvoice(inv); err != nil {
		return invoice.Invoice{}, errors.Wrap(err, "add invoice failed")
	}
	return inv, nil
}

// CreateInvoicesWithStatuses generates and stores a collection of invoices.
// Every invoice generated according to provided status.
func (api *Invoice) CreateInvoicesWithStatuses(statuses ...invoice.Status) ([]invoice.Invoice, error) {
	invoices := make([]invoice.Invoice, 0, len(statuses))
	for _, status := range statuses {
		inv, err := api.CreateInvoice(WithStatus(status))
		if err != nil {
			return nil, err
		}
		invoices = append(invoices, inv)
	}
	return invoices, nil
}

// CreateInvoiceWithNItems generates and stores invoice with n items. It returns
// created invoice or any error occurred. In error case an empty invoice
// returned.
func (api *Invoice) CreateInvoiceWithNItems(n int, opts ...InvoiceOptions) (invoice.Invoice, error) {
	items := make([]invoice.Item, 0, n)
	for i := 0; i < n; i++ {
		items = append(items, api.ItemFactory())
	}

	opts = append(opts, WithItems(items...))

	return api.CreateInvoice(opts...)
}

// ItemFactory generates invoice items with default product name, price and
// quantity. Every method call generates unique item ID and created at time.
func (api *Invoice) ItemFactory() invoice.Item {
	return defaultItem()
}

func defaultInvoice() invoice.Invoice {
	id := uuid.NewString()
	now := time.Now()

	return invoice.Invoice{
		ID:           id,
		CustomerName: "John Doe",
		Date:         nil,
		Status:       invoice.Open,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

type InvoiceOptions interface {
	apply(*invoice.Invoice)
}

type funcInvoiceOptions struct {
	f func(*invoice.Invoice)
}

func (fio *funcInvoiceOptions) apply(inv *invoice.Invoice) {
	fio.f(inv)
}

func newFuncInvoiceOptions(f func(*invoice.Invoice)) InvoiceOptions {
	return &funcInvoiceOptions{f: f}
}

func WithID(id string) InvoiceOptions {
	return newFuncInvoiceOptions(func(inv *invoice.Invoice) {
		inv.ID = id
	})
}

func WithCustomerName(cn string) InvoiceOptions {
	return newFuncInvoiceOptions(func(inv *invoice.Invoice) {
		inv.CustomerName = cn
	})
}

func WithIssueaDate(date *time.Time) InvoiceOptions {
	return newFuncInvoiceOptions(func(inv *invoice.Invoice) {
		inv.Date = date
	})
}

func WithStatus(status invoice.Status) InvoiceOptions {
	return newFuncInvoiceOptions(func(inv *invoice.Invoice) {
		inv.Status = status
	})
}

func WithItems(items ...invoice.Item) InvoiceOptions {
	return newFuncInvoiceOptions(func(inv *invoice.Invoice) {
		inv.Items = items
	})
}

func WithCreatedAt(date time.Time) InvoiceOptions {
	return newFuncInvoiceOptions(func(inv *invoice.Invoice) {
		inv.CreatedAt = date
	})
}

func WithUpdatedAt(date time.Time) InvoiceOptions {
	return newFuncInvoiceOptions(func(inv *invoice.Invoice) {
		inv.UpdatedAt = date
	})
}

func defaultItem() invoice.Item {
	id := uuid.NewString()
	now := time.Now()

	return invoice.Item{
		ID:          id,
		ProductName: "Pen",
		Price:       123, // nolint:gomnd
		Qty:         2,   // nolint:gomnd
		CreatedAt:   now,
	}
}
