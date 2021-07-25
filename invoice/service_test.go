package invoice_test

import (
	"errors"
	"testing"
	"time"

	"github.com/antklim/go-invoice/invoice"
)

type testStorage struct {
	e   error
	inv invoice.Invoice
}

func (storage *testStorage) AddInvoice(inv invoice.Invoice) error {
	if storage.e != nil {
		return storage.e
	}

	storage.inv = inv
	return nil
}

func (storage testStorage) GetInvoice(id string) (invoice.Invoice, error) {
	return storage.inv, storage.e
}

var testErrStorage = &testStorage{e: errors.New("storage error")}

func TestCreateInvoice(t *testing.T) {
	invDate, err := time.Parse("2006-01-02", "2021-05-01")
	if err != nil {
		t.Fatalf("Error parsing invoice date: %v", err)
	}

	t.Run("creates valid invoice", func(t *testing.T) {
		srv := invoice.NewService(&testStorage{})
		inv, err := srv.CreateInvoice("John Doe", invDate)
		if err != nil {
			t.Fatalf("Error creating invoice: %v", err)
		}

		if inv.ID == "" {
			t.Fatal("invoice.ID should not be empty")
		}

		if inv.CustomerName != "John Doe" {
			t.Fatalf("invoice.CustomerName = %s, want John Doe", inv.CustomerName)
		}

		if !inv.Date.Equal(invDate) {
			t.Fatalf("invoice.Date = %s, want = %s", inv.Date.Format(time.RFC3339), invDate.Format(time.RFC3339))
		}

		if status := inv.Status; status != "open" {
			t.Fatalf("invoice.Status should be open, but got=%s", status)
		}

		if !inv.CreatedAt.Equal(inv.UpdatedAt) {
			t.Fatalf("invoice.CreatedAt = %s is not equalt to invoice.UpdatedAt = %s",
				inv.CreatedAt.Format(time.RFC3339),
				inv.UpdatedAt.Format(time.RFC3339))
		}
	})

	t.Run("stores invoice in data storage", func(t *testing.T) {
		srv := invoice.NewService(&testStorage{})
		inv, err := srv.CreateInvoice("John Doe", invDate)
		if err != nil {
			t.Fatalf("Error creating invoice: %v", err)
		}

		storedInvoice, err := srv.ViewInvoice(inv.ID)
		if err != nil {
			t.Fatalf("Error viewing invoice: %v", err)
		}

		if storedInvoice != inv {
			t.Fatalf("stored invoice is %v, want %v", storedInvoice, inv)
		}
	})

	t.Run("propagates data storage errors", func(t *testing.T) {
		srv := invoice.NewService(testErrStorage)
		inv, err := srv.CreateInvoice("Doe John", invDate)
		expectedErr := "failed to store invoice: storage error"
		if err.Error() != expectedErr {
			t.Fatalf("err is %v, want = %s", err, expectedErr)
		}

		expInv := invoice.Invoice{}
		if inv != expInv {
			t.Fatalf("invoice is %v, want = %v", inv, expInv)
		}
	})
}

func TestGetInvoice(t *testing.T) {
	t.Run("returns nothing when no invoice found", func(t *testing.T) {
		srv := invoice.NewService(&testStorage{})
		inv, err := srv.ViewInvoice("37f86bef-041d-4e50-aaf7-b1a066123751")
		if err != nil {
			t.Fatalf("Error viewing invoice: %v", err)
		}

		expInv := invoice.Invoice{}
		if inv != expInv {
			t.Fatalf("invoice is %v, want = %v", inv, expInv)
		}
	})

	t.Run("propagates data storage errors", func(t *testing.T) {
		srv := invoice.NewService(testErrStorage)
		inv, err := srv.ViewInvoice("37f86bef-041d-4e50-aaf7-b1a066123751")
		expectedErr := "failed to get invoice: storage error"
		if err.Error() != expectedErr {
			t.Fatalf("err is %v, want = %s", err, expectedErr)
		}

		expInv := invoice.Invoice{}
		if inv != expInv {
			t.Fatalf("invoice is %v, want = %v", inv, expInv)
		}
	})
}

func TestUpdateInvoice(t *testing.T) {
	t.Run("propagates data storage errors", func(t *testing.T) {
		t.Log("not implemented")
		t.Fail()
	})
}

func TestCloseInvoice(t *testing.T) {
	t.Run("propagates data storage errors", func(t *testing.T) {
		t.Log("not implemented")
		t.Fail()
	})
}

func TestOpenInvoice(t *testing.T) {
	t.Run("can be viewed", func(t *testing.T) {
		t.Log("not implemented")
		t.Fail()
	})

	t.Run("can be updated", func(t *testing.T) {
		t.Run("customer name can be updated", func(t *testing.T) {
			t.Log("not implemented")
			t.Fail()
		})

		t.Run("date can be updated", func(t *testing.T) {
			t.Log("not implemented")
			t.Fail()
		})

		t.Run("items can be added", func(t *testing.T) {
			t.Log("not implemented")
			t.Fail()
		})

		t.Run("items can be deleted", func(t *testing.T) {
			// when deleting non existent item it does not return error
			// error returned only in case of data access layer
			t.Log("not implemented")
			t.Fail()
		})
	})

	t.Run("can be issued", func(t *testing.T) {
		t.Log("not implemented")
		t.Fail()
	})

	t.Run("can be closed (cancel)", func(t *testing.T) {
		t.Log("not implemented")
		t.Fail()
	})

	t.Run("cannot be closed (paid)", func(t *testing.T) {
		t.Log("not implemented")
		t.Fail()
	})
}

func TestIssuedInvoice(t *testing.T) {
	t.Run("can be viewed", func(t *testing.T) {
		t.Log("not implemented")
		t.Fail()
	})

	t.Run("cannot be updated", func(t *testing.T) {
		t.Run("customer name cannot updated", func(t *testing.T) {
			t.Log("not implemented")
			t.Fail()
		})

		t.Run("date cannot be updated", func(t *testing.T) {
			t.Log("not implemented")
			t.Fail()
		})

		t.Run("items cannot be added", func(t *testing.T) {
			t.Log("not implemented")
			t.Fail()
		})

		t.Run("items cannot be deleted", func(t *testing.T) {
			t.Log("not implemented")
			t.Fail()
		})
	})

	t.Run("cannot be issued", func(t *testing.T) {
		t.Log("not implemented")
		t.Fail()
	})

	t.Run("can be closed (cancel)", func(t *testing.T) {
		t.Log("not implemented")
		t.Fail()
	})

	t.Run("can be closed (paid)", func(t *testing.T) {
		t.Log("not implemented")
		t.Fail()
	})
}

func TestClosedInvoice(t *testing.T) {
	t.Run("can be viewed", func(t *testing.T) {
		t.Log("not implemented")
		t.Fail()
	})

	t.Run("cannot be updated", func(t *testing.T) {
		t.Run("customer name cannot updated", func(t *testing.T) {
			t.Log("not implemented")
			t.Fail()
		})

		t.Run("date cannot be updated", func(t *testing.T) {
			t.Log("not implemented")
			t.Fail()
		})

		t.Run("items cannot be added", func(t *testing.T) {
			t.Log("not implemented")
			t.Fail()
		})

		t.Run("items cannot be deleted", func(t *testing.T) {
			t.Log("not implemented")
			t.Fail()
		})
	})

	t.Run("cannot be issued", func(t *testing.T) {
		t.Log("not implemented")
		t.Fail()
	})

	t.Run("cannot be closed (cancel)", func(t *testing.T) {
		t.Log("not implemented")
		t.Fail()
	})

	t.Run("cannot be closed (paid)", func(t *testing.T) {
		t.Log("not implemented")
		t.Fail()
	})
}
