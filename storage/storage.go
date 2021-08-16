package storage

import (
	"fmt"

	"github.com/antklim/go-invoice/invoice"
	"github.com/antklim/go-invoice/storage/memory"
)

func Factory(kind string) (invoice.Storage, error) {
	switch kind {
	case "memory":
		return memory.New(), nil
	default:
		return nil, fmt.Errorf("unknown storage %q", kind)
	}
}
