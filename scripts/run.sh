#!/usr/bin/env bash

BEARER_TOKEN=$(curl -s https://hubject.stoplight.io/api/v1/projects/cHJqOjk0NTg5/nodes/6bb8b3bc79c2e-authorization-token | jq -r .data | sed -n '/Bearer/s/^.*Bearer //p')

# fall back to BEARER_TOKEN if no arg
CSO_OPCP_TOKEN="${1:-$BEARER_TOKEN}"
MO_OPCP_TOKEN="${2:-$BEARER_TOKEN}"

export MO_OPCP_TOKEN=$MO_OPCP_TOKEN; export CSO_OPCP_TOKEN=$CSO_OPCP_TOKEN;docker-compose up "${@:2}" --build
