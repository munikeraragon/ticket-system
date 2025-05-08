#!/bin/bash

set -e

echo "🚀 Running tests in Docker..."

docker compose -f docker-compose.test.yml up --build --abort-on-container-exit

echo "🧼 Cleaning up..."
docker compose -f docker-compose.test.yml down -v
