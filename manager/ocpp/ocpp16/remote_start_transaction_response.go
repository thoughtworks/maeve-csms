package ocpp16

type RemoteStartTransactionResponseJsonStatus string

type RemoteStartTransactionResponseJson struct {
	// Status corresponds to the JSON schema field "status".
	Status RemoteStartTransactionResponseJsonStatus `json:"status" yaml:"status" mapstructure:"status"`
}

const RemoteStartTransactionResponseJsonStatusAccepted RemoteStartTransactionResponseJsonStatus = "Accepted"
const RemoteStartTransactionResponseJsonStatusRejected RemoteStartTransactionResponseJsonStatus = "Rejected"

func (*RemoteStartTransactionResponseJson) IsResponse() {}
