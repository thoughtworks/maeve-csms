package ocpp16

type AuthorizeJson struct {
	// IdTag corresponds to the JSON schema field "idTag".
	IdTag string `json:"idTag" yaml:"idTag" mapstructure:"idTag"`
}

func (*AuthorizeJson) IsRequest() {}
