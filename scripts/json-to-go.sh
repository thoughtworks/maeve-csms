#!/usr/bin/env bash

function usage() {
    echo "Wrapper to extract Go code from JSON schema"
    echo
    echo "Usage: json-to-go.sh <ocpp-version> <type-name>"
    echo "  where <ocpp-version>" is either ocpp16 or ocpp201
    echo "  where <type-name> is the name of the type you want to extract from the json schema"
    echo
    echo "NOTE: this requires gojsonschema to be installed:"
    echo "go install github.com/atombender/go-jsonschema/cmd/gojsonschema@latest"
}

function snake_case() {
  echo "$1" \
    | sed 's/\([^A-Z]\)\([A-Z0-9]\)/\1_\2/g' \
    | sed 's/\([A-Z0-9]\)\([A-Z0-9]\)\([^A-Z]\)/\1_\2\3/g' \
    | tr '[:upper:]' '[:lower:]'
}

if [[ -z "$1" || -z "$2" ]]
then
  usage
  exit 1
fi

file=$(snake_case "$2")

$(go env GOPATH)/bin/gojsonschema -p "$1" manager/schemas/"$1"/"$2".json > manager/ocpp/"$1"/"$file".go

