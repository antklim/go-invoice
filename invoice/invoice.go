package invoice

import "time"

const (
	invoiceOpen   = "open"
	invoiceIssued = "issued"
	invoiceClosed = "closed"
)

type Invoice struct {
	ID           string
	CustomerName string
	Date         *time.Time // issue date
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func createInvoice(id, customer string) Invoice {
	now := time.Now()
	return Invoice{
		ID:           id,
		CustomerName: customer,
		Status:       invoiceOpen,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
