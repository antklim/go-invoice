package invoice

import "time"

type Status int

const (
	Open Status = iota
	Issued
	Closed
)

var statuses = [...]string{Open: "open", Issued: "issued", Closed: "closed"}

type Invoice struct {
	ID           string
	CustomerName string
	Date         *time.Time // issue date
	Status       Status
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func createInvoice(id, customer string) Invoice {
	now := time.Now()
	return Invoice{
		ID:           id,
		CustomerName: customer,
		Status:       Open,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
