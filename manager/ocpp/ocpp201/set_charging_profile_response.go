// SPDX-License-Identifier: Apache-2.0

package ocpp201

import "github.com/thoughtworks/maeve-csms/manager/ocpp/ocpp201"

type ChargingProfileStatusEnumType string

const (
	Accepted ChargingProfileStatusEnumType = "Accepted"
	Rejected ChargingProfileStatusEnumType = "Rejected"
)

type SetChargingProfileResponseJson struct {
	Status     ChargingProfileStatusEnumType
	StatusInfo ocpp201.StatusInfoType
}

func (*SetChargingProfileResponseJson) IsResponse() {}
