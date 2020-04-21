#!/bin/bash

echo "гасим"
docker-compose down


echo "поднимаем"
docker-compose up -d

echo "поясняем"
sh/greetings.sh