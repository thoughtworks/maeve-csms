package ocpp16

type HeartbeatResponseJson struct {
	// CurrentTime corresponds to the JSON schema field "currentTime".
	CurrentTime string `json:"currentTime" yaml:"currentTime" mapstructure:"currentTime"`
}

func (*HeartbeatResponseJson) IsResponse() {}
