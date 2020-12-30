
log-monitor
========
Сбор и анализ логов приложений RG.RU в Elasticsearch 

-------------------------------
## Ссылки

- Конечная точка Elasticsearch для GET запросов <https://log-monitor.rg.ru/elasticsearch/>

- POST запросы к Elasticsearch ограничены только SQL запросами <https://log-monitor.rg.ru/elasticsearch/_sql>

    ```bash
    curl -XPOST "https://log-monitor.rg.ru/elasticsearch/_sql?format=txt" \
    -H 'Content-Type: application/json' \
    -d'{  "query": "SELECT \"@timestamp\", log.file.path, message  FROM \"app-logs*\"  LIMIT 40"}'    
    ```

- Индексы   <https://log-monitor.rg.ru/elasticsearch/_cat/indices?v&format=txt>

- Узлы <https://log-monitor.rg.ru/elasticsearch/_cat/nodes?format=txt&v>

- Kibana  <https://log-monitor.rg.ru>


- Управление докер-контейнерами 
<https://log-monitor.rg.ru/portainer/#/dashboard/>


- Пример запросов к Elasticsearch из Javascript <br>
  <https://log-monitor.rg.ru/www/><br> 
  <https://observablehq.com/d/9e4bdac324ef3667>




Файлы проекта размещены на: `dockerweb.rgwork.ru:/home/gitupdater/log-monitor-prod`


Мотивация
------------

**log-monitor** предназначен для сбора и анализа логов приложений. Это отличает данное приложение от <https://monitor.rg.ru> показывающего текущее состояние приложений.



Схема приложения
--------------


<!-- <img src="images/log-**monitor**.png"> -->
<img src="images/schema.png">
<!-- <img src="images/ELK.png"> -->

Контейнер log-monitor содержит кластер Elasticsearch, Kibana,  и Caddy.

- Кластер Elasticsearch - для хранения, поиска и анализа данных логов
- Kibana - для представления данных Elasticsearch в требуемом виде.
- Caddy - для обеспечения базовой аутентификации к Kibana и Elastic.



Приложения пишут логи в общий вольюм. Служба filebeat следит за файлами логов отправляет новые строки в Эластик.

<img src="images/log-monitor-ideas.png">


Настройка приложений для посылки логов в Эластик
-----------------------------------------------

Приложение должно сохранять логи в общий вольюм app-logs,  
за файлами которого следит специально настроенный filebeat.
В docker-compose.yml добавьте

    services:
        app-name:
        ...
            volumes: 
                - app-logs:/logs
            ...
    
    # внешний вольюм для логов
    volumes:
        app-logs:
            external: true

Простейший способ заставить приложение сохранять сообщения об ошибках в файл –
перенаправить stderr приложения в файл. Например в docker-compose.yml

    command: ./myapp 2>/logs/myapp.log


Другой способ - внутри самого приложения организовать запись в файл
лога, например средствами пакета ruslog. 
При этом нужно соблюдать следующие правила:

1. Каждая запись в логе должна занимать одну строку.
2. Сохраняйте логи в поддиректории названной по имени
    программы или включайте имя программы в имя файла лога. Так легче будет различать записи логов различных программ.

3. (необязательно) Вы можете сохранять логи в формате JSON. 
   Для того чтобы JSON был разобран на поля перед посылкой в Эластик имя файла должно содержать подстроку "json".




Просмотр логов приложений
-------------------------
Логи приложений сохраняются в индексе Эластик `app-logs-*`.
Запись лога в индексе имеет поля 
- **@timestamp** - Штамп даты времени записи лога. 
- **log.file.path** - Название файла лога. 
  Для фильтрации логов различных приложений.
- **message** - Содержание записи лога приложения. 
  Обычно содержит сообщение об ошибке.


Логи приложений могут быть просмотрены через **API** логов
1. **API** `END_POINT` = <https://log-monitor.rg.ru/elasticsearch>.
   Запросы должны иметь заголовок `Content-Type: application/json`.

- **POST SQL API**. Конечная точка :
    `END_POINT/_sql`.
    Допускаются только POST запросы в SQL формате.  

    **Пример**: Найти логи приложения `auth-proxy`,
    в сообщении которых присутствует слово `QueryRowMap`,
    между указанными датами. Отсортировать результат в порядке
    убывания даты и ограничить выдачу двадцатью записями.

```sql
POST https://log-monitor.rg.ru/elasticsearch/_sql?format=json&pretty
{
  "query": """
  SELECT "@timestamp", log.file.path, message  
  FROM "app-logs"
  WHERE 
  log.file.path = '/logs/auth-proxy.log'
  AND MATCH(message, 'QueryRowMap')
  AND "@timestamp" > '2020-12-28T23:03:08'
  AND "@timestamp" < '2020-12-30T01:03:20.15'
  ORDER BY "@timestamp" DESC
  LIMIT 20
  """
}
```

- **GET API**. Конечная точка:
    `END_POINT/_search`. 
    
    **Пример**: тот же запрос что и выше.

```json
GET https://log-monitor.rg.ru/elasticsearch/app-logs/_search
{
  "size" : 20,
  "query" : {
    "bool" : {
      "must" : [
        {
          "term" : {
            "log.file.path" : {
              "value" : "/logs/auth-proxy.log"
            }
          }
        },
        {
          "match" : {
            "message" : {
              "query" : "QueryRowMap"
            }
          }
        },
        {
          "range" : {
            "@timestamp" : {
              "from" : "2020-12-28T23:03:08",
              "to" : "2020-12-30T01:03:20.15"
            }
          }
        }
      ]
    }
  },
  "_source" : {
    "includes" : ["log.file.path", "message","@timestamp"]
  },
  "sort" : [
    {
      "@timestamp" : {
        "order" : "desc"
      }
    }
  ]
}

```

- **URI запросы в формате Lucene**. 
  
  Такие GET запросы укладываются в одну строку, выглядят как 
  обычная ссылка и во многих случаях более удобны. Пример:

  <a href='https://log-monitor.rg.ru/elasticsearch/app-logs/_search?pretty&size=3&sort=@timestamp:desc&_source_includes=@timestamp,message,log.file.path&q=log.file.path:"/logs/auth-proxy.log" AND @timestamp:["2020-12-28T23:07" TO "2020-12-29T01:10"] AND message:QueryRowMap'>
  https://log-monitor.rg.ru/elasticsearch/app-logs/_search?pretty&size=3&sort=@timestamp:desc&_source_includes=@timestamp,message,log.file.path&q=log.file.path:"/logs/auth-proxy.log" AND @timestamp:["2020-12-28T23:07" TO "2020-12-29T01:10"] AND message:QueryRowMap
  </a>
  <br><br>

  **Пояснение**:

  - `https://log-monitor.rg.ru/elasticsearch` - конечная точка
  - `/app-logs/_search` - искать в индексе app-logs
  - `?pretty` - красиво форматировать ответ 
  - `&size=3` - ограничить ответ тремя записями
  - `&sort=@timestamp:desc` - сортировать по времени в порядке убывания
  - `&_source_includes=@timestamp,message,log.file.path` - включать в ответ только эти поля
  - `&q=` - начало запроса Lucene
  - `log.file.path:"/logs/auth-proxy.log"` - искать в логах auth-proxy
  - `AND @timestamp:["2020-12-28T23:07" TO "2020-12-29T01:10"]` - искать между указанными датами
  - `AND message:QueryRowMap` - сообщение лога должно содержать строчку "QueryRowMap"

Можно опрашивать API с помощью команды curl.

Получить логи JSON формате

```bash
curl -XPOST "https://log-monitor.rg.ru/elasticsearch/_sql?format=json&pretty" \
-H 'Content-Type: application/json' \
-d'{  "query": "SELECT \"@timestamp\", log.file.path, message  FROM \"app-logs*\"  LIMIT 40"}'    
```
Получить логи текстовом формате

```bash
curl -XPOST "https://log-monitor.rg.ru/elasticsearch/_sql?format=txt" \
-H 'Content-Type: application/json' \
-d'{  "query": "SELECT \"@timestamp\", log.file.path, message  FROM \"app-logs*\"  LIMIT 40"}'    
```
   
1. Доступ к логам из Кибаны <https://log-monitor.rg.ru/app/kibana#/dev_tools/console>
2. Доступ к логам с помощью плугина Chrome [Elasticsearch Head](https://chrome.google.com/webstore/detail/elasticsearch-head/ffmkiejjmecolpfloofpjologoblkegm?hl=en-GB&utm_source=chrome-ntp-launcher). 
   
   Подсоединитесь к 
   <https://log-monitor.rg.ru/elasticsearch/>. На вкладке Any Request сделайте GET запрос к конечной точке _search, 
   или POST  к конечной точке _sql. Другие запросы запрещены.



<br><br><br>

Требования к системе для запуска Elasticsearch
-----------------------------------------

* **Минимум 4 ГБ RAM выделено для Docker**

    Для работы Elasticsearch требуется как минимум 2 ГБ оперативной памяти.

    Для Mac объем оперативной памяти, выделенной для Docker, можно установить с помощью пользовательского интерфейса

* **Лимит на mmap должен быть больше или равен 262,144**

    Это самая частая причина, по которой Elasticsearch не запускается с момента выпуска Elasticsearch версии 5.

    В Linux используйте `sysctl vm.max_map_count` на хосте, чтобы просмотреть текущее значение. Обратите внимание, что ограничения должны быть изменены на хосте; они не могут быть изменены из контейнера.

    Если вы используете Docker для Mac, вам потребуется запустить контейнер с переменной среды `MAX_MAP_COUNT`, установленной как минимум в 262144 (с использованием, например, опции `-e` докера), чтобы Elasticsearch установил ограничения на число `mmap` в время запуска.

* **Доступ к TCP-порту 5044 от клиентов, генерирующих логи**


Виртуальная память для запуска Elasticsearch
-----------

Elasticsearch по умолчанию использует директорию `hybrid mmapfs / niofs` для хранения своих индексов. По умолчанию ограничения  `mmap` слишком малы, что может привести к нехватке памяти.

В Linux вы можете увеличить ограничения, выполнив следующую команду от имени root:

    sysctl -w vm.max_map_count=262144

Чтобы установить это значение навсегда, обновите параметр `vm.max_map_count` в `/etc/sysctl.conf`. Для проверки после перезагрузки, запустите 

    sysctl vm.max_map_count

Пакеты RPM и Debian настроят этот параметр автоматически. Никаких дополнительных настроек не требуется.


<!--  
<br><br><br>

--------------------------

## Логи приложений

Для посылки логов в Elasticsearch приложения могут использовать две схемы:

1. **С использованием filebeat**. В контейнере работающего приложения должна быть запущена программа filebeat, 
   которая следит за файлами логов приложения и посылает новые записи логов в Эластик.
2. **Приложение напрямую посылает информацию** о логгируемых событиях в Эластик, возможно с параллельной записью в файл лога. 
   Сделать это возможно с помощью плагинов к logrus <https://github.com/interactive-solutions/go-logrus-elasticsearch>, 
   <https://github.com/sohlich/elogrus>, <https://github.com/go-extras/elogrus> 


Инструкции
===========
**Как добавить анализ логов к существующему приложению**

В docker-compose.yml существующего приложения необходимо внести следующие изменения:

1. Необходимо чтобы приложение работало в сети `auth_proxy_network`, что скорее всего уже выполнено,
поскольку предполагается, что приложения rgru работают под защитой сервера авторизации auth-proxy.
Если этого нет добавьте следующие строчки в верхний уровень файла docker-compose.yml:

    ```yml
    networks:
        default:
            external:
                name: auth_proxy_network    
    ```

2. Если приложение пишет файлы логов добавьте службу:

    ```yml
        filebeat:
            image: docker.elastic.co/beats/filebeat:7.6.2
            restart: unless-stopped
            volumes: 
                # настроечный файл Filebeat
                - ./configs/filebeat.yml:/usr/share/filebeat/filebeat.yml
                # директория где Filebeat следит за файлами логов
                - ./logs:/logs
            command: filebeat -e --strict.perms=false
    ```
    и добавьте настроечный файл filebeat.yml, пример которого можно посмотреть в [configs/filebeat.yml](configs/filebeat.yml).

**или**

3. Измените способ логирования приложения так, чтобы оно  посылало логи в Elasticsearch напрямую. 
   Сделать это можно вызвав следующую функцию во время инициализации логгера. 
   ```go
        addElasticHookToLogger(logger *logrus.Logger)
   ```
   Пример кода можно посмотреть в приложении log-generator 
   ( <https://git.rgwork.ru/ivlev/log-generator/blob/master/main.go> )



Рекомендации по форматам логов
---------------------------

1. logrus Лучше настроить так, чтобы формат выдачи был не текстовый, как установлено
по умолчанию, а JSON. 

2. Формат даты logrus лучше выбирать таким какой Эластик понимает по умолчанию.
   
Ниже приведены настройки форматера logrus, которые рекомендуется сделать.


```go
	// Java формат штампа даты времени, какой принят в Эластик по умолчанию.
	timeFormat := "2006-01-02T15:04:05.999Z"

	// Предпочитаем JSON формат вместо текстового, что удобно для анализа логов в Эластик.
	stdoutLog.SetFormatter(&logrus.JSONFormatter{TimestampFormat: timeFormat})
	
	// Текстовый формат логов
	// stdoutLog.SetFormatter(&logrus.TextFormatter{FullTimestamp:true, TimestampFormat: timeFormat})

```

Логирование метрик
-----------

Elasticsearch можно использовать для логирования метрик приложения, 
как альтернатива  Prometheus-Grafana. 

-->

<br><br><br>

Дополнительная информация
----

**Запросы к Эластик**

https://dzone.com/articles/23-useful-elasticsearch-example-queries


https://docs.google.com/document/d/1Q1ExyY36btdnTNe5co_pB4UdWNk41gY3rP1geg1LJBo/edit?usp=sharing



<br><br><br>

--------------------------

Порядок работы
==============

1. Изменить код
2. Запустить докер
3. Проверить
4. Запушить
5. Отдеплоить
   

Команды
-------
В директории `sh/` находятся следующие команды для облегчения работы.


|   |   |
|---|---|
Подъем                                      | `sh/up.sh`
Приостановка контейнера                     | `sh/stop.sh`
Старт приостановленного контейнера          | `sh/start.sh`
Полный останов контейнера                   | `sh/down.sh`
Подготовка директории deploy                | `sh/build-deploy-directory.sh`
Деплой                                      | `sh/deploy.sh`



Перезапуск Caddy и перестройка контейнера если что то изменилост в docker-compose 
```
dc restart log-monitor-caddy  
dc up -d --build log-monitor-caddy     
```