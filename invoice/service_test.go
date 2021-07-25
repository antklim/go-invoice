package invoice_test

import (
	"testing"
	"time"

	"github.com/antklim/go-invoice/invoice"
)

func TestCreateInvoice(t *testing.T) {
	srv := &invoice.Service{}

	t.Run("created invoice should be valid", func(t *testing.T) {
		date, err := time.Parse("2006-01-02", "2021-05-01")
		if err != nil {
			t.Fatalf("Error parsing invoice date: %v", err)
		}

		inv, err := srv.CreateInvoice("John Doe", date)
		if err != nil {
			t.Fatalf("Error creating invoice: %v", err)
		}

		if inv.ID == "" {
			t.Fatal("invoice.ID should not be empty")
		}

		if inv.CustomerName != "John Doe" {
			t.Fatalf("invoice.CustomerName = %s, want John Doe", inv.CustomerName)
		}

		if !inv.Date.Equal(date) {
			t.Fatalf("invoice.Date = %s, want = %s", inv.Date.Format(time.RFC3339), date.Format(time.RFC3339))
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
		// t.Log("not implemented")
		// t.Fail()
	})

	t.Run("propagates data storage errors", func(t *testing.T) {
		// t.Log("not implemented")
		// t.Fail()
	})
}

func TestGetInvoice(t *testing.T) {
	t.Run("returns nothing when no invoice found", func(t *testing.T) {
		t.Log("not implemented")
		t.Fail()
	})

	t.Run("propagates data storage errors", func(t *testing.T) {
		t.Log("not implemented")
		t.Fail()
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
