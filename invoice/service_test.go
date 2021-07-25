package invoice_test

import "testing"

func TestCreateInvoice(t *testing.T) {
	t.Run("created invoice is in open status", func(t *testing.T) {
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
