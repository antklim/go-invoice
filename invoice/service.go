package invoice

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var (
	errCreateFailed = "create invoice failed"
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
	if err != nil {
		err = errors.Wrap(err, errCreateFailed)
	}
	return inv, err
}

// ViewInvoice finds an invoice by invoice ID. It returns non nil pointer to the
// found invoice or nil in case when no invoices selected by ID. Nil invoice
// pointer also returned in error case.
func (s *Service) ViewInvoice(id string) (*Invoice, error) {
	return s.findInvoice(id)
}

// UpdateInvoiceCustomer updates invoice's customer name. If invoice not found
// by provided ID or any issue occurred during invoice lookup or update an error
// returned. Only invoices in "open" status are allowed to be updated.
func (s *Service) UpdateInvoiceCustomer(id, name string) error {
	inv, err := s.mustFindInvoice(id)
	if err != nil {
		return err
	}

	if err := inv.UpdateCustomerName(name); err != nil {
		return err
	}

	if err := s.strg.UpdateInvoice(*inv); err != nil {
		return errors.Wrapf(err, errUpdateFailed, id)
	}

	return nil
}

// AddInvoiceItem adds invoice item to the invoice. If invoice not found
// by provided ID or any issue occurred during invoice lookup or update an error
// returned. Only invoices in "open" status are allowed to be updated.
func (s *Service) AddInvoiceItem(invID, productName string, price, qty int) (Item, error) {
	// TODO: add price and qty validation

	inv, err := s.mustFindInvoice(invID)
	if err != nil {
		return Item{}, err
	}

	item := NewItem(uuid.NewString(), productName, price, qty)
	if err := inv.AddItem(item); err != nil {
		return Item{}, err
	}

	if err := s.strg.UpdateInvoice(*inv); err != nil {
		return Item{}, errors.Wrapf(err, errUpdateFailed, invID)
	}

	return item, nil
}

// DeleteInvoiceItem deletes invoice item to the invoice. If invoice not found
// by provided ID or any issue occurred during invoice lookup or update an error
// returned. Only invoices in "open" status are allowed to be updated.
func (s *Service) DeleteInvoiceItem(invID, itemID string) error {
	inv, err := s.mustFindInvoice(invID)
	if err != nil {
		return err
	}

	ok, err := inv.DeleteItem(itemID)
	if err != nil {
		return err
	}

	if ok {
		// update storage only when item collection was changed
		if err := s.strg.UpdateInvoice(*inv); err != nil {
			return errors.Wrapf(err, errUpdateFailed, invID)
		}
	}

	return nil
}

// IssueInvoice sets invoice the the issued status. If invoice not found
// by provided ID or any issue occurred during invoice lookup or update an error
// returned. Only invoices in "open" status are allowed to be issued.
func (s *Service) IssueInvoice(id string) error {
	inv, err := s.mustFindInvoice(id)
	if err != nil {
		return err
	}

	if err := inv.Issue(); err != nil {
		return err
	}

	if err := s.strg.UpdateInvoice(*inv); err != nil {
		return errors.Wrapf(err, errUpdateFailed, inv.ID)
	}

	return nil
}

// CancelInvoice sets invoice to the canceled status. If invoice not found
// by provided ID or any issue occurred during invoice lookup or update an error
// returned. Canceled or paid invoices cannot be canceled.
func (s *Service) CancelInvoice(id string) error {
	inv, err := s.mustFindInvoice(id)
	if err != nil {
		return err
	}

	if err := inv.Cancel(); err != nil {
		return err
	}

	if err := s.strg.UpdateInvoice(*inv); err != nil {
		return errors.Wrapf(err, errUpdateFailed, inv.ID)
	}

	return nil
}

// PayInvoice sets invoice to the paid status. If invoice not found
// by provided ID or any issue occurred during invoice lookup or update an error
// returned. Only invoices in "issued" status are allowed to be paid.
func (s *Service) PayInvoice(id string) error {
	inv, err := s.mustFindInvoice(id)
	if err != nil {
		return err
	}

	if err := inv.Pay(); err != nil {
		return err
	}

	if err := s.strg.UpdateInvoice(*inv); err != nil {
		return errors.Wrapf(err, errUpdateFailed, inv.ID)
	}

	return nil
}

// mustFindInvoice searches for the invoice by id. If invoice not found or other
// issues occurred during invoice lookup an error returned. It returns a non-nil
// pointer to the found invoice.
func (s *Service) mustFindInvoice(id string) (*Invoice, error) {
	inv, err := s.findInvoice(id)
	if err != nil {
		return nil, err
	}
	if inv == nil {
		return nil, fmt.Errorf(errNotFound, id)
	}
	return inv, nil
}

// findInvoice searches for the invoice by id. It returns error only when
// storage failure occurred during invoice lookup.
//
// Invoice pointer is nil in error case or when invoice not found. Otherwise a
// non-nil pointer to the found invoice returned.
func (s *Service) findInvoice(id string) (*Invoice, error) {
	inv, err := s.strg.FindInvoice(id)
	if err != nil {
		return nil, errors.Wrapf(err, errFindFailed, id)
	}
	return inv, nil
}
