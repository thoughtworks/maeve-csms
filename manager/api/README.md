To regenerate code based on the OpenAPI spec:

```shell
$ go generate api.go
```

To generate the OpenAPI markdown docs:

```shell
$ npm install -g widdershins
$ widdershins api-spec.yaml -o API.md -c true
```