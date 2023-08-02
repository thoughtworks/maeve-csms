#!/usr/bin/env bash

command_exists() {
  command -v "$1" >/dev/null 2>&1
}

BEARER_TOKEN=$(curl -s https://hubject.stoplight.io/api/v1/projects/cHJqOjk0NTg5/nodes/6bb8b3bc79c2e-authorization-token | jq -r .data | sed -n '/Bearer/s/^.*Bearer //p')

# fall back to BEARER_TOKEN if no arg
CSO_OPCP_TOKEN="${1:-$BEARER_TOKEN}"
MO_OPCP_TOKEN="${2:-$BEARER_TOKEN}"

shift

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

export MO_OPCP_TOKEN=$MO_OPCP_TOKEN; export CSO_OPCP_TOKEN=$CSO_OPCP_TOKEN;$DOCKER_COMPOSE_CMD up "${@:2}"
