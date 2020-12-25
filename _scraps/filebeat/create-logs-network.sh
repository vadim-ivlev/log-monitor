#!/bin/bash

echo "creating logs_network used in docker-compose files for elk and filebeats"
docker network create logs_network
docker network ls
