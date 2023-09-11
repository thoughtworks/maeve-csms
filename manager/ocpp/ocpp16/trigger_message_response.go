// SPDX-License-Identifier: Apache-2.0

package ocpp16

type TriggerMessageResponseJsonStatus string

type TriggerMessageResponseJson struct {
	// Status corresponds to the JSON schema field "status".
	Status TriggerMessageResponseJsonStatus `json:"status" yaml:"status" mapstructure:"status"`
}

const TriggerMessageResponseJsonStatusAccepted TriggerMessageResponseJsonStatus = "Accepted"
const TriggerMessageResponseJsonStatusNotImplemented TriggerMessageResponseJsonStatus = "NotImplemented"
const TriggerMessageResponseJsonStatusRejected TriggerMessageResponseJsonStatus = "Rejected"

func (*TriggerMessageResponseJson) IsResponse() {}
