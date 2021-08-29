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

func (api *Invoice) CreateInvoice(opts ...InvoiceOptions) (string, error) {
	inv := defaultInvoice()

	for _, o := range opts {
		o.apply(&inv)
	}

	if err := api.strg.AddInvoice(inv); err != nil {
		return "", errors.Wrap(err, "add invoice failed")
	}
	return inv.ID, nil
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
