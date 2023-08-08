#!/usr/bin/env bash

command_exists() {
  command -v "$1" >/dev/null 2>&1
}

# Check if 'docker compose' is available (with space)
if command_exists "docker compose"; then
  DOCKER_COMPOSE_CMD="docker compose"
else
  # Check if 'docker-compose' is available
  if command_exists docker-compose; then
    DOCKER_COMPOSE_CMD="docker-compose"
  else
    echo "Error: Neither 'docker-compose' nor 'docker compose' is available. Please install Docker Compose."
    exit 1
  fi
fi

$DOCKER_COMPOSE_CMD up "${@:1}"
