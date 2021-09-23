#!/bin/bash

aws dynamodb put-item --table-name invoices \
    --item file://test/data/invoice.json \
    --endpoint-url ${AWS_ENDPOINT_URL} \
    --return-consumed-capacity TOTAL
