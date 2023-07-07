// SPDX-License-Identifier: Apache-2.0

package ocpp16

type BootNotificationResponseJsonStatus string

type BootNotificationResponseJson struct {
	// CurrentTime corresponds to the JSON schema field "currentTime".
	CurrentTime string `json:"currentTime" yaml:"currentTime" mapstructure:"currentTime"`

	// Interval corresponds to the JSON schema field "interval".
	Interval int `json:"interval" yaml:"interval" mapstructure:"interval"`

	// Status corresponds to the JSON schema field "status".
	Status BootNotificationResponseJsonStatus `json:"status" yaml:"status" mapstructure:"status"`
}

const BootNotificationResponseJsonStatusAccepted BootNotificationResponseJsonStatus = "Accepted"
const BootNotificationResponseJsonStatusPending BootNotificationResponseJsonStatus = "Pending"
const BootNotificationResponseJsonStatusRejected BootNotificationResponseJsonStatus = "Rejected"

func (*BootNotificationResponseJson) IsResponse() {}
