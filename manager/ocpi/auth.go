// SPDX-License-Identifier: Apache-2.0

package ocpi

import (
	"context"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"net/http"
	"regexp"
)

var authzHeaderRegexp = regexp.MustCompile(`^Token (.*)$`)

func NewTokenAuthenticationFunc(engine store.Engine) openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		authzHeader := input.RequestValidationInput.Request.Header.Get("Authorization")
		matches := authzHeaderRegexp.FindStringSubmatch(authzHeader)
		if len(matches) != 2 {
			return input.NewError(nil)
		}

		reg, err := engine.GetRegistrationDetails(ctx, matches[1])
		if err != nil {
			return input.NewError(err)
		}
		if reg == nil {
			return input.NewError(fmt.Errorf("unknown token"))
		}
		if reg.Status != store.OcpiRegistrationStatusRegistered {
			allowed := false
			switch input.RequestValidationInput.Request.Method {
			case http.MethodGet:
				switch input.RequestValidationInput.Request.URL.Path {
				case "/ocpi/versions":
					allowed = true
				case "/ocpi/2.2":
					allowed = true
				}
			case http.MethodPost:
				switch input.RequestValidationInput.Request.URL.Path {
				case "/ocpi/2.2/credentials":
					allowed = true
				}
			}

			if !allowed {
				return input.NewError(fmt.Errorf("unregistered token"))
			}
		}

		return nil
	}
}
