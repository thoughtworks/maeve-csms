setup: setup-rfid setup-contract

.PHONY: setup-rfid
setup-rfid:
	@curl -i http://localhost:9410/api/v0/token -H 'content-type: application/json' -d '{\
	  "countryCode": "GB",\
	  "partyId": "TWK",\
	  "type": "RFID",\
	  "uid": "DEADBEEF",\
	  "contractId": "GBTWK012345678V",\
	  "issuer": "Thoughtworks",\
	  "valid": true,\
	  "cacheMode": "ALWAYS"\
	}'

.PHONY: setup-contract
setup-contract:
	@curl -i http://localhost:9410/api/v0/token -H 'content-type: application/json' -d '{\
	"countryCode": "GB",\
	"partyId": "TWK",\
	"type": "RFID",\
	"uid": "EMP77TWTW99999",\
	"contractId": "GBTWK012345678V",\
	"issuer": "Thoughtworks",\
	"valid": true,\
	"cacheMode": "ALWAYS"\
	}'

.PHONY: debug
debug:
	@curl -i http://localhost:9410/api/v0/cs/cs001 -H 'content-type: application/json' -d '{"securityProfile":2}' || docker-compose restart lb
