#!/bin/bash

aws dynamodb scan --table-name invoices \
  --endpoint-url http://localhost:8000 \
  --region ap-southeast-2
