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
  "base64SHA256Password": "string",
  "invalidUsernameAllowed": true
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

## reconfigureChargeStation

<a id="opIdreconfigureChargeStation"></a>

`POST /cs/{csId}/reconfigure`

*Reconfigure the charge station*

Supplies new configuration that should be applied to the charge station. This is not
intended to be used as a general charge station provisioning mechanism, it is intended
for one time changes required during testing. After reconfiguration, the charge station
will be rebooted so the new configuration can take effect if instructed to.

> Body parameter

```json
{
  "property1": "string",
  "property2": "string"
}
```

<h3 id="reconfigurechargestation-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|csId|path|string|false|The charge station identifier|
|body|body|[ChargeStationSettings](#schemachargestationsettings)|true|none|

> Example responses

> default Response

```json
{
  "status": "string",
  "error": "string"
}
```

<h3 id="reconfigurechargestation-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|None|
|default|Default|Unexpected error|[Status](#schemastatus)|

<aside class="success">
This operation does not require authentication
</aside>

## installChargeStationCertificates

<a id="opIdinstallChargeStationCertificates"></a>

`POST /cs/{csId}/certificates`

*Install certificates on the charge station*

> Body parameter

```json
{
  "certificates": [
    {
      "type": "V2G",
      "certificate": "string",
      "status": "Accepted"
    }
  ]
}
```

<h3 id="installchargestationcertificates-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|csId|path|string|false|The charge station identifier|
|body|body|[ChargeStationInstallCertificates](#schemachargestationinstallcertificates)|true|none|

> Example responses

> default Response

```json
{
  "status": "string",
  "error": "string"
}
```

<h3 id="installchargestationcertificates-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|None|
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
  "base64SHA256Password": "string",
  "invalidUsernameAllowed": true
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

## triggerChargeStation

<a id="opIdtriggerChargeStation"></a>

`POST /cs/{csId}/trigger`

> Body parameter

```json
{
  "trigger": "BootNotification"
}
```

<h3 id="triggerchargestation-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|csId|path|string|false|The charge station identifier|
|body|body|[ChargeStationTrigger](#schemachargestationtrigger)|true|none|

> Example responses

> default Response

```json
{
  "status": "string",
  "error": "string"
}
```

<h3 id="triggerchargestation-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|None|
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

## listTokens

<a id="opIdlistTokens"></a>

`GET /token`

*List authorization tokens*

Lists all tokens that can be used to authorize a charge

<h3 id="listtokens-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|offset|query|integer|false|none|
|limit|query|integer|false|none|

> Example responses

> 200 Response

```json
[
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
]
```

<h3 id="listtokens-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|List of tokens|Inline|
|default|Default|Unexpected error|[Status](#schemastatus)|

<h3 id="listtokens-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|[[Token](#schematoken)]|false|none|[An authorization token]|
|» countryCode|string|true|none|The country code of the issuing eMSP|
|» partyId|string|true|none|The party id of the issuing eMSP|
|» type|string|true|none|The type of token|
|» uid|string|true|none|The unique token id|
|» contractId|string|true|none|The contract ID (eMAID) associated with the token (with optional component separators)|
|» visualNumber|string|false|none|The visual/readable number/identification printed on an RFID card|
|» issuer|string|true|none|Issuing company, most of the times the name of the company printed on the RFID card, not necessarily the eMSP|
|» groupId|string|false|none|This id groups a couple of tokens to make two or more tokens work as one|
|» valid|boolean|true|none|Is this token valid|
|» languageCode|string|false|none|The preferred language to use encoded as ISO 639-1 language code|
|» cacheMode|string|true|none|Indicates what type of token caching is allowed|
|» lastUpdated|string(date-time)|false|none|The date the record was last updated (ignored on create/update)|

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

## registerParty

<a id="opIdregisterParty"></a>

`POST /register`

*Registers an OCPI party with the CSMS*

Registers an OCPI party with the CSMS. Depending on the configuration provided the CSMS will
either initiate a registration with the party or the party will wait for the party to initiate 
a registration with the CSMS.

> Body parameter

```json
{
  "token": "string",
  "url": "http://example.com",
  "status": "PENDING"
}
```

<h3 id="registerparty-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[Registration](#schemaregistration)|true|none|

> Example responses

> default Response

```json
{
  "status": "string",
  "error": "string"
}
```

<h3 id="registerparty-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|201|[Created](https://tools.ietf.org/html/rfc7231#section-6.3.2)|Created|None|
|default|Default|Unexpected error|[Status](#schemastatus)|

<aside class="success">
This operation does not require authentication
</aside>

## registerLocation

<a id="opIdregisterLocation"></a>

`POST /location/{locationId}`

*Registers a location with the CSMS*

Registers a location with the CSMS.

> Body parameter

```json
{
  "country_code": "string",
  "party_id": "string",
  "name": "string",
  "address": "string",
  "city": "string",
  "postal_code": "string",
  "country": "string",
  "coordinates": {
    "latitude": "string",
    "longitude": "string"
  },
  "parking_type": "ALONG_MOTORWAY",
  "evses": [
    {
      "uid": "string",
      "evse_id": "string",
      "connectors": [
        {
          "id": "string",
          "standard": "CHADEMO",
          "format": "SOCKET",
          "power_type": "AC_1_PHASE",
          "max_voltage": 0,
          "max_amperage": 0
        }
      ]
    }
  ]
}
```

<h3 id="registerlocation-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|locationId|path|string|false|The location identifier|
|body|body|[Location](#schemalocation)|true|none|

> Example responses

> default Response

```json
{
  "status": "string",
  "error": "string"
}
```

<h3 id="registerlocation-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|201|[Created](https://tools.ietf.org/html/rfc7231#section-6.3.2)|Created|None|
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
  "base64SHA256Password": "string",
  "invalidUsernameAllowed": true
}

```

Connection details for a charge station

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|securityProfile|integer|true|none|The security profile to use for the charge station: * `0` - unsecured transport with basic auth * `1` - TLS with basic auth * `2` - TLS with client certificate|
|base64SHA256Password|string|false|none|The base64 encoded, SHA-256 hash of the charge station password|
|invalidUsernameAllowed|boolean|false|none|If set to true then an invalid username will not prevent the charge station connecting|

<h2 id="tocS_ChargeStationSettings">ChargeStationSettings</h2>
<!-- backwards compatibility -->
<a id="schemachargestationsettings"></a>
<a id="schema_ChargeStationSettings"></a>
<a id="tocSchargestationsettings"></a>
<a id="tocschargestationsettings"></a>

```json
{
  "property1": "string",
  "property2": "string"
}

```

Settings for a charge station

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|**additionalProperties**|string|false|none|The key is the name of the setting. For OCPP 2.0.1 the name should have the following pattern:<br><component>/<variable>. The component name can include an optional component instance name and evse id<br>separated by semi-colons. The variable name can include an optional variable instance name and attribute<br>type separated by semi-colons. The maximum length for OCPP 1.6 is 500 characters.|

<h2 id="tocS_ChargeStationInstallCertificates">ChargeStationInstallCertificates</h2>
<!-- backwards compatibility -->
<a id="schemachargestationinstallcertificates"></a>
<a id="schema_ChargeStationInstallCertificates"></a>
<a id="tocSchargestationinstallcertificates"></a>
<a id="tocschargestationinstallcertificates"></a>

```json
{
  "certificates": [
    {
      "type": "V2G",
      "certificate": "string",
      "status": "Accepted"
    }
  ]
}

```

The set of certificates to install on the charge station. The certificates will be sent
to the charge station asynchronously.

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|certificates|[object]|true|none|none|
|» type|string|true|none|none|
|» certificate|string|true|none|The PEM encoded certificate with newlines replaced by `\n`|
|» status|string|false|none|The status, defaults to Pending|

#### Enumerated Values

|Property|Value|
|---|---|
|type|V2G|
|type|MO|
|type|MF|
|type|CSMS|
|status|Accepted|
|status|Rejected|
|status|Pending|

<h2 id="tocS_ChargeStationTrigger">ChargeStationTrigger</h2>
<!-- backwards compatibility -->
<a id="schemachargestationtrigger"></a>
<a id="schema_ChargeStationTrigger"></a>
<a id="tocSchargestationtrigger"></a>
<a id="tocschargestationtrigger"></a>

```json
{
  "trigger": "BootNotification"
}

```

Trigger a charge station action

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|trigger|string|true|none|none|

#### Enumerated Values

|Property|Value|
|---|---|
|trigger|BootNotification|
|trigger|StatusNotification|
|trigger|SignV2GCertificate|
|trigger|SignChargingStationCertificate|
|trigger|SignCombinedCertificate|

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

<h2 id="tocS_Registration">Registration</h2>
<!-- backwards compatibility -->
<a id="schemaregistration"></a>
<a id="schema_Registration"></a>
<a id="tocSregistration"></a>
<a id="tocsregistration"></a>

```json
{
  "token": "string",
  "url": "http://example.com",
  "status": "PENDING"
}

```

Defines the initial connection details for the OCPI registration process

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|token|string|true|none|The token to use for communicating with the eMSP (CREDENTIALS_TOKEN_A).|
|url|string(uri)|false|none|The URL of the eMSP versions endpoint. If provided the CSMS will act as the sender of the versions request.|
|status|string|false|none|The status of the registration request. If the request is marked as `REGISTERED` then the token will be allowed to<br>be used to access all endpoints avoiding the need for the OCPI registration process. If the request is marked as <br>`PENDING` then the token will only be allowed to access the `/ocpi/versions`, `/ocpi/2.2` and `/ocpi/2.2/credentials`<br>endpoints.|

#### Enumerated Values

|Property|Value|
|---|---|
|status|PENDING|
|status|REGISTERED|

<h2 id="tocS_Location">Location</h2>
<!-- backwards compatibility -->
<a id="schemalocation"></a>
<a id="schema_Location"></a>
<a id="tocSlocation"></a>
<a id="tocslocation"></a>

```json
{
  "country_code": "string",
  "party_id": "string",
  "name": "string",
  "address": "string",
  "city": "string",
  "postal_code": "string",
  "country": "string",
  "coordinates": {
    "latitude": "string",
    "longitude": "string"
  },
  "parking_type": "ALONG_MOTORWAY",
  "evses": [
    {
      "uid": "string",
      "evse_id": "string",
      "connectors": [
        {
          "id": "string",
          "standard": "CHADEMO",
          "format": "SOCKET",
          "power_type": "AC_1_PHASE",
          "max_voltage": 0,
          "max_amperage": 0
        }
      ]
    }
  ]
}

```

A charge station location

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|country_code|string|true|none|none|
|party_id|string|true|none|none|
|name|string¦null|false|none|none|
|address|string|true|none|none|
|city|string|true|none|none|
|postal_code|string¦null|false|none|none|
|country|string|true|none|none|
|coordinates|[GeoLocation](#schemageolocation)|true|none|none|
|parking_type|string¦null|false|none|none|
|evses|[[Evse](#schemaevse)]¦null|false|none|none|

#### Enumerated Values

|Property|Value|
|---|---|
|parking_type|ALONG_MOTORWAY|
|parking_type|PARKING_GARAGE|
|parking_type|PARKING_LOT|
|parking_type|ON_DRIVEWAY|
|parking_type|ON_STREET|
|parking_type|UNDERGROUND_GARAGE|

<h2 id="tocS_GeoLocation">GeoLocation</h2>
<!-- backwards compatibility -->
<a id="schemageolocation"></a>
<a id="schema_GeoLocation"></a>
<a id="tocSgeolocation"></a>
<a id="tocsgeolocation"></a>

```json
{
  "latitude": "string",
  "longitude": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|latitude|string|true|none|none|
|longitude|string|true|none|none|

<h2 id="tocS_Evse">Evse</h2>
<!-- backwards compatibility -->
<a id="schemaevse"></a>
<a id="schema_Evse"></a>
<a id="tocSevse"></a>
<a id="tocsevse"></a>

```json
{
  "uid": "string",
  "evse_id": "string",
  "connectors": [
    {
      "id": "string",
      "standard": "CHADEMO",
      "format": "SOCKET",
      "power_type": "AC_1_PHASE",
      "max_voltage": 0,
      "max_amperage": 0
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|uid|string|true|none|Uniquely identifies the EVSE within the CPOs platform (and<br>suboperator platforms).|
|evse_id|string¦null|false|none|none|
|connectors|[[Connector](#schemaconnector)]|true|none|none|

<h2 id="tocS_Connector">Connector</h2>
<!-- backwards compatibility -->
<a id="schemaconnector"></a>
<a id="schema_Connector"></a>
<a id="tocSconnector"></a>
<a id="tocsconnector"></a>

```json
{
  "id": "string",
  "standard": "CHADEMO",
  "format": "SOCKET",
  "power_type": "AC_1_PHASE",
  "max_voltage": 0,
  "max_amperage": 0
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string|true|none|none|
|standard|string|true|none|none|
|format|string|true|none|none|
|power_type|string|true|none|none|
|max_voltage|integer(int32)|true|none|none|
|max_amperage|integer(int32)|true|none|none|

#### Enumerated Values

|Property|Value|
|---|---|
|standard|CHADEMO|
|standard|CHAOJI|
|standard|DOMESTIC_A|
|standard|DOMESTIC_B|
|standard|DOMESTIC_C|
|standard|DOMESTIC_D|
|standard|DOMESTIC_E|
|standard|DOMESTIC_F|
|standard|DOMESTIC_G|
|standard|DOMESTIC_H|
|standard|DOMESTIC_I|
|standard|DOMESTIC_J|
|standard|DOMESTIC_K|
|standard|DOMESTIC_L|
|standard|GBT_AC|
|standard|GBT_DC|
|standard|IEC_60309_2_single_16|
|standard|IEC_60309_2_three_16|
|standard|IEC_60309_2_three_32|
|standard|IEC_60309_2_three_64|
|standard|IEC_62196_T1|
|standard|IEC_62196_T1_COMBO|
|standard|IEC_62196_T2|
|standard|IEC_62196_T2_COMBO|
|standard|IEC_62196_T3A|
|standard|IEC_62196_T3C|
|standard|NEMA_5_20|
|standard|NEMA_6_30|
|standard|NEMA_6_50|
|standard|NEMA_10_30|
|standard|NEMA_10_50|
|standard|NEMA_14_30|
|standard|NEMA_14_50|
|standard|PANTOGRAPH_BOTTOM_UP|
|standard|PANTOGRAPH_TOP_DOWN|
|standard|TESLA_R|
|standard|TESLA_S|
|standard|UNKNOWN|
|format|SOCKET|
|format|CABLE|
|power_type|AC_1_PHASE|
|power_type|AC_3_PHASE|
|power_type|DC|

