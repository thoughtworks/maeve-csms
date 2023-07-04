#!/bin/sh
#
command_exists() {
  command -v "$1" >/dev/null 2>&1
}

BEARER_TOKEN="$1"
if [[ "$BEARER_TOKEN" == "" ]]; then
  echo "You must provide a bearer token"
  echo "Usage: run.sh <BEARER_TOKEN>"
  echo "       BEARER_TOKEN can be obtained from the Hubject test environment: "
  echo "       https://hubject.stoplight.io/docs/open-plugncharge/6bb8b3bc79c2e-authorization-token"
  exit 1
fi
BEARER_TOKEN=${BEARER_TOKEN#"Bearer "}

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

HUBJECT_TOKEN="Bearer "$BEARER_TOKEN docker-compose up
