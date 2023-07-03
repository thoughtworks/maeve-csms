# Scripts

## Get CA Certificates

The [get-ca-certs.sh](get-ca-cert.sh) script retrieves the root V2G CA certificate and the CPO intermediate CA
certificates and writes them to the [config/certificates](../config/certificates) directory.

The script takes 1 argument:
1. The Hubject bearer token

The Hubject bearer token changes every day. You can extract the bearer token from
[Open Plug&Charge Testing Environment](https://hubject.stoplight.io/docs/open-plugncharge/6bb8b3bc79c2e-authorization-token).

The files will be called:
* `cpo_sub_ca1.pem`
* `cpo_sub_ca2.pem`
* `root-V2G-cert.pem`

Requires:
* curl
* openssl

## Generate TLS Certificate

The [generate-tls-cert.sh](../../ocpp2-broker-core/scripts/generate-tls-cert.sh) script creates a self-signed certificate that the CSMS can
use as its server certificate and writes it to the [config/certificates](../config/certificates) directory.

The script takes no arguments.

The output file is called `csms.pem`

Requires:
* openssl

## JSON to go

The [json-to-go.sh](../../ocpp2-broker-core/scripts/json-to-go.sh) script generates go types from the specified JSON schema.

It takes two arguments:
1. The OCPP version: either `ocpp16` or `ocpp201`
2. The name of the schema file from manager/schemas/<ocpp-version>

Requires:
* [gojsonschema](https://github.com/xeipuuv/gojsonschema)