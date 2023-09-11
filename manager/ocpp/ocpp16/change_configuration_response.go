// SPDX-License-Identifier: Apache-2.0

package ocpp16

type ChangeConfigurationResponseJsonStatus string

type ChangeConfigurationResponseJson struct {
	// Status corresponds to the JSON schema field "status".
	Status ChangeConfigurationResponseJsonStatus `json:"status" yaml:"status" mapstructure:"status"`
}

const ChangeConfigurationResponseJsonStatusAccepted ChangeConfigurationResponseJsonStatus = "Accepted"
const ChangeConfigurationResponseJsonStatusNotSupported ChangeConfigurationResponseJsonStatus = "NotSupported"
const ChangeConfigurationResponseJsonStatusRebootRequired ChangeConfigurationResponseJsonStatus = "RebootRequired"
const ChangeConfigurationResponseJsonStatusRejected ChangeConfigurationResponseJsonStatus = "Rejected"

func (*ChangeConfigurationResponseJson) IsResponse() {}
