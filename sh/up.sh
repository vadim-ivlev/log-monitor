#!/bin/bash

echo "deleting log"
rm ./logs/generated.log

echo "поднимаем" 
docker-compose up -d

echo "поясняем"
sh/greetings.sh