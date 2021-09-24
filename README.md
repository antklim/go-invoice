# go-invoice

Go-invoice is a simple invoice management application. A user can create, update, issue, pay and cancel invoice.
The invoice contains such information as ID, customer name and a list of items. The invoice item has its own ID and a product information such as product name, price and quantity.

When invoice created its open to updates:
- customer name and date can be updated
- items can be added and deleted

Invoice in any status can be viewed. But only invoices in open status can be updated.

The following diagram shows invoices statuses (in square brackets `[]`) and actions that cause status change (in parentheses `()`)
```
(Create Invoice)---> [OPEN] ---(Issue Invoice)---> [ISSUED] ---(Pay Invoice)---> [PAID]
                        |                              |
                 (Cancel Invoice)               (Cancel Invoice)
                           \                      /
                            +---> [CANCELED] <---+
```

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
To run tests simply call the following command:
```
$ make test
```

This command above runs all tests and calculates coverage. By default all tests run using in-memory storage.

Running tests using DynamoDB storage requires additional configuration. First, an instance of DynamoDB should be available for the test. The following command launches a local DynamoDB and creates `invoices` table:
```
$ docker-compose up
```
_Note_: DynamoDB docker container configured to bind host port `8000`. Make sure this port is available before you start a container or update ports binding settings.

The following command runs tests and uses local DynamoDB instance as storage:
```
$ AWS_PROFILE=local TEST_STORAGE=dynamo TEST_AWS_ENDPOINT=http://localhost:8000 make test
```
AWS SDK uses `AWS_PROFILE` to access user credentials to open a connection to AWS Resources. Use [aws configure](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html) to set AWS profile. The following is an example of `local` profile settings:
```
# ~/.aws/config
[profile local]
region = ap-southeast-2

# ~/.aws/credentials
[local]
aws_access_key_id = DUMMYIDEXAMPLE
aws_secret_access_key = DUMMYEXAMPLEKEY
```

Because of all resources running locally (in a docker container) a connection should be configured to use a custom endpoint URL. `TEST_AWS_ENDPOINT` parameter enables to do it.  
`TEST_STORAGE` tells to test suite what storage should it be using (by default in-memory storage used).

The following is the list of all supported test configuration options:
<table>
<thead><tr><th>Env variable</th><th>Description</th></thead>
<tbody>
<tr><td>
  TEST_AWS_ENDPOINT
</td><td>
  <p>A custom AWS endpoint URL. Used when runnning tests against local DynamoDB instance. The port number value should match ports binding settings in <b><i>docker-compose.yml</i></b>.</p>
</td></tr>
<tr><td>
  TEST_STORAGE
</td><td>
  <p>A storage to use when running test. Supported storages are</p>
  <ul>
    <li>memory - in-memory storage</li>
    <li>dynamo - DynamoDB storage</li>
  </ul>
  <p>By default in-memory storage used.</p>
</td></tr>
<tr><td>
  TEST_STORAGE_TABLE
</td><td>
  <p>A storage table name. Used only when <b><i>dynamo</i></b> storage selected. Default value is <b><i>invoices</i></b>.</p>
</td></tr>
</tbody>
</table>


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
