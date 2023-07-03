package schemas

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/santhosh-tekuri/jsonschema"
	"github.com/santhosh-tekuri/jsonschema/loader"
	"io"
	"io/fs"
	"net/url"
)

type FSLoader struct {
	FS fs.FS
}

func (f FSLoader) Load(u string) (io.ReadCloser, error) {
	parsedUrl, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	if parsedUrl.Scheme == "fs" {
		return f.FS.Open(parsedUrl.Path[1:])
	}
	return nil, nil
}

func Validate(data json.RawMessage, schemaFs fs.FS, schemaFile string) error {
	compiler := jsonschema.NewCompiler()
	compiler.Draft = jsonschema.Draft6
	loader.Register("fs", FSLoader{
		FS: schemaFs,
	})
	schema, err := compiler.Compile(fmt.Sprintf("fs:///%s", schemaFile))
	if err != nil {
		return err
	}

	return schema.Validate(bytes.NewReader(data))
}
