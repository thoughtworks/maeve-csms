# Scripts

## Get CA Certificates

The [get-ca-cert.sh](get-ca-cert.sh) script retrieves the root V2G CA certificate and the CPO intermediate CA
certificates and writes them to the [config/certificates](../config/certificates) directory.

The files will be called:
* `cpo_sub_ca1.pem`
* `cpo_sub_ca2.pem`
* `root-V2G-cert.pem`

Requires:
* curl
* openssl

## JSON to go

The [json-to-go.sh](json-to-go.sh) script generates go types from the specified JSON schema.

It takes two arguments:
1. The OCPP version: either `ocpp16` or `ocpp201`
2. The name of the schema file from manager/schemas/<ocpp-version>

Requires:
* [gojsonschema](https://github.com/xeipuuv/gojsonschema)