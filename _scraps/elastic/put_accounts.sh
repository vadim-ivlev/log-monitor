#!/bin/bash

curl -XDELETE "localhost:9200/bank"
curl -H "Content-Type: application/json" -XPOST "localhost:9200/bank/_bulk?pretty&refresh" --data-binary "@accounts.json" > /dev/null
curl "localhost:9200/_cat/indices?v"

