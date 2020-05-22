Некоторые запросы к Эластик выполненые из Кибаны


```sh

# Получить все записи индекса.
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

# Выполните SQL запрос.
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

# Перевести SQL запрос в форматы эластик.
POST _sql/translate
{
  "query": """
  SELECT data.flight_number as flight ,count(*) as n FROM "log-generator-logrus-2020-05-19"
  WHERE data.wait < 3000
  GROUP BY data.flight_number
  ORDER BY flight DESC
  LIMIT 3
  """
}


GET /log-generator-logrus-2020-05-19/_search
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


# list indices
GET _cat/indices?v&format=txt


GET _cat/nodes?format=json

GET _cat/shards

GET _indices

GET _xpack


```