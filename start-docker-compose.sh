#!/bin/bash
# start-docker-compose.sh

# Загружаем значения из config.yml в переменные окружения
export $(cat config.yml | grep -v '^#' | awk '/database:/ {getline; getline; print "POSTGRES_DB="$2; getline; print "POSTGRES_USER="$2; getline; print "POSTGRES_PASSWORD="$2}' | xargs)

# Запускаем docker-compose с этими переменными
docker-compose up -d