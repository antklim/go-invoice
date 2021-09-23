#!/bin/bash

aws dynamodb create-table --table-name invoices \
  --attribute-definitions AttributeName=pk,AttributeType=S \
  --key-schema AttributeName=pk,KeyType=HASH \
  --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
  --endpoint-url ${AWS_ENDPOINT_URL}
