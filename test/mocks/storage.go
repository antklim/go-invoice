package mocks

import (
	"sync"

	"github.com/antklim/go-invoice/invoice"
)

type operation int

const (
	addInvoice operation = iota
	findInvoice
	updateInvoice
)

// Storage describes storage mock.
type Storage struct {
	errors       map[operation]error
	foundInvoice *invoice.Invoice

	sync.RWMutex // guards calls
	calls        map[operation]int
}

func NewStorage(opts ...StorageOptions) *Storage {
	strg := &Storage{
		errors: make(map[operation]error),
		calls:  make(map[operation]int),
	}

	for _, o := range opts {
		o.apply(strg)
	}

	return strg
}

func (strg *Storage) AddInvoice(inv invoice.Invoice) error {
	strg.Lock()
	defer strg.Unlock()
	strg.calls[addInvoice]++

	return strg.errors[addInvoice]
}

func (strg *Storage) FindInvoice(id string) (*invoice.Invoice, error) {
	strg.Lock()
	defer strg.Unlock()
	strg.calls[findInvoice]++

	if err := strg.errors[findInvoice]; err != nil {
		return nil, err
	}

	return strg.foundInvoice, nil
}

func (strg *Storage) UpdateInvoice(inv invoice.Invoice) error {
	strg.Lock()
	defer strg.Unlock()
	strg.calls[updateInvoice]++

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
