.PHONY: up
up: register-cs setup-rfid setup-contract
	OCPP_VERSION=ocpp docker-compose --profile everest up -d

.PHONY: test
test: # password 123456
	websocat --client-pkcs12-der ./config/certificates/cs001.pem --client-pkcs12-passwd 123456 wss://localhost:443

.PHONY: setup-rfid
setup-rfid:
	curl -i http://localhost:9410/api/v0/token -H 'content-type: application/json' -d '{"countryCode": "GB","partyId": "TWK","type": "RFID","uid": "DEADBEEF","contractId": "GBTWK012345678V","issuer": "Thoughtworks","valid": true,"cacheMode": "ALWAYS"}'

.PHONY: setup-contract
setup-contract:
	curl -i http://localhost:9410/api/v0/token -H 'content-type: application/json' -d '{"countryCode": "GB","partyId": "TWK","type": "RFID","uid": "EMP77TWTW99999","contractId": "GBTWK012345678V","issuer": "Thoughtworks","valid": true,"cacheMode": "ALWAYS"}'

.PHONY: register-cs
register-cs:
	curl -i http://localhost:9410/api/v0/cs/cs001 -H 'content-type: application/json' -d '{"securityProfile":2}'
