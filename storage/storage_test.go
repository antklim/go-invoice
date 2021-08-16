package storage_test

import (
	"fmt"
	"testing"

	"github.com/antklim/go-invoice/storage"
)

func TestFactoryError(t *testing.T) {
	kind := "foo"
	strg, err := storage.Factory(kind)
	if strg != nil {
		t.Errorf("Factory(%q) returned not nil storage", kind)
	}
	if got, want := err.Error(), fmt.Sprintf("unknown storage %q", kind); got != want {
		t.Errorf("Factory(%q) unexpected error %v, want %v", kind, got, want)
	}
}

func TestFactory(t *testing.T) {
	testCases := []struct {
		desc string
		kind string
	}{
		{
			desc: "returns in memory storage",
			kind: "memory",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			strg, err := storage.Factory(tC.kind)
			if err != nil {
				t.Errorf("Factory(%q) failed: %v", tC.kind, err)
			}
			if strg == nil {
				t.Errorf("Factory(%q) storage is nil", tC.kind)
			}
			if got, want := fmt.Sprintf("%T", strg), "*memory.Memory"; got != want {
				t.Errorf("Factory(%q) storage type is %q, want %q", tC.kind, got, want)
			}
		})
	}
}
