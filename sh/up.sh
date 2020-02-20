#!/bin/bash

# поднимаем 
docker-compose -f deploy/docker-compose.yml up -d

# поясняем
sh/greetings.sh