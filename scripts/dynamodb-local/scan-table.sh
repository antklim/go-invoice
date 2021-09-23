#!/bin/bash

aws dynamodb scan --table-name invoices \
  --endpoint-url ${AWS_ENDPOINT_URL}
