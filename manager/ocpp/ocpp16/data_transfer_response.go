package ocpp16

type DataTransferResponseJsonStatus string

var enumValues_DataTransferResponseJsonStatus = []interface{}{
	"Accepted",
	"Rejected",
	"UnknownMessageId",
	"UnknownVendorId",
}

type DataTransferResponseJson struct {
	// Data corresponds to the JSON schema field "data".
	Data *string `json:"data,omitempty" yaml:"data,omitempty" mapstructure:"data,omitempty"`

	// Status corresponds to the JSON schema field "status".
	Status DataTransferResponseJsonStatus `json:"status" yaml:"status" mapstructure:"status"`
}

const DataTransferResponseJsonStatusAccepted DataTransferResponseJsonStatus = "Accepted"
const DataTransferResponseJsonStatusRejected DataTransferResponseJsonStatus = "Rejected"
const DataTransferResponseJsonStatusUnknownMessageId DataTransferResponseJsonStatus = "UnknownMessageId"
const DataTransferResponseJsonStatusUnknownVendorId DataTransferResponseJsonStatus = "UnknownVendorId"

func (*DataTransferResponseJson) IsResponse() {}
