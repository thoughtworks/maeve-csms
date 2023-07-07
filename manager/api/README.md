To regenerate code based on the OpenAPI spec:

```shell
$ go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
$ oapi-codegen -config api/cfg.yaml api/api-spec.yaml
```