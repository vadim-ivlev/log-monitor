#!/bin/bash

# гасим 
docker-compose -f deploy/docker-compose.yml down


# поднимаем 
docker-compose -f deploy/docker-compose.yml up -d

# поясняем
sh/greetings.sh