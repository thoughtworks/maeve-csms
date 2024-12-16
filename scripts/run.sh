#!/usr/bin/env bash

# Check if 'docker compose' is available
if docker compose version >/dev/null 2>&1; then
  DOCKER_COMPOSE_CMD="docker compose"
elif docker-compose version >/dev/null 2>&1; then
  # Fallback to 'docker-compose' if 'docker compose' is not available
  DOCKER_COMPOSE_CMD="docker-compose"
else
  echo "Error: Neither 'docker-compose' nor 'docker compose' is available. Please install Docker Compose."
  exit 1
fi

# Run the Docker Compose command with any passed arguments
$DOCKER_COMPOSE_CMD up "${@:1}"

