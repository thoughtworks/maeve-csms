package api

import "net/http"

func (c ChargeStationAuth) Bind(r *http.Request) error {
	return nil
}

func (c ChargeStationAuth) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
