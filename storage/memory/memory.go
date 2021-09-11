package memory

import (
	"fmt"
	"sync"
	"time"

	"github.com/antklim/go-invoice/invoice"
)

type Memory struct {
	sync.RWMutex // guards records
	records      map[string]invoice.Invoice
}

var _ invoice.Storage = (*Memory)(nil)

func New() *Memory {
	return &Memory{records: make(map[string]invoice.Invoice)}
}

func (memo *Memory) AddInvoice(inv invoice.Invoice) error {
	memo.Lock()
	defer memo.Unlock()
	if _, ok := memo.records[inv.ID]; ok {
		return fmt.Errorf("ID %q exists", inv.ID)
	}
	memo.records[inv.ID] = inv
	return nil
}

func (memo *Memory) FindInvoice(id string) (*invoice.Invoice, error) {
	memo.RLock()
	defer memo.RUnlock()

	inv, ok := memo.records[id]
	if !ok {
		return nil, nil
	}

	return &inv, nil
}

func (memo *Memory) UpdateInvoice(inv invoice.Invoice) error {
	memo.Lock()
	defer memo.Unlock()

	r, ok := memo.records[inv.ID]
	if !ok {
		return fmt.Errorf("invoice %q not found", inv.ID)
	}

	inv.CreatedAt = r.CreatedAt
	inv.UpdatedAt = time.Now()
	memo.records[inv.ID] = inv

	return nil
}
