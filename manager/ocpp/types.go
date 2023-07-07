// SPDX-License-Identifier: Apache-2.0

package ocpp

type Request interface {
	IsRequest()
}

type Response interface {
	IsResponse()
}
