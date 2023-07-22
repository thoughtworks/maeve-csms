---
title: MaEVe CSMS v0.0.0
language_tabs:
  - shell: Shell
  - http: HTTP
  - javascript: JavaScript
  - ruby: Ruby
  - python: Python
  - php: PHP
  - java: Java
  - go: Go
toc_footers: []
includes: []
search: true
highlight_theme: darkula
headingLevel: 2

---

<!-- Generator: Widdershins v4.0.1 -->

<h1 id="maeve-csms">MaEVe CSMS v0.0.0</h1>

> Scroll down for code samples, example requests and responses. Select a language for code samples from the tabs above or the mobile navigation menu.

Internal API to interact with the MaEVe CSMS, external clients should use OCPI.

Base URLs:

* <a href="http://localhost:9410/api/v0">http://localhost:9410/api/v0</a>

Email: <a href="mailto:maeve-team@thoughtworks.com">MaEVe team</a> 
 License: Apache 2.0

<h1 id="maeve-csms-default">Default</h1>

## registerChargeStation

<a id="opIdregisterChargeStation"></a>

`POST /cs/{csId}`

*Register a new charge station*

Registers a new charge station. The system will assume that the charge station
has not yet been provisioned and will place the charge station into a pending state
so it can been configured when it sends a boot notification.

> Body parameter

```json
{
  "securityProfile": 0,
  "base64SHA256Password": "string"
}
```

<h3 id="registerchargestation-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|csId|path|string|false|The charge station identifier|
|body|body|[ChargeStationAuth](#schemachargestationauth)|true|none|

> Example responses

> default Response

```json
{
  "status": "string",
  "error": "string"
}
```

<h3 id="registerchargestation-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|201|[Created](https://tools.ietf.org/html/rfc7231#section-6.3.2)|Created|None|
|default|Default|Unexpected error|[Status](#schemastatus)|

<aside class="success">
This operation does not require authentication
</aside>

## lookupChargeStationAuth

<a id="opIdlookupChargeStationAuth"></a>

`GET /cs/{csId}/auth`

*Returns the authentication details*

Returns the details required by the CSMS gateway to determine how to authenticate
the charge station

<h3 id="lookupchargestationauth-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|csId|path|string|false|The charge station identifier|

> Example responses

> 200 Response

```json
{
  "securityProfile": 0,
  "base64SHA256Password": "string"
}
```

<h3 id="lookupchargestationauth-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Charge station auth response|[ChargeStationAuth](#schemachargestationauth)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Unknown charge station|[Status](#schemastatus)|
|default|Default|Unexpected error|[Status](#schemastatus)|

<aside class="success">
This operation does not require authentication
</aside>

## setToken

<a id="opIdsetToken"></a>

`POST /token`

*Create/update an authorization token*

Creates or updates a token that can be used to authorize a charge

> Body parameter

```json
{
  "countryCode": "st",
  "partyId": "str",
  "type": "AD_HOC_USER",
  "uid": "string",
  "contractId": "string",
  "visualNumber": "string",
  "issuer": "string",
  "groupId": "string",
  "valid": true,
  "languageCode": "st",
  "cacheMode": "ALWAYS",
  "lastUpdated": "2019-08-24T14:15:22Z"
}
```

<h3 id="settoken-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[Token](#schematoken)|true|none|

> Example responses

> default Response

```json
{
  "status": "string",
  "error": "string"
}
```

<h3 id="settoken-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|201|[Created](https://tools.ietf.org/html/rfc7231#section-6.3.2)|Created|None|
|default|Default|Unexpected error|[Status](#schemastatus)|

<aside class="success">
This operation does not require authentication
</aside>

## lookupToken

<a id="opIdlookupToken"></a>

`GET /token/{tokenUid}`

*Lookup an authorization token*

Lookup a token that can be used to authorize a charge

<h3 id="lookuptoken-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|tokenUid|path|string|true|none|

> Example responses

> 200 Response

```json
{
  "countryCode": "st",
  "partyId": "str",
  "type": "AD_HOC_USER",
  "uid": "string",
  "contractId": "string",
  "visualNumber": "string",
  "issuer": "string",
  "groupId": "string",
  "valid": true,
  "languageCode": "st",
  "cacheMode": "ALWAYS",
  "lastUpdated": "2019-08-24T14:15:22Z"
}
```

<h3 id="lookuptoken-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Authorization token details|[Token](#schematoken)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found|[Status](#schemastatus)|
|default|Default|Unexpected error|[Status](#schemastatus)|

<aside class="success">
This operation does not require authentication
</aside>

## uploadCertificate

<a id="opIduploadCertificate"></a>

`POST /certificate`

*Upload a certificate*

Uploads a client certificate to the CSMS. The CSMS can use the certificate to authenticate
the charge station using mutual TLS when the TLS operations are being offloaded to a load-balancer.

> Body parameter

```json
{
  "certificate": "string"
}
```

<h3 id="uploadcertificate-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[Certificate](#schemacertificate)|true|none|

> Example responses

> default Response

```json
{
  "status": "string",
  "error": "string"
}
```

<h3 id="uploadcertificate-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|201|[Created](https://tools.ietf.org/html/rfc7231#section-6.3.2)|Created|None|
|default|Default|Unexpected error|[Status](#schemastatus)|

<aside class="success">
This operation does not require authentication
</aside>

## lookupCertificate

<a id="opIdlookupCertificate"></a>

`GET /certificate/{certificateHash}`

*Lookup a certificate*

Lookup a client certificate that has been uploaded to the CSMS using a base64 encoded SHA-256 hash
of the DER bytes.

<h3 id="lookupcertificate-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|certificateHash|path|string|true|none|

> Example responses

> 200 Response

```json
{
  "certificate": "string"
}
```

<h3 id="lookupcertificate-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Certificate details|[Certificate](#schemacertificate)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found|[Status](#schemastatus)|
|default|Default|Unexpected error|[Status](#schemastatus)|

<aside class="success">
This operation does not require authentication
</aside>

## deleteCertificate

<a id="opIddeleteCertificate"></a>

`DELETE /certificate/{certificateHash}`

*Delete a certificate*

Deletes a client certificate that has been uploaded to the CSMS using a base64 encoded SHA-256 hash
of the DER bytes.

<h3 id="deletecertificate-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|certificateHash|path|string|true|none|

> Example responses

> 404 Response

```json
{
  "status": "string",
  "error": "string"
}
```

<h3 id="deletecertificate-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|204|[No Content](https://tools.ietf.org/html/rfc7231#section-6.3.5)|No content|None|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found|[Status](#schemastatus)|
|default|Default|Unexpected error|[Status](#schemastatus)|

<aside class="success">
This operation does not require authentication
</aside>

# Schemas

<h2 id="tocS_ChargeStationAuth">ChargeStationAuth</h2>
<!-- backwards compatibility -->
<a id="schemachargestationauth"></a>
<a id="schema_ChargeStationAuth"></a>
<a id="tocSchargestationauth"></a>
<a id="tocschargestationauth"></a>

```json
{
  "securityProfile": 0,
  "base64SHA256Password": "string"
}

```

Connection details for a charge station

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|securityProfile|integer|true|none|The security profile to use for the charge station: * `0` - unsecured transport with basic auth * `1` - TLS with basic auth * `2` - TLS with client certificate|
|base64SHA256Password|string|false|none|The base64 encoded, SHA-256 hash of the charge station password|

<h2 id="tocS_Token">Token</h2>
<!-- backwards compatibility -->
<a id="schematoken"></a>
<a id="schema_Token"></a>
<a id="tocStoken"></a>
<a id="tocstoken"></a>

```json
{
  "countryCode": "st",
  "partyId": "str",
  "type": "AD_HOC_USER",
  "uid": "string",
  "contractId": "string",
  "visualNumber": "string",
  "issuer": "string",
  "groupId": "string",
  "valid": true,
  "languageCode": "st",
  "cacheMode": "ALWAYS",
  "lastUpdated": "2019-08-24T14:15:22Z"
}

```

An authorization token

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|countryCode|string|true|none|The country code of the issuing eMSP|
|partyId|string|true|none|The party id of the issuing eMSP|
|type|string|true|none|The type of token|
|uid|string|true|none|The unique token id|
|contractId|string|true|none|The contract ID (eMAID) associated with the token (with optional component separators)|
|visualNumber|string|false|none|The visual/readable number/identification printed on an RFID card|
|issuer|string|true|none|Issuing company, most of the times the name of the company printed on the RFID card, not necessarily the eMSP|
|groupId|string|false|none|This id groups a couple of tokens to make two or more tokens work as one|
|valid|boolean|true|none|Is this token valid|
|languageCode|string|false|none|The preferred language to use encoded as ISO 639-1 language code|
|cacheMode|string|true|none|Indicates what type of token caching is allowed|
|lastUpdated|string(date-time)|false|none|The date the record was last updated (ignored on create/update)|

#### Enumerated Values

|Property|Value|
|---|---|
|type|AD_HOC_USER|
|type|APP_USER|
|type|OTHER|
|type|RFID|
|cacheMode|ALWAYS|
|cacheMode|ALLOWED|
|cacheMode|ALLOWED_OFFLINE|
|cacheMode|NEVER|

<h2 id="tocS_Status">Status</h2>
<!-- backwards compatibility -->
<a id="schemastatus"></a>
<a id="schema_Status"></a>
<a id="tocSstatus"></a>
<a id="tocsstatus"></a>

```json
{
  "status": "string",
  "error": "string"
}

```

HTTP status

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|status|string|true|none|The status description|
|error|string|false|none|The error details|

<h2 id="tocS_Certificate">Certificate</h2>
<!-- backwards compatibility -->
<a id="schemacertificate"></a>
<a id="schema_Certificate"></a>
<a id="tocScertificate"></a>
<a id="tocscertificate"></a>

```json
{
  "certificate": "string"
}

```

A client certificate

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|certificate|string|true|none|The PEM encoded certificate with newlines replaced by `\n`|

