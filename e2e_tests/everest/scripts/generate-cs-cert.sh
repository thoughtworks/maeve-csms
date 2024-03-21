#!/usr/bin/env bash

BEARER_TOKEN="$1"
CS_NAME="${2:-cs001}"

if [[ "$BEARER_TOKEN" == "" ]]; then
  echo "You must provide a bearer token"
  exit 1
fi

BEARER_TOKEN=${BEARER_TOKEN#"Bearer "}

script_dir=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

openssl ecparam -name prime256v1 -genkey -noout -out "${script_dir}"/../config/certificates/"${CS_NAME}".key
openssl req -new -key "${script_dir}"/../config/certificates/"${CS_NAME}".key \
  -subj "/CN=${CS_NAME}/O=Thoughtworks" \
  -out "${script_dir}"/../config/certificates/"${CS_NAME}".csr \
  -outform DER \
  -sha256

curl -s https://open.plugncharge-test.hubject.com/cpo/simpleenroll/ISO15118-2 \
  -H 'Accept: application/pkcs7' \
  -H "Authorization: Bearer ${BEARER_TOKEN}" \
  -H 'Content-Type: application/pkcs10' \
  -d "$(cat "${script_dir}"/../config/certificates/"${CS_NAME}".csr | base64)" | openssl enc -base64 -d > "${script_dir}"/../config/certificates/"${CS_NAME}".p7

openssl pkcs7 -in "${script_dir}"/../config/certificates/"${CS_NAME}".p7 -inform DER -print_certs -out "${script_dir}"/../config/certificates/"${CS_NAME}".pem
