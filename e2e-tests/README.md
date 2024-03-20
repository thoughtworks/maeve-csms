# End-to-End test

We currently use two different versions of EVerest and the MQTT API for RFID authorisation
is different between them: if using the older "v5" image (the default):

```shell
$ NO_AUTH_TOKEN_TYPE_PREFIX=1 go test -v ./... -count=1
```

Otherwise,

```shell
$ go test -v ./... -count=1
```

