package invoice

import "time"

type Status int

// Supported invoice statuses
const (
	Open Status = iota
	Issued
	Paid
	Canceled
)

var statuses = [...]string{
	Open:     "open",
	Issued:   "issued",
	Paid:     "paid",
	Canceled: "canceled",
}

type Invoice struct {
	ID           string
	CustomerName string
	Date         *time.Time // issue date
	Status       Status
	Items        []Item
	CreatedAt    time.Time
	UpdatedAt    time.Time
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

// FormatStatus returns formatted invoice status.
func (inv *Invoice) FormatStatus() string {
	return statuses[inv.Status]
}

func (inv *Invoice) IsClosed() bool {
	return inv.Status == Paid || inv.Status == Canceled
}

type Item struct {
	ID          string
	ProductName string
	Price       uint // price in cents
	Qty         uint
	CreatedAt   time.Time
}
