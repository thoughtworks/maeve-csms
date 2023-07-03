package ocpp

type Request interface {
	IsRequest()
}

type Response interface {
	IsResponse()
}
