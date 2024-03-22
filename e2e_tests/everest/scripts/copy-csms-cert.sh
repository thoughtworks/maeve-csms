#!/usr/bin/env bash

script_dir=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
default_csms_dir="${script_dir}"/../../..
csms_dir="${1:-$default_csms_dir}"

cp "$csms_dir"/config/certificates/csms.pem "$script_dir"/../config/certificates
cp "$csms_dir"/config/certificates/csms.pem "$script_dir"/../config/everest/certs/ca/csms
