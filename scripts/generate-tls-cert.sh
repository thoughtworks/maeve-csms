#!/usr/bin/env bash

script_dir=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

openssl ecparam -name prime256v1 -genkey -noout -out "${script_dir}"/../config/certificates/csms.key
openssl req -new -nodes -key "${script_dir}"/../config/certificates/csms.key \
  -subj "/CN=CSMS/O=Thoughtworks" \
  -addext "subjectAltName = DNS:localhost, DNS:gateway, DNS:lb" \
  -out "${script_dir}"/../config/certificates/csms.csr
openssl x509 -req -in "${script_dir}"/../config/certificates/csms.csr \
  -out "${script_dir}"/../config/certificates/csms.pem \
  -signkey "${script_dir}"/../config/certificates/csms.key \
  -days 365 \
  -extfile <(printf "basicConstraints = critical, CA:false\n\
keyUsage = critical, digitalSignature, keyEncipherment\n\
subjectAltName = DNS:localhost, DNS:gateway, DNS:lb")
