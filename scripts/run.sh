#!/usr/bin/env bash

command_exists() {
  command -v "$1" >/dev/null 2>&1
}

CSO_OPCP_TOKEN="$1"
MO_OPCP_TOKEN="$2"
if [[ "$CSO_OPCP_TOKEN" == "" ]]; then
  echo "You must provide a bearer token"
  echo "Usage: run.sh <CSO_OPCP_TOKEN> <MO_OPCP_TOKEN>"
  echo "       CSO_OPCP_TOKEN and MO_OPCP_TOKEN can be obtained from the Hubject test environment: "
  echo "       https://hubject.stoplight.io/docs/open-plugncharge/6bb8b3bc79c2e-authorization-token"
  exit 1
fi
CSO_OPCP_TOKEN=${CSO_OPCP_TOKEN#"Bearer "}
MO_OPCP_TOKEN=${MO_OPCP_TOKEN#"Bearer "}

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

MO_OPCP_TOKEN="Bearer "$MO_OPCP_TOKEN ;CSO_OPCP_TOKEN="Bearer "$CSO_OPCP_TOKEN ;$DOCKER_COMPOSE_CMD up "${@:2}"
