// SPDX-License-Identifier: Apache-2.0

package ocpi

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
	"k8s.io/utils/clock"
	"net/http"
	"time"
)

type Server struct {
	ocpi  Api
	clock clock.PassiveClock
}

func NewServer(ocpi Api, clock clock.PassiveClock) (*Server, error) {
	return &Server{
		ocpi:  ocpi,
		clock: clock,
	}, nil
}

// VERSIONS

func (s *Server) GetVersions(w http.ResponseWriter, r *http.Request, params GetVersionsParams) {
	versions, err := s.ocpi.GetVersions(r.Context())
	if err != nil {
		_ = render.Render(w, r, OcpiResponseListVersion{
			StatusCode: StatusGenericServerFailure,
			Timestamp:  s.clock.Now().Format(time.RFC3339),
		})
		return
	}

	_ = render.Render(w, r, OcpiResponseListVersion{
		Data:          &versions,
		StatusCode:    StatusSuccess,
		StatusMessage: &StatusSuccessMessage,
		Timestamp:     s.clock.Now().Format(time.RFC3339),
	})
}

func (s *Server) GetVersion(w http.ResponseWriter, r *http.Request, params GetVersionParams) {
	version, err := s.ocpi.GetVersion(r.Context())
	if err != nil {
		if err != nil {
			_ = render.Render(w, r, OcpiResponseListVersion{
				StatusCode: StatusGenericServerFailure,
				Timestamp:  s.clock.Now().Format(time.RFC3339),
			})
			return
		}
	}
	_ = render.Render(w, r, OcpiResponseVersionDetail{
		StatusCode:    StatusSuccess,
		Data:          &version,
		Timestamp:     s.clock.Now().Format(time.RFC3339),
		StatusMessage: &StatusSuccessMessage,
	})
}

// CREDENTIALS

func (s *Server) PostCredentials(w http.ResponseWriter, r *http.Request, params PostCredentialsParams) {
	creds := new(Credentials)
	if err := render.Bind(r, creds); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	matches := authzHeaderRegexp.FindStringSubmatch(params.Authorization)
	if len(matches) != 2 {
		_ = render.Render(w, r, ErrInvalidRequest(fmt.Errorf("invalid authorization header")))
		return
	}

	err := s.ocpi.SetCredentials(r.Context(), matches[1], *creds)
	if err != nil {
		slog.Error("Error setting credentials", "err", err)
		_ = render.Render(w, r, ErrInternalError(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// TOKEN RECEIVER

func (s *Server) GetClientOwnedToken(w http.ResponseWriter, r *http.Request, countryCode string, partyID string, tokenUID string, params GetClientOwnedTokenParams) {
	token, err := s.ocpi.GetToken(r.Context(), countryCode, partyID, tokenUID)
	if err != nil {
		_ = render.Render(w, r, ErrInternalError(err))
		return
	}

	_ = render.Render(w, r, OcpiResponseToken{
		StatusCode:    StatusSuccess,
		StatusMessage: &StatusSuccessMessage,
		Timestamp:     s.clock.Now().Format(time.RFC3339),
		Data:          token,
	})
}

func (s *Server) PutClientOwnedToken(w http.ResponseWriter, r *http.Request, countryCode string, partyID string, tokenUID string, params PutClientOwnedTokenParams) {
	tok := new(Token)
	if err := render.Bind(r, tok); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if tok.CountryCode != countryCode {
		_ = render.Render(w, r, ErrInvalidRequest(fmt.Errorf("token country code mismatch")))
		return
	}
	if tok.PartyId != partyID {
		_ = render.Render(w, r, ErrInvalidRequest(fmt.Errorf("token party id mismatch")))
		return
	}
	if tok.Uid != tokenUID {
		_ = render.Render(w, r, ErrInvalidRequest(fmt.Errorf("token uid mismatch")))
		return
	}

	err := s.ocpi.SetToken(r.Context(), *tok)
	if err != nil {
		_ = render.Render(w, r, OcpiResponseListVersion{
			StatusCode: StatusGenericServerFailure,
			Timestamp:  s.clock.Now().Format(time.RFC3339),
		})
	}
}

func (s *Server) PatchClientOwnedToken(w http.ResponseWriter, r *http.Request, countryCode string, partyID string, tokenUID string, params PatchClientOwnedTokenParams) {
	var patch map[string]any
	err := json.NewDecoder(r.Body).Decode(&patch)
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	tok, err := s.ocpi.GetToken(r.Context(), countryCode, partyID, tokenUID)
	if err != nil {
		_ = render.Render(w, r, ErrInternalError(err))
		return
	}
	if tok == nil {
		_ = render.Render(w, r, ErrNotFound)
		return
	}

	for k, v := range patch {
		switch k {
		case "contract_id":
			contractID := v.(string)
			tok.ContractId = contractID
		case "group_id":
			groupID := v.(string)
			tok.GroupId = &groupID
		case "issuer":
			issuer := v.(string)
			tok.Issuer = issuer
		case "language":
			language := v.(string)
			tok.Language = &language
		case "type":
			typ := v.(string)
			tok.Type = TokenType(typ)
		case "valid":
			valid := v.(bool)
			tok.Valid = valid
		case "visual_number":
			visualNumber := v.(string)
			tok.VisualNumber = &visualNumber
		case "whitelist":
			whitelist := v.(string)
			tok.Whitelist = TokenWhitelist(whitelist)
		default:
			_ = render.Render(w, r, ErrInvalidRequest(fmt.Errorf("unknown field %s", k)))
		}
	}

	err = s.ocpi.SetToken(r.Context(), *tok)
	if err != nil {
		_ = render.Render(w, r, ErrInternalError(err))
		return
	}
}

func (s *Server) DeleteCredentials(w http.ResponseWriter, r *http.Request, params DeleteCredentialsParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) GetCredentials(w http.ResponseWriter, r *http.Request, params GetCredentialsParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) PutCredentials(w http.ResponseWriter, r *http.Request, params PutCredentialsParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) DeleteReceiverChargingProfile(w http.ResponseWriter, r *http.Request, sessionId string, params DeleteReceiverChargingProfileParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) GetReceiverChargingProfile(w http.ResponseWriter, r *http.Request, sessionId string, params GetReceiverChargingProfileParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) PutReceiverChargingProfile(w http.ResponseWriter, r *http.Request, sessionId string, params PutReceiverChargingProfileParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) PostGenericChargingProfileResult(w http.ResponseWriter, r *http.Request, uid string, params PostGenericChargingProfileResultParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) PutSenderChargingProfile(w http.ResponseWriter, r *http.Request, sessionId string, params PutSenderChargingProfileParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) PostClientOwnedCdr(w http.ResponseWriter, r *http.Request, params PostClientOwnedCdrParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) GetClientOwnedCdr(w http.ResponseWriter, r *http.Request, cdrID string, params GetClientOwnedCdrParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) PostCancelReservation(w http.ResponseWriter, r *http.Request, params PostCancelReservationParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) PostReserveNow(w http.ResponseWriter, r *http.Request, params PostReserveNowParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) PostStartSession(w http.ResponseWriter, r *http.Request, params PostStartSessionParams) {
	_ = render.Render(w, r, OcpiResponseCommandResponse{
		StatusCode:    StatusSuccess,
		StatusMessage: &StatusSuccessMessage,
		Timestamp:     s.clock.Now().Format(time.RFC3339),
		Data:          &CommandResponse{Result: "ACCEPTED"},
	})
}

func (s *Server) PostStopSession(w http.ResponseWriter, r *http.Request, params PostStopSessionParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) PostUnlockConnector(w http.ResponseWriter, r *http.Request, params PostUnlockConnectorParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) GetClientOwnedLocation(w http.ResponseWriter, r *http.Request, countryCode string, partyID string, locationID string, params GetClientOwnedLocationParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) PatchClientOwnedLocation(w http.ResponseWriter, r *http.Request, countryCode string, partyID string, locationID string, params PatchClientOwnedLocationParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) PutClientOwnedLocation(w http.ResponseWriter, r *http.Request, countryCode string, partyID string, locationID string, params PutClientOwnedLocationParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) GetClientOwnedEvse(w http.ResponseWriter, r *http.Request, countryCode string, partyID string, locationID string, evseUID string, params GetClientOwnedEvseParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) PatchClientOwnedEvse(w http.ResponseWriter, r *http.Request, countryCode string, partyID string, locationID string, evseUID string, params PatchClientOwnedEvseParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) PutClientOwnedEvse(w http.ResponseWriter, r *http.Request, countryCode string, partyID string, locationID string, evseUID string, params PutClientOwnedEvseParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) GetClientOwnedConnector(w http.ResponseWriter, r *http.Request, countryCode string, partyID string, locationID string, evseUID string, connectorID string, params GetClientOwnedConnectorParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) PatchClientOwnedConnector(w http.ResponseWriter, r *http.Request, countryCode string, partyID string, locationID string, evseUID string, connectorID string, params PatchClientOwnedConnectorParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) PutClientOwnedConnector(w http.ResponseWriter, r *http.Request, countryCode string, partyID string, locationID string, evseUID string, connectorID string, params PutClientOwnedConnectorParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) GetClientOwnedSession(w http.ResponseWriter, r *http.Request, countryCode string, partyID string, sessionID string, params GetClientOwnedSessionParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) PatchClientOwnedSession(w http.ResponseWriter, r *http.Request, countryCode string, partyID string, sessionID string, params PatchClientOwnedSessionParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) PutClientOwnedSession(w http.ResponseWriter, r *http.Request, countryCode string, partyID string, sessionID string, params PutClientOwnedSessionParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) DeleteClientOwnedTariff(w http.ResponseWriter, r *http.Request, countryCode string, partyID string, tariffID string, params DeleteClientOwnedTariffParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) GetClientOwnedTariff(w http.ResponseWriter, r *http.Request, countryCode string, partyID string, tariffID string, params GetClientOwnedTariffParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) PutClientOwnedTariff(w http.ResponseWriter, r *http.Request, countryCode string, partyID string, tariffID string, params PutClientOwnedTariffParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) GetCdrsFromDataOwner(w http.ResponseWriter, r *http.Request, params GetCdrsFromDataOwnerParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) GetCdrPageFromDataOwner(w http.ResponseWriter, r *http.Request, uid string, params GetCdrPageFromDataOwnerParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) PostAsyncResponse(w http.ResponseWriter, r *http.Request, command PostAsyncResponseParamsCommand, uid string, params PostAsyncResponseParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) GetLocationListFromDataOwner(w http.ResponseWriter, r *http.Request, params GetLocationListFromDataOwnerParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) GetLocationPageFromDataOwner(w http.ResponseWriter, r *http.Request, uid string, params GetLocationPageFromDataOwnerParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) GetLocationObjectFromDataOwner(w http.ResponseWriter, r *http.Request, locationID string, params GetLocationObjectFromDataOwnerParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) GetEvseObjectFromDataOwner(w http.ResponseWriter, r *http.Request, locationID string, evseUID string, params GetEvseObjectFromDataOwnerParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) GetConnectorObjectFromDataOwner(w http.ResponseWriter, r *http.Request, locationID string, evseUID string, connectorID string, params GetConnectorObjectFromDataOwnerParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) GetSessionsFromDataOwner(w http.ResponseWriter, r *http.Request, params GetSessionsFromDataOwnerParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) GetSessionsPageFromDataOwner(w http.ResponseWriter, r *http.Request, uid string, params GetSessionsPageFromDataOwnerParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) PutChargingPreferences(w http.ResponseWriter, r *http.Request, sessionID string, params PutChargingPreferencesParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) GetTariffsFromDataOwner(w http.ResponseWriter, r *http.Request, params GetTariffsFromDataOwnerParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) GetTariffsPageFromDataOwner(w http.ResponseWriter, r *http.Request, uid string, params GetTariffsPageFromDataOwnerParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) GetTokensFromDataOwner(w http.ResponseWriter, r *http.Request, params GetTokensFromDataOwnerParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) GetTokensPageFromDataOwner(w http.ResponseWriter, r *http.Request, uid string, params GetTokensPageFromDataOwnerParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s *Server) PostRealTimeTokenAuthorization(w http.ResponseWriter, r *http.Request, tokenUID string, params PostRealTimeTokenAuthorizationParams) {
	w.WriteHeader(http.StatusNotImplemented)
}
