#!/bin/bash

echo "stdin -> stdout "
echo "type a line and press ENTER"
# docker run --rm -it docker.elastic.co/logstash/logstash:7.6.2 logstash --log.level fatal -e 'input { stdin { } } output { stdout {} }'
docker run --rm -it  -v "$(pwd)/pipeline/:/usr/share/logstash/pipeline/"   docker.elastic.co/logstash/logstash:7.6.2 bin/logstash -f pipeline/first-pipeline.conf --config.reload.automatic