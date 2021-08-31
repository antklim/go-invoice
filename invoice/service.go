package invoice

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// TODO: add error interfaces

var (
	errFindFailed   = "find invoice %q failed"
	errUpdateFailed = "update invoice %q failed"
	errNotFound     = "invoice %q not found"
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
	inv, err := s.strg.FindInvoice(id)
	if err != nil {
		return errors.Wrapf(err, errFindFailed, id)
	}
	if inv == nil {
		return fmt.Errorf(errNotFound, id)
	}

	if inv.Status != Open {
		return fmt.Errorf("%q invoice cannot be updated", inv.FormatStatus())
	}

	inv.CustomerName = name
	if err := s.strg.UpdateInvoice(*inv); err != nil {
		return errors.Wrapf(err, errUpdateFailed, id)
	}

	return nil
}

func (s *Service) AddInvoiceItem(id string, item Item) error {
	inv, err := s.strg.FindInvoice(id)
	if err != nil {
		return errors.Wrapf(err, errFindFailed, id)
	}
	if inv == nil {
		return fmt.Errorf(errNotFound, id)
	}

	if inv.Status != Open {
		return fmt.Errorf("item cannot be added to %q invoice", inv.FormatStatus())
	}

	inv.Items = append(inv.Items, item)
	if err := s.strg.UpdateInvoice(*inv); err != nil {
		return errors.Wrapf(err, errUpdateFailed, id)
	}

	return nil
}

func (s *Service) DeleteInvoiceItem() error {
	return errors.New("not implemented")
}

func (s *Service) IssueInvoice(id string) error {
	return errors.New("not implemented")
}

func (s *Service) CancelInvoice(id string) error {
	return errors.New("not implemented")
}

func (s *Service) PayInvoice(id string) error {
	return errors.New("not implemented")
}
