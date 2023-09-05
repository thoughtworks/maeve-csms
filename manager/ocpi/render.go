// SPDX-License-Identifier: Apache-2.0

package ocpi

import "net/http"

func (OcpiResponseCommandResponse) Render(http.ResponseWriter, *http.Request) error {
	return nil
}

func (OcpiResponseListVersion) Render(http.ResponseWriter, *http.Request) error {
	return nil
}

func (OcpiResponseVersionDetail) Render(http.ResponseWriter, *http.Request) error {
	return nil
}

func (OcpiResponseToken) Render(http.ResponseWriter, *http.Request) error {
	return nil
}

func (Credentials) Bind(r *http.Request) error {
	return nil
}

func (Token) Bind(r *http.Request) error {
	return nil
}
