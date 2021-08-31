package invoice

import (
	"sort"
	"time"
)

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

// FormatStatus returns formatted status.
func FormatStatus(st Status) string {
	return statuses[st]
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

func (inv *Invoice) Equal(other *Invoice) bool {
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
		inv.itemsEqual(other.Items) &&
		inv.CreatedAt.Equal(other.CreatedAt) &&
		inv.UpdatedAt.Equal(other.UpdatedAt)
}

// FindItemIndex returns the index of the first item in the collection that
// satisfies the provided testing function. Testing function should returns true
// to indicate that the satisfying item was found.
func (inv *Invoice) FindItemIndex(f func(item Item) bool) int {
	for i, item := range inv.Items {
		if f(item) {
			return i
		}
	}
	return -1
}

// ContainsItem returns true when invoice contains item with the provided ID.
func (inv *Invoice) ContainsItem(id string) bool {
	idx := inv.FindItemIndex(func(item Item) bool {
		return item.ID == id
	})
	return idx != -1
}

// AddItem adds an item to the invoice.
func (inv *Invoice) AddItem(item Item) {
	inv.Items = append(inv.Items, item)
}

// DeleteItem deletes an item by ID. This operation is idempotent, repeatable
// item delete supported. Returns true when the item was found and deleted from
// the items collection.
func (inv *Invoice) DeleteItem(id string) bool {
	idx := inv.FindItemIndex(func(item Item) bool {
		return item.ID == id
	})

	if idx == -1 {
		return false
	}

	inv.Items = append(inv.Items[:idx], inv.Items[idx+1:]...)
	return true
}

// Issue sets invoice to issued state.
func (inv *Invoice) Issue() {
	inv.Status = Issued
	now := time.Now()
	inv.Date = &now
}

// Pay sets invoice to paid state.
func (inv *Invoice) Pay() {
	inv.Status = Paid
}

// Pay sets invoice to canceled state.
func (inv *Invoice) Cancel() {
	inv.Status = Canceled
}

func (inv *Invoice) itemsEqual(otherItems []Item) bool {
	if len(inv.Items) != len(otherItems) {
		return false
	}

	sort.Sort(byItemID(inv.Items))
	sort.Sort(byItemID(otherItems))

	for i, item := range inv.Items {
		item := item
		if !otherItems[i].Equal(&item) {
			return false
		}
	}

	return true
}

type Item struct {
	ID          string
	ProductName string
	Price       uint // price in cents
	Qty         uint
	CreatedAt   time.Time
}

func (item *Item) Equal(other *Item) bool {
	return item.ID == other.ID &&
		item.ProductName == other.ProductName &&
		item.Price == other.Price &&
		item.Qty == other.Qty &&
		item.CreatedAt.Equal(other.CreatedAt)
}

type byItemID []Item

func (x byItemID) Len() int           { return len(x) }
func (x byItemID) Less(i, j int) bool { return x[i].ID < x[j].ID }
func (x byItemID) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
