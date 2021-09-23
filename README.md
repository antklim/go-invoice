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


# Project layout
```
.
+-- cli                 # interactive CLI implementation
+-- invoice             # core of the application
|   +-- invoice.go      # entities definitions
|   +-- service.go      # application logic (business rules) implementation
|   +-- storage.go      # application storage and storage factory interface definitions
|
+-- scripts             # misc scripts
|   +-- dynamodb        # dynamodb operations scripts such as create table, put item, etc.
|
+-- storage             # application storage concrete implementations
|   +-- dynamo          # DynamoDB storage implementation
|   +-- memory          # In memory storage implementation
|   +-- storage.go      # Storage factory implementation
|
+-- test                # test utilities, mocks, and fixtures
|   +-- api             # convinence APIs/DSL to set application in the state required by the test
|   +-- fixtures        # test data fixtures
|   +-- mocks           # various APIs mocks
|
+-- docker-compose.yml  # local DynamoDB service
+-- main.go             # go-invoice application entry point
+-- Makefile            # test, build and release tools
```

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
