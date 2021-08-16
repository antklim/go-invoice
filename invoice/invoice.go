package invoice

import "time"

type Status int

// Supported invoice statuses
const (
	Open Status = iota
	Issued
	Closed
)

type Invoice struct {
	ID           string
	CustomerName string
	Date         *time.Time // issue date
	Status       Status
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (inv *Invoice) Equal(other Invoice) bool {
	var invDatesEqual bool
	if inv.Date == nil && other.Date == nil {
		invDatesEqual = true
	} else if inv.Date != nil && other.Date != nil {
		invDatesEqual = inv.Date.Equal(*other.Date)
	}

	return inv.ID == other.ID &&
		inv.CustomerName == other.CustomerName &&
		invDatesEqual &&
		inv.Status == other.Status &&
		inv.CreatedAt.Equal(other.CreatedAt) &&
		inv.UpdatedAt.Equal(other.UpdatedAt)
}

func NewInvoice(id, customer string) Invoice {
	now := time.Now()
	return Invoice{
		ID:           id,
		CustomerName: customer,
		Status:       Open,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
