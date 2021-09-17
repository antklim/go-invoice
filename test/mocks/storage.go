package mocks

import (
	"github.com/antklim/go-invoice/invoice"
)

type memoryOp int

const (
	addInvoice memoryOp = iota
	findInvoice
	updateInvoice
)

// Storage describes storage mock.
type Storage struct {
	errors       map[memoryOp]error
	foundInvoice *invoice.Invoice
}

func NewStorage(opts ...StorageOption) *Storage {
	strg := &Storage{
		errors: make(map[memoryOp]error),
	}

	for _, o := range opts {
		o.apply(strg)
	}

	return strg
}

func (strg *Storage) AddInvoice(inv invoice.Invoice) error {
	return strg.errors[addInvoice]
}

func (strg *Storage) FindInvoice(id string) (*invoice.Invoice, error) {
	if err := strg.errors[findInvoice]; err != nil {
		return nil, err
	}

	return strg.foundInvoice, nil
}

func (strg *Storage) UpdateInvoice(inv invoice.Invoice) error {
	return strg.errors[updateInvoice]
}

var _ invoice.Storage = (*Storage)(nil)

type StorageOption interface {
	apply(*Storage)
}

type funcStorageOption struct {
	f func(strg *Storage)
}

func (fso *funcStorageOption) apply(strg *Storage) {
	fso.f(strg)
}

func newFuncStorageOption(f func(strg *Storage)) StorageOption {
	return &funcStorageOption{f: f}
}

func WithAddInvoiceError(err error) StorageOption {
	return newFuncStorageOption(func(strg *Storage) {
		strg.errors[addInvoice] = err
	})
}

func WithFindInvoiceError(err error) StorageOption {
	return newFuncStorageOption(func(strg *Storage) {
		strg.errors[findInvoice] = err
	})
}

func WithUpdateInvoiceError(err error) StorageOption {
	return newFuncStorageOption(func(strg *Storage) {
		strg.errors[updateInvoice] = err
	})
}

func WithFoundInvoice(inv *invoice.Invoice) StorageOption {
	return newFuncStorageOption(func(strg *Storage) {
		strg.foundInvoice = inv
	})
}
