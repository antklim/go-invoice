package invoice

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Storage interface {
	AddInvoice(Invoice) error
	GetInvoice(string) (Invoice, error)
}

type Service struct {
	storage Storage
}

func NewService(storage Storage) Service {
	return Service{storage: storage}
}

func (s Service) CreateInvoice(customerName string) (Invoice, error) {
	invID := uuid.NewString()
	inv := createInvoice(invID, customerName)
	err := s.storage.AddInvoice(inv)
	if err != nil {
		return Invoice{}, errors.Wrap(err, "failed to store invoice")
	}
	return inv, nil
}

func (s Service) ViewInvoice(id string) (Invoice, error) {
	inv, err := s.storage.GetInvoice(id)
	if err != nil {
		return Invoice{}, errors.Wrap(err, "failed to get invoice")
	}
	return inv, nil
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
