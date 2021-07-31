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
	Date         time.Time // issue date
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func createInvoice(id, customer string, date time.Time) Invoice {
	now := time.Now()
	return Invoice{
		ID:           id,
		CustomerName: customer,
		Date:         date,
		Status:       invoiceOpen,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
