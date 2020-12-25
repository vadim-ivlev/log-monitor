GET /_cluster/allocation/explain




GET /_cluster/allocation/explain
{
  "index": "customer",
  "shard": 0,
  "primary": true
}


PUT /customer/_doc/1
{
  "name": "John Doe"
}

curl -X PUT "localhost:9200/customer/_doc/1?pretty" -H 'Content-Type: application/json' -d'{  "name": "John Doe"}'



GET /_cluster/allocation/explain
{
  "index": "customer",
  "shard": 0,
  "primary": true
}






GET _cat/indices?v&format=json



 list indices
GET _cat/indices?v&format=txt


GET _cat/nodes?format=json

GET _cat/shards

GET _indices

GET _xpack




# Выполните SQL запрос.
POST _sql?format=txt
{
  "query": " SELECT *  FROM \"bank\" "
}

GET /bank/_search
{
    "query": {
        "bool": {
            "must": [
                {
                    "match_all": {}
                }
            ],
            "must_not": [],
            "should": []
        }
    },
    "from": 0,
    "size": 10,
    "sort": [],
    "aggs": {}
}


# set number of shards and replicas for index

PUT /kibana_sample_data_ecommerce
{
  "settings":{
    "number_of_shards":6,
    "number_of_replicas":2
  }
}