// SPDX-License-Identifier: Apache-2.0

package ocpp201

type FirmwareStatusEnumType string

const FirmwareStatusEnumTypeDownloadFailed FirmwareStatusEnumType = "DownloadFailed"
const FirmwareStatusEnumTypeDownloadPaused FirmwareStatusEnumType = "DownloadPaused"
const FirmwareStatusEnumTypeDownloadScheduled FirmwareStatusEnumType = "DownloadScheduled"
const FirmwareStatusEnumTypeDownloaded FirmwareStatusEnumType = "Downloaded"
const FirmwareStatusEnumTypeDownloading FirmwareStatusEnumType = "Downloading"
const FirmwareStatusEnumTypeIdle FirmwareStatusEnumType = "Idle"
const FirmwareStatusEnumTypeInstallRebooting FirmwareStatusEnumType = "InstallRebooting"
const FirmwareStatusEnumTypeInstallScheduled FirmwareStatusEnumType = "InstallScheduled"
const FirmwareStatusEnumTypeInstallVerificationFailed FirmwareStatusEnumType = "InstallVerificationFailed"
const FirmwareStatusEnumTypeInstallationFailed FirmwareStatusEnumType = "InstallationFailed"
const FirmwareStatusEnumTypeInstalled FirmwareStatusEnumType = "Installed"
const FirmwareStatusEnumTypeInstalling FirmwareStatusEnumType = "Installing"
const FirmwareStatusEnumTypeInvalidSignature FirmwareStatusEnumType = "InvalidSignature"
const FirmwareStatusEnumTypeSignatureVerified FirmwareStatusEnumType = "SignatureVerified"

type FirmwareStatusNotificationRequestJson struct {
	// CustomData corresponds to the JSON schema field "customData".
	CustomData *CustomDataType `json:"customData,omitempty" yaml:"customData,omitempty" mapstructure:"customData,omitempty"`

	// The request id that was provided in the
	// UpdateFirmwareRequest that started this firmware update.
	// This field is mandatory, unless the message was triggered by a
	// TriggerMessageRequest AND there is no firmware update ongoing.
	//
	RequestId *int `json:"requestId,omitempty" yaml:"requestId,omitempty" mapstructure:"requestId,omitempty"`

	// Status corresponds to the JSON schema field "status".
	Status FirmwareStatusEnumType `json:"status" yaml:"status" mapstructure:"status"`
}

func (*FirmwareStatusNotificationRequestJson) IsRequest() {}
