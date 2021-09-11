package storage

import (
	"github.com/antklim/go-invoice/invoice"
	"github.com/antklim/go-invoice/storage/dynamo"
	"github.com/antklim/go-invoice/storage/memory"
)

type Memory struct{}

func (Memory) MakeStorage() invoice.Storage {
	return memory.New()
}

var _ invoice.StorageFactory = new(Memory)

type Dynamo struct {
	client dynamo.API
}

func (s *Dynamo) MakeStorage() invoice.Storage {
	return dynamo.New(s.client)
}

var _ invoice.StorageFactory = (*Dynamo)(nil)
