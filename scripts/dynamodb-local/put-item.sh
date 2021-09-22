#!/bin/bash

aws dynamodb put-item --table-name invoices \
    --item file://test/data/invoice.json \
    --endpoint-url http://localhost:8000 \
    --region ap-southeast-2 \
    --return-consumed-capacity TOTAL
