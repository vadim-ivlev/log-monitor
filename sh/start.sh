#!/bin/bash



echo "Starting stopped containers..."

docker-compose -f deploy/docker-compose.yml start

# поясняем
sh/greetings.sh