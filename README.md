# go-invoice

TODO: update with the latest info

Go-invoice is a simple invoice management application. A user can create, update, issue and close the invoice.
The invoice contains such information as ID, customer name and a list of items. The invoice item has its own ID and a product information such as SKU, product name, price and quantity.

When invoice created its open to updates:
- customer name and date can be updated
- items can either be deleted or added

After invoice being issued it cannot be updated. Issued invoice can be viewed.
Closing the invoice can be done in two ways: cancel the invoice and comfirm invoice payment. Closed invoice cannot be updated.

|        | Can be viewed | Can be updated |
|--------|---------------|----------------|
|Open    | YES           | YES            |
|Issued  | YES           | NO             |
|Closed  | YES           | NO             |


The following table shows the invoice status transitions:

|        | Open | Issued | Closed |
|--------|------|--------|--------|
|Open    | NO   | YES    | YES    |
|Issued  | NO   | NO     | YES    |
|Closed  | NO   | NO     | NO     |


TODO: add project structure

# Testing
To run test simply call:
```
$ make test
```

This will run all tests and calculate coverage. By default all tests run using in-memory storage. To run tests using DynamoDB storage, additional parameters should be provided:
```
$ AWS_PROFILE=local TEST_STORAGE=dynamo TEST_AWS_ENDPOINT=http://localhost:8000 make test
```

TODO: explain prerequisites
TODO: add supported env var flags and values for testing

# Usage
To launch the `go-invoice` application run the following command:
```
$ go run main.go
```

By default the application uses in-memory storage. To configure application to use DynamoDB, additional parameters shuld be provided:
```
$ AWS_PROFILE=local go run main.go -storage=dynamo -endpoint=http://localhost:8000
```

_Note_: it's important to provide protocol when configuring an endpoint. Just `localhost:8000` does not work.
