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

func NewStorage(opts ...StorageOptions) *Storage {
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

type StorageOptions interface {
	apply(*Storage)
}

type funcStorageOptions struct {
	f func(strg *Storage)
}

func (fso *funcStorageOptions) apply(strg *Storage) {
	fso.f(strg)
}

func newFuncStorageOptions(f func(strg *Storage)) StorageOptions {
	return &funcStorageOptions{f: f}
}

func WithAddInvoiceError(err error) StorageOptions {
	return newFuncStorageOptions(func(strg *Storage) {
		strg.errors[addInvoice] = err
	})
}

func WithFindInvoiceError(err error) StorageOptions {
	return newFuncStorageOptions(func(strg *Storage) {
		strg.errors[findInvoice] = err
	})
}

func WithUpdateInvoiceError(err error) StorageOptions {
	return newFuncStorageOptions(func(strg *Storage) {
		strg.errors[updateInvoice] = err
	})
}

func WithFoundInvoice(inv *invoice.Invoice) StorageOptions {
	return newFuncStorageOptions(func(strg *Storage) {
		strg.foundInvoice = inv
	})
}
