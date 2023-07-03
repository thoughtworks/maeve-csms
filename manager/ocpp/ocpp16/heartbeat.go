package ocpp16

type HeartbeatJson map[string]interface{}

func (*HeartbeatJson) IsRequest() {}
