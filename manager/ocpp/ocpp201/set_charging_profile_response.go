// SPDX-License-Identifier: Apache-2.0

package ocpp201

type ChargingProfileStatusEnumType string

const (
	Accepted ChargingProfileStatusEnumType = "Accepted"
	Rejected ChargingProfileStatusEnumType = "Rejected"
)

type SetChargingProfileResponseJson struct {
	Status     ChargingProfileStatusEnumType
	StatusInfo StatusInfoType
}

func (*SetChargingProfileResponseJson) IsResponse() {}
