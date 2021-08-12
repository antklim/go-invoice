package invoice

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var (
	mu       sync.RWMutex // guards invoices
	invoices map[string]Invoice
)

type Service struct{}

func (s Service) CreateInvoice(customerName string) (Invoice, error) {
	invID := uuid.NewString()
	inv := createInvoice(invID, customerName)

	mu.Lock()
	defer mu.Unlock()
	if invoices == nil {
		invoices = make(map[string]Invoice)
	}
	if _, ok := invoices[inv.ID]; ok {
		return Invoice{}, fmt.Errorf("store invoice: ID %q exists", inv.ID)
	}
	invoices[inv.ID] = inv

	return inv, nil
}

func (s Service) ViewInvoice(id string) (Invoice, error) {
	mu.RLock()
	defer mu.RUnlock()
	if invoices == nil {
		return Invoice{}, nil
	}
	return invoices[id], nil
}

func (s Service) UpdateInvoiceCustomer(id, name string) error {
	return errors.New("not implemented")
}

func (s Service) CancelInvoice(id string) error {
	return errors.New("not implemented")
}

func (s Service) PayInvoice(id string) error {
	return errors.New("not implemented")
}
