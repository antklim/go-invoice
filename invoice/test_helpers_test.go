package invoice_test

import (
	"os"

	"github.com/antklim/go-invoice/invoice"
	"github.com/antklim/go-invoice/storage"
	testapi "github.com/antklim/go-invoice/test/api"
)

func storageSetup() invoice.Storage {
	var f invoice.StorageFactory
	switch os.Getenv("TEST_STORAGE") {
	case "dynamo":
		tableName := "invoices"
		if os.Getenv("TEST_STORAGE_TABLE") != "" {
			tableName = os.Getenv("TEST_STORAGE_TABLE")
		}
		f = storage.NewDynamo(tableName, storage.WithEndpoint(os.Getenv("TEST_AWS_ENDPOINT")))
	default:
		f = new(storage.Memory)
	}
	strg := f.MakeStorage()
	return strg
}

func serviceSetup() (*invoice.Service, *testapi.Invoice) {
	strg := storageSetup()
	srv := invoice.New(strg)
	api := testapi.NewIvoiceAPI(strg)
	return srv, api
}
