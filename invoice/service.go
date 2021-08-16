package invoice

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Service struct {
	strg Storage
}

func New(strg Storage) *Service {
	return &Service{strg: strg}
}

func (s *Service) CreateInvoice(customerName string) (Invoice, error) {
	invID := uuid.NewString()
	inv := NewInvoice(invID, customerName)
	err := s.strg.AddInvoice(inv)
	return inv, err
}

func (s *Service) ViewInvoice(id string) (*Invoice, error) {
	return s.strg.FindInvoice(id)
}

func (s *Service) UpdateInvoiceCustomer(id, name string) error {
	return errors.New("not implemented")
}

func (s *Service) CancelInvoice(id string) error {
	return errors.New("not implemented")
}

func (s *Service) PayInvoice(id string) error {
	return errors.New("not implemented")
}
