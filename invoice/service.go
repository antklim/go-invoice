package invoice

import (
	"time"

	"github.com/google/uuid"
)

type Service struct{}

func (s Service) CreateInvoice(customerName string, date time.Time) (Invoice, error) {
	invID := uuid.NewString()
	inv := createInvoice(invID, customerName, date)
	return inv, nil
}
