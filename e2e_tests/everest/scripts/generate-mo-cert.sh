#!/usr/bin/env bash

BEARER_TOKEN="$1"
EMAID="${2}"

if [[ "$BEARER_TOKEN" == "" ]]; then
  echo "You must provide a bearer token"
  exit 1
fi

BEARER_TOKEN=${BEARER_TOKEN#"Bearer "}

if [[ "$EMAID" == "" ]]; then
  echo "You must provide a EMAID"
  exit 1
fi

script_dir=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

openssl ecparam -name prime256v1 -genkey -noout -out "${script_dir}"/../config/certificates/"${EMAID}".key
openssl req -new -key "${script_dir}"/../config/certificates/"${EMAID}".key \
  -subj "/CN=${EMAID}/O=Thoughtworks" \
  -out "${script_dir}"/../config/certificates/"${EMAID}".csr \
  -outform DER \
  -sha256

curl -s https://open.plugncharge-test.hubject.com/mo/simpleenroll/ISO15118-2 \
  -H 'Accept: application/pkcs7' \
  -H "Authorization: Bearer ${BEARER_TOKEN}" \
  -H 'Content-Type: application/pkcs10' \
  -d "$(cat "${script_dir}"/../config/certificates/"${EMAID}".csr | base64)" | openssl enc -base64 -d > "${script_dir}"/../config/certificates/"${EMAID}".p7

openssl pkcs7 -in "${script_dir}"/../config/certificates/"${EMAID}".p7 -inform DER -print_certs -out "${script_dir}"/../config/certificates/"${EMAID}".pem
openssl pkcs12 -export -inkey "${script_dir}"/../config/certificates/"${EMAID}".key -in "${script_dir}"/../config/certificates/"${EMAID}".pem -name contract_cert -out "${script_dir}"/../config/certificates/"${EMAID}".p12 -password pass:123456
