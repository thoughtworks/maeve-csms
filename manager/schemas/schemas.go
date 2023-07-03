package schemas

import "embed"

//go:embed */*.json
var OcppSchemas embed.FS
