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

// New initiates a new instance of the service.
func New(strg Storage) *Service {
	return &Service{strg: strg}
}

// CreateInvoice generates and stores an invoice. A new invoice generated with
// the provided customer name. Invoice and any occurred error returned.
func (s *Service) CreateInvoice(customerName string) (Invoice, error) {
	invID := uuid.NewString()
	inv := NewInvoice(invID, customerName)
	err := s.strg.AddInvoice(inv)
	return inv, err
}

// ViewInvoice finds an invoice by invoice ID. It returns non nil pointer to the
// found invoice or nil in case when no invoices selected by ID. Nil invoice
// pointer also returned in error case.
func (s *Service) ViewInvoice(id string) (*Invoice, error) {
	return s.strg.FindInvoice(id)
}

// UpdateInvoiceCustomer updates invoice's customer name. If invoice not found
// by provided ID or any issue occurred during invoice lookup or update an error
// returned. Only invoices in "open" status are allowed to be updated.
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

// AddInvoiceItem adds invoice item to the invoice. If invoice not found
// by provided ID or any issue occurred during invoice lookup or update an error
// returned. Only invoices in "open" status are allowed to be updated.
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

// DeleteInvoiceItem deletes invoice item to the invoice. If invoice not found
// by provided ID or any issue occurred during invoice lookup or update an error
// returned. Only invoices in "open" status are allowed to be updated.
func (s *Service) DeleteInvoiceItem(invID, itemID string) error {
	inv, err := s.strg.FindInvoice(invID)
	if err != nil {
		return errors.Wrapf(err, errFindFailed, invID)
	}
	if inv == nil {
		return fmt.Errorf(errNotFound, invID)
	}

	if inv.Status != Open {
		return fmt.Errorf("item cannot be deleted from %q invoice", inv.FormatStatus())
	}

	idx := inv.FindItemIndex(func(item Item) bool {
		return item.ID == itemID
	})

	if idx == -1 {
		return nil
	}

	inv.Items = append(inv.Items[:idx], inv.Items[idx+1:]...)
	if err := s.strg.UpdateInvoice(*inv); err != nil {
		return errors.Wrapf(err, errUpdateFailed, invID)
	}

	return nil
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
