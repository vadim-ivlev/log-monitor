#!/bin/bash

echo "logs -> stdout "
echo "add a line to ./logs/logstash-tutorial.log"
docker run --name=filebeat --rm -it --volume="$(pwd)/filebeat-logs-console.yml:/usr/share/filebeat/filebeat.yml" --volume="$(pwd)/logs:/logs" docker.elastic.co/beats/filebeat:7.6.2  filebeat