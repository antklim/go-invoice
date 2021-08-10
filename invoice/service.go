package invoice

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Service struct{}

func (s Service) CreateInvoice(customerName string) (Invoice, error) {
	invID := uuid.NewString()
	return createInvoice(invID, customerName), nil
}

func (s Service) ViewInvoice(id string) (Invoice, error) {
	return Invoice{}, errors.New("not implemented")
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
