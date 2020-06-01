Некоторые запросы к Эластик выполненные из Кибаны


```sh

# show indices
GET _cat/indices?v&format=txt

# show index mapping
GET /log-generator-logrus-2020-05-27/_mapping

# show nodes
GET _cat/nodes?format=txt&v

# show shards
GET _cat/shards?v

# show xpack plugins
GET _xpack


#GET /_cluster/allocation/explain

GET /_cluster/allocation/explain
{
  "index": "log-generator-logrus-2020-05-27",
  "shard": 0,
  "primary": true
}

# create a doc 
PUT /customer/_doc/1
{
  "name": "John Doe"
}






GET /log-generator-logrus-2020-*/_search
{
  "query": {
    "match_all": {}
    
  },
  "size": 2, 
  
  "sort": [
    {
      "@timestamp": {
        "order": "asc"
      }
    }
  ]
}

# SQL query
POST _sql?format=txt
{
  "query": """
  SELECT 
  DATE_TRUNC('mi', "@timestamp") AS mins,
  ("@timestamp"::long / 60000 ) as div,
  * 
  FROM "log-generator-logrus*"
  order by "@timestamp" DESC
  limit 40
  """
}


# SQL query hist by 60 sec interval
POST _sql?format=txt
{
  "query": """
  SELECT 
  ("@timestamp"::long / 60000 ) as interval,
  (("@timestamp"::long /60000)*60000)::datetime as t,
  count(*) as N 
  FROM "log-generator-logrus*"
  GROUP BY interval,t
  order by interval DESC
  limit 40
  """
}



# SQL histogram query
POST _sql?format=txt
{
  "query": """
  SELECT HISTOGRAM("@timestamp", INTERVAL 30 SECOND) as interval, 
  count(*) as count
  FROM "log-generator-logrus*"
  GROUP by interval
  ORDER BY interval DESC
  LIMIT 30
  """
}



POST _sql?format=txt
{
  "query": """
  SELECT data.flight_number as flight ,count(*) as n FROM "log-generator-logrus-*"
  WHERE data.wait < 3000
  GROUP BY data.flight_number
  ORDER BY flight DESC
  LIMIT 3
  """
}

POST _sql/translate
{
  "query": """
  SELECT data.flight_number as flight ,count(*) as n FROM "log-generator-logrus-*"
  WHERE data.wait < 3000
  GROUP BY data.flight_number
  ORDER BY flight DESC
  LIMIT 3
  """
}


GET /log-generator-logrus-*/_search
{
  "size" : 1,
  "query" : {
    "range" : {
      "data.wait" : {
        "from" : null,
        "to" : 3000,
        "include_lower" : false,
        "include_upper" : false,
        "boost" : 1.0
      }
    }
  },
  "_source" : false,
  "stored_fields" : "_none_",
  "aggregations" : {
    "groupby" : {
      "composite" : {
        "size" : 3,
        "sources" : [
          {
            "e07b3457" : {
              "terms" : {
                "field" : "data.flight_number",
                "missing_bucket" : true,
                "order" : "desc"
              }
            }
          }
        ]
      }
    }
  }
}



```