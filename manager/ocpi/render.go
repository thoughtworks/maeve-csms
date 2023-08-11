package ocpi

import "net/http"

func (OcpiResponseListVersion) Render(http.ResponseWriter, *http.Request) error {
	return nil
}

func (OcpiResponseVersionDetail) Render(http.ResponseWriter, *http.Request) error {
	return nil
}

func (Credentials) Bind(r *http.Request) error {
	return nil
}
