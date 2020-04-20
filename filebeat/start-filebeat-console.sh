#!/bin/bash

echo "stdin -> stdout "
echo "type a line and press ENTER"
docker run --name=filebeat --rm -it -v "$(pwd)/filebeat-console.yml:/usr/share/filebeat/filebeat.yml" docker.elastic.co/beats/filebeat:7.6.2 filebeat