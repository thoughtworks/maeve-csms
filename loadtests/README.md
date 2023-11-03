# Project Name

Load Testing with K6

## Table of Contents

- [Project Description](#project-description)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Getting Started](#getting-started)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Project Description

To set up a load test against the CSMS back-end using a number of simulated charge stations. 

## Prerequisites
1. Ensure that the CSMS backend is running via docker

## Installation

1. Install k6 onto your local machine
```bash
brew install k6
```
2. Download xk6 using go
```bash
go install go.k6.io/xk6/cmd/xk6@latest
```

3. Build the binary
```bash
xk6 build --with github.com/grafana/xk6-dashboard@latest
```

## Getting Started

1. Extract and copy the base64SHA256Password using the following command
```bash
(cd manager && go run main.go auth encode-password fiddlesticks_fishsticks)
```
2. Register a number of charge stations to the CSMS (e.g. cs1, cs2, cs3 etc). Replace 'BASE64_SHA256_PASSWORD' with the base64SHA256Password that was extracted
```bash
curl http://localhost:9410/api/v0/cs/cs1 -H 'content-type: application/json' -d '{"securityProfile":0,"base64SHA256Password":"BASE64_SHA256_PASSWORD"}' &&
curl http://localhost:9410/api/v0/cs/cs2 -H 'content-type: application/json' -d '{"securityProfile":0,"base64SHA256Password":"BASE64_SHA256_PASSWORD"}' &&
curl http://localhost:9410/api/v0/cs/cs3 -H 'content-type: application/json' -d '{"securityProfile":0,"base64SHA256Password":"BASE64_SHA256_PASSWORD"}'
```

3. Register the contract token to the CSMS. Replace 'UID' with the value of the idTag that is found in loadtests/ws_load_test.js. This is used in the websocket messages: Authorise, StartTransaction and StopTransaction.
```bash 
curl -i http://localhost:9410/api/v0/token -H 'content-type: application/json' -d '{"countryCode": "GB","partyId": "TWK","type": "RFID","uid": "UID","contractId": "GBTWK012345678V","issuer": "Thoughtworks","valid": true,"cacheMode": "ALWAYS"}'
```

4. Set the load simulation for ramping virtual users in loadtests/ws_load_test.js. Please refer to https://k6.io/docs/using-k6/scenarios/executors/ramping-vus/ for guidance. Please note that 1 virtual user is the equivalent to 1 charge station.


5. Run the load test using the script file and output the results to the k6 dashboard
```bash
./k6 run --out dashboard loadtests/ws_load_test.js 
```




