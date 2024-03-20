#!/usr/bin/env bash

set -euo pipefail

bearer_token=$(curl -s https://hubject.stoplight.io/api/v1/projects/cHJqOjk0NTg5/nodes/6bb8b3bc79c2e-authorization-token | jq -r .data | sed -n '/Bearer/s/^.*Bearer //p')
cs_id=cs001
store_pass=123456
emaid="EMP77TWTW99999"
pcid="HUBOPENPROVCERT999"

while getopts ":b:c:e:p:" opt; do
  case "${opt}" in
  b )
    bearer_token="$OPTARG"
    ;;
  c )
    cs_id="$OPTARG"
    ;;
  e)
    emaid="$OPTARG"
    ;;
  p)
    pcid="$OPTARG"
    ;;
  \? )
    echo "Invalid option: $OPTARG" 1>&2
    echo "Usage: $0 [-b <hubject-token>] [-c <cs-id>] [-e <emaid>] [-p <pcid>]"
    exit 1
    ;;
  : )
    echo "Invalid option: $OPTARG requires an argument" 1>&2
    exit 1
  esac
done
shift $((OPTIND -1))

script_dir=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
cert_dir="${script_dir}"/../config/certificates
everest_dir="${script_dir}"/../config/everest

if [ ! -f "${cert_dir}"/cpo_sub_ca1.pem ]; then
  echo "Retrieving CPO cert chain..."
  "$script_dir"/get-ca-cert.sh "${bearer_token}"
fi

echo "Settings CS CPO cert chain..."
cp "${cert_dir}"/cpo_sub_ca1.pem "$everest_dir"/certs/ca/cso/CPO_SUB_CA1.pem
cp "${cert_dir}"/cpo_sub_ca2.pem "$everest_dir"/certs/ca/cso/CPO_SUB_CA2.pem

echo "Setting CS V2G root..."
cp "${cert_dir}"/root-V2G-cert.pem "$everest_dir"/certs/ca/v2g/V2G_ROOT_CA.pem

if [ ! -f "${cert_dir}/evccTruststore.jks" ]; then
  echo "Creating EVCC trust store with V2G root"
  keytool -import -keystore "${cert_dir}/evccTruststore.jks" -storepass "${store_pass}" -alias v2g_root -noprompt -file "${cert_dir}"/root-V2G-cert.pem
fi

echo "Setting EV trust store..."
cp "${cert_dir}"/evccTruststore.jks "${everest_dir}"/certs/client/oem/EVCC_TRUSTSTORE.jks

echo "Setting CSMS server certificate as EVerest TLS trust root"

if [ ! -f "${cert_dir}/csms.pem" ]; then
  echo "CSMS server certificate must be copied into the config/certificates directory"
  exit 1
fi

cp "${cert_dir}"/csms.pem "${everest_dir}"/certs/ca/csms

if [ ! -f "${cert_dir}"/${cs_id}.pem ]; then
  echo "Creating CS certificate..."
  "$script_dir"/generate-cs-cert.sh "${bearer_token}" "${cs_id}"
fi

echo "Setting CS certificate as EVerest SECC client certificate"
cp "${cert_dir}"/${cs_id}.pem "${everest_dir}"/certs/client/csms/CSMS_LEAF.pem
cp "${cert_dir}"/${cs_id}.key "${everest_dir}"/certs/client/csms/CSMS_LEAF.key

echo "Setting CS certificate as EVerest SECC server certificate"
cp "${cert_dir}"/${cs_id}.pem "${everest_dir}"/certs/client/cso/SECC_LEAF.pem
cp "${cert_dir}"/${cs_id}.key "${everest_dir}"/certs/client/cso/SECC_LEAF.key
echo "" > "${everest_dir}/certs/client/cso/SECC_LEAF_PASSWORD.txt"

if [ ! -f "${cert_dir}"/"${pcid}".p12 ]; then
  echo "Creating OEM provisioning certificate..."
  "$script_dir"/generate-oem-cert.sh "${bearer_token}" "${pcid}"
  "$script_dir"/register-oem-cert.sh "${bearer_token}" "${pcid}"
fi

if [ ! -f "${cert_dir}"/evccKeystore.jks ]; then
  echo "Creating EVCC key store with OEM provisioning certificate..."
  keytool -importkeystore -srckeystore "${cert_dir}"/${pcid}.p12 -srcstoretype pkcs12 -srcstorepass "${store_pass}" \
    -srcalias oem_prov_cert -destalias oem_prov_cert -destkeystore "${cert_dir}"/evccKeystore.jks -storepass "${store_pass}" -noprompt
fi

if [[ "${emaid}" != "" ]]; then
  if [ ! -f "${cert_dir}"/"${emaid}".p12 ]; then
    echo "Creating contract certificate..."
    "$script_dir"/generate-mo-cert.sh "${bearer_token}" "${emaid}"
  fi

  keytool -importkeystore -srckeystore "${cert_dir}"/"${emaid}".p12 -srcstoretype pkcs12 -srcstorepass "${store_pass}" \
    -srcalias contract_cert -destalias contract_cert -destkeystore "${cert_dir}"/evccKeystore.jks -storepass "${store_pass}" -noprompt
fi

echo "Setting EV key store..."
cp "${cert_dir}"/evccKeystore.jks "${everest_dir}"/certs/client/oem/EVCC_KEYSTORE.jks
