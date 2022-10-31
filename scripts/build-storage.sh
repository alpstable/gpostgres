#!/bin/bash

set -e

# remove volumes
rm -rf .db

# drop existing containers
docker compose -f "docker-compose.yml" down

# prune containers
docker system prune --force

docker-compose -f "docker-compose.yml" up -d \
	--remove-orphans \
	--force-recreate \
	--build postgres1
