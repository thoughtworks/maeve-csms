# End-to-end tests

## Description

Runs end-to-end tests against Everest (Charge station and EV simulator) along with CSMS. This performs the following tests:
- Plug-in a connector
- Authorise a charge (using RFID or ISO-15118)
- Start a charge
- Stop a charge
- Unplug a connector


## Steps

1. Run the following bash script to start the docker containers for Everest and CSMS and execute the end-to-end tests
```shell
./run-e2e-tests.sh
```





