package memory

import (
	"fmt"
	"sync"

	"github.com/antklim/go-invoice/invoice"
)

type Memory struct {
	sync.RWMutex // guards records
	records      map[string]invoice.Invoice
}

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
