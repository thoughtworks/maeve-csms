#!/usr/bin/env bash

BEARER_TOKEN="$1"
if [[ "$BEARER_TOKEN" == "" ]]; then
  echo "You must provide a bearer token"
  exit 1
fi

BEARER_TOKEN=${BEARER_TOKEN#"Bearer "}

script_dir=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

certs=$(curl -s https://open.plugncharge-test.hubject.com/cpo/cacerts/ISO15118-2 \
  -H 'Accept: application/pkcs10, application/pkcs7' \
  -H "Authorization: Bearer ${BEARER_TOKEN}" \
  -H 'Content-Transfer-Encoding: application/pkcs10' | openssl enc -base64 -d | openssl pkcs7 -inform DER -print_certs)

echo "${certs}" | awk '/subject.*CN.*=.*CPO Sub1 CA QA G1.2/,/END CERTIFICATE/' > "${script_dir}"/../config/certificates/cpo_sub_ca1.pem
echo "${certs}" | awk '/subject.*CN.*=.*CPO Sub2 CA QA G1.2.1/,/END CERTIFICATE/' > "${script_dir}"/../config/certificates/cpo_sub_ca2.pem
echo "${certs}" | awk '/subject.*CN.*=.*V2G Root CA QA G1/,/END CERTIFICATE/' > "${script_dir}"/../config/certificates/root-V2G-cert.pem
cat "${script_dir}"/../config/certificates/cpo_sub_ca1.pem "${script_dir}"/../config/certificates/cpo_sub_ca2.pem > "${script_dir}"/../config/certificates/trust.pem

certs=$(curl -s https://open.plugncharge-test.hubject.com/mo/cacerts/ISO15118-2 \
  -H 'Accept: application/pkcs10, application/pkcs7' \
  -H "Authorization: Bearer ${BEARER_TOKEN}" \
  -H 'Content-Transfer-Encoding: application/pkcs10' | openssl enc -base64 -d | openssl pkcs7 -inform DER -print_certs)

echo "${certs}" | awk '/subject.*CN.*=.*MO Sub1 CA QA G1.2/,/END CERTIFICATE/' > "${script_dir}"/../config/certificates/mo_sub_ca1.pem
echo "${certs}" | awk '/subject.*CN.*=.*MO Sub2 CA QA G1.2.1/,/END CERTIFICATE/' > "${script_dir}"/../config/certificates/mo_sub_ca2.pem

