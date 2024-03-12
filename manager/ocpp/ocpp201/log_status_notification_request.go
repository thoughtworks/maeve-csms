// SPDX-License-Identifier: Apache-2.0

package ocpp201

type LogStatusNotificationRequestJson struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// The request id that was provided in GetLogRequest that started this log upload.
	// This field is mandatory,
	// unless the message was triggered by a TriggerMessageRequest AND there is no log
	// upload ongoing.
	//
	RequestId *int `json:"requestId,omitempty" yaml:"requestId,omitempty" mapstructure:"requestId,omitempty"`

	// Status corresponds to the JSON schema field "status".
	Status UploadLogStatusEnumType `json:"status" yaml:"status" mapstructure:"status"`
}

func (*LogStatusNotificationRequestJson) IsRequest() {}

type UploadLogStatusEnumType string

const UploadLogStatusEnumTypeAcceptedCanceled UploadLogStatusEnumType = "AcceptedCanceled"
const UploadLogStatusEnumTypeBadMessage UploadLogStatusEnumType = "BadMessage"
const UploadLogStatusEnumTypeIdle UploadLogStatusEnumType = "Idle"
const UploadLogStatusEnumTypeNotSupportedOperation UploadLogStatusEnumType = "NotSupportedOperation"
const UploadLogStatusEnumTypePermissionDenied UploadLogStatusEnumType = "PermissionDenied"
const UploadLogStatusEnumTypeUploadFailure UploadLogStatusEnumType = "UploadFailure"
const UploadLogStatusEnumTypeUploaded UploadLogStatusEnumType = "Uploaded"
const UploadLogStatusEnumTypeUploading UploadLogStatusEnumType = "Uploading"
