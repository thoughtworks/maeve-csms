#!/usr/bin/env bash

BEARER_TOKEN="$1"
PCID="${2}"

if [[ "$BEARER_TOKEN" == "" ]]; then
  echo "You must provide a bearer token"
  exit 1
fi

BEARER_TOKEN=${BEARER_TOKEN#"Bearer "}

if [[ "$PCID" == "" ]]; then
  echo "You must provide a PCID"
  exit 1
fi

script_dir=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

vehicleCertificate=$(cat "${script_dir}"/../config/certificates/"${PCID}".pem | awk '/subject.*CN.*=.*'${PCID}'/,/END CERTIFICATE/' | openssl x509 -outform DER | openssl enc -base64 -A)
subCA1Cert=$(cat "${script_dir}"/../config/certificates/"${PCID}".pem | awk '/subject.*CN.*=.*OEM Sub1 CA QA G1.2/,/END CERTIFICATE/' | openssl x509 -outform DER | openssl enc -base64 -A)
subCA2Cert=$(cat "${script_dir}"/../config/certificates/"${PCID}".pem | awk '/subject.*CN.*=.*OEM Sub2 CA QA G1.2.1/,/END CERTIFICATE/' | openssl x509 -outform DER | openssl enc -base64 -A)

curl --request PUT \
  --url https://open.plugncharge-test.hubject.com/v1/oem/provCerts \
  --header 'Accept: application/json, application/xml' \
  --header "Authorization: Bearer ${BEARER_TOKEN}" \
  --header 'Content-Type: application/json' \
  --data '{
  "subCA1Certificate": "'$subCA1Cert'",
  "subCA2Certificate": "'$subCA2Cert'",
  "vehicleCertificate": "'$vehicleCertificate'",
  "xsdMsgDefNamespace": "urn:iso:15118:2:2013:MsgDef",
  "rootIssuerDistinguishedName": "CN=V2G Root CA QA G1, DC=V2G, O=Hubject GmbH, C=DE",
  "rootIssuerSerialNumber": "69ab00d259bbdf42ce80529ad30ce5ed",
  "v2gRootAuthorityKeyIdentifier": "4b:45:ff:82:25:fc:10:96",
  "rootAuthorityKeyIdentifier": "4b:45:ff:82:25:fc:10:96"}'
