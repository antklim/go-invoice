package invoice

import (
	"fmt"
	"sort"
	"strings"
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

var statusName = map[Status]string{
	Open:     "open",
	Issued:   "issued",
	Paid:     "paid",
	Canceled: "canceled",
}

func (s Status) String() string { return statusName[s] }

type Invoice struct {
	ID           string
	CustomerName string
	Date         *time.Time // issue date
	Status       Status
	Items        []Item
	CreatedAt    time.Time
	UpdatedAt    time.Time
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

// UpdateCustomerName sets new customer name. It returns error when invoice
// cannot be updated.
func (inv *Invoice) UpdateCustomerName(name string) error {
	if inv.Status != Open {
		return fmt.Errorf("%q invoice cannot be updated", inv.Status)
	}

	inv.CustomerName = name
	return nil
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

// AddItem adds an item to the invoice. It returns error when invoice item
// cannot be added.
func (inv *Invoice) AddItem(item Item) error {
	if inv.Status != Open {
		return fmt.Errorf("item cannot be added to %q invoice", inv.Status)
	}

	inv.Items = append(inv.Items, item)
	return nil
}

// DeleteItem deletes an item by ID. This operation is idempotent, repeatable
// item delete supported. Returns true when the item was found and deleted from
// the items collection.
func (inv *Invoice) DeleteItem(id string) (bool, error) {
	if inv.Status != Open {
		return false, fmt.Errorf("item cannot be deleted from %q invoice", inv.Status)
	}

	idx := inv.FindItemIndex(func(item Item) bool {
		return item.ID == id
	})

	if idx == -1 {
		return false, nil
	}

	inv.Items = append(inv.Items[:idx], inv.Items[idx+1:]...)
	return true, nil
}

// Issue sets invoice to issued state. It returns error when invoice is not
// issueable.
func (inv *Invoice) Issue() error {
	if inv.Status != Open {
		return fmt.Errorf("%q invoice cannot be issued", inv.Status)
	}

	inv.Status = Issued
	now := time.Now()
	inv.Date = &now
	return nil
}

// Pay sets invoice to paid state. It returns error when invoice is not payable.
func (inv *Invoice) Pay() error {
	if inv.Status != Issued {
		return fmt.Errorf("%q invoice cannot be paid", inv.Status)
	}

	inv.Status = Paid
	return nil
}

// Cancel sets invoice to canceled state. It returns error when invoice is not
// cancelable.
func (inv *Invoice) Cancel() error {
	if inv.Status == Canceled || inv.Status == Paid {
		return fmt.Errorf("%q invoice cannot be canceled", inv.Status)
	}

	inv.Status = Canceled
	return nil
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
	Price       int // price in cents
	Qty         int
	CreatedAt   time.Time
}

func (item *Item) Equal(other *Item) bool {
	return item.ID == other.ID &&
		item.ProductName == other.ProductName &&
		item.Price == other.Price &&
		item.Qty == other.Qty &&
		item.CreatedAt.Equal(other.CreatedAt)
}

func (item *Item) Validate() error {
	var errors []string

	if item.ProductName == "" {
		errors = append(errors, "product name cannot be blank")
	}

	if item.Price < 1 {
		errors = append(errors, "price should be positive")
	}

	if item.Qty < 1 {
		errors = append(errors, "qty should be positive")
	}

	if len(errors) == 0 {
		return nil
	}

	return fmt.Errorf("item details not valid: %s", strings.Join(errors, ", "))
}

type byItemID []Item

func (x byItemID) Len() int           { return len(x) }
func (x byItemID) Less(i, j int) bool { return x[i].ID < x[j].ID }
func (x byItemID) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
