# go-invoice

Go-invoice is a simple invoice management application. A user can create, update, issue and close the invoice.
The invoice contains such information as ID, number, customer name and a list of items. The invoice item has its own ID and a product information such as SKU, product name, price and quantity.

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
