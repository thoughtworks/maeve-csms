// SPDX-License-Identifier: Apache-2.0

package api

import (
	"fmt"
	handlers "github.com/thoughtworks/maeve-csms/manager/handlers/ocpp201"
	"github.com/thoughtworks/maeve-csms/manager/ocpi"
	"net/http"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/render"
	"github.com/thoughtworks/maeve-csms/manager/ocpp"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"k8s.io/utils/clock"
)

type Server struct {
	store   store.Engine
	clock   clock.PassiveClock
	swagger *openapi3.T
	ocpi    ocpi.Api
}

func NewServer(engine store.Engine, clock clock.PassiveClock, ocpi ocpi.Api) (*Server, error) {
	swagger, err := GetSwagger()
	if err != nil {
		return nil, err
	}
	return &Server{
		store:   engine,
		clock:   clock,
		ocpi:    ocpi,
		swagger: swagger,
	}, nil
}

func (s *Server) RegisterChargeStation(w http.ResponseWriter, r *http.Request, csId string) {
	req := new(ChargeStationAuth)
	if err := render.Bind(r, req); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	var pwd string
	if req.Base64SHA256Password != nil {
		pwd = *req.Base64SHA256Password
	}
	err := s.store.SetChargeStationAuth(r.Context(), csId, &store.ChargeStationAuth{
		SecurityProfile:      store.SecurityProfile(req.SecurityProfile),
		Base64SHA256Password: pwd,
	})
	if err != nil {
		_ = render.Render(w, r, ErrInternalError(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) ReconfigureChargeStation(w http.ResponseWriter, r *http.Request, csId string) {
	req := new(ChargeStationSettings)
	if err := render.Bind(r, req); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	chargeStationSettings := make(map[string]*store.ChargeStationSetting)
	for k, v := range *req {
		chargeStationSettings[k] = &store.ChargeStationSetting{
			Value:  v,
			Status: store.ChargeStationSettingStatusPending,
		}
	}

	err := s.store.UpdateChargeStationSettings(r.Context(), csId, &store.ChargeStationSettings{
		Settings: chargeStationSettings,
	})
	if err != nil {
		_ = render.Render(w, r, ErrInternalError(err))
		return
	}
}

func (s *Server) InstallChargeStationCertificates(w http.ResponseWriter, r *http.Request, csId string) {
	req := new(ChargeStationInstallCertificates)
	if err := render.Bind(r, req); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	var certs []*store.ChargeStationInstallCertificate
	for _, cert := range req.Certificates {
		certId, err := handlers.GetCertificateId(cert.Certificate)
		if err != nil {
			_ = render.Render(w, r, ErrInvalidRequest(fmt.Errorf("invalid certificate: %w", err)))
			return
		}

		status := store.CertificateInstallationPending
		if cert.Status != nil {
			status = store.CertificateInstallationStatus(*cert.Status)
		}

		certs = append(certs, &store.ChargeStationInstallCertificate{
			CertificateType:               store.CertificateType(cert.Type),
			CertificateId:                 certId,
			CertificateData:               cert.Certificate,
			CertificateInstallationStatus: status,
		})
	}

	err := s.store.UpdateChargeStationInstallCertificates(r.Context(), csId, &store.ChargeStationInstallCertificates{
		ChargeStationId: csId,
		Certificates:    certs,
	})
	if err != nil {
		_ = render.Render(w, r, ErrInternalError(err))
		return
	}
}

func (s *Server) LookupChargeStationAuth(w http.ResponseWriter, r *http.Request, csId string) {
	auth, err := s.store.LookupChargeStationAuth(r.Context(), csId)
	if err != nil {
		_ = render.Render(w, r, ErrInternalError(err))
		return
	}
	if auth == nil {
		_ = render.Render(w, r, ErrNotFound)
		return
	}

	resp := &ChargeStationAuth{
		SecurityProfile: int(auth.SecurityProfile),
	}
	if auth.Base64SHA256Password != "" {
		resp.Base64SHA256Password = &auth.Base64SHA256Password
	}

	_ = render.Render(w, r, resp)
}

func (s *Server) TriggerChargeStation(w http.ResponseWriter, r *http.Request, csId string) {
	req := new(ChargeStationTrigger)
	if err := render.Bind(r, req); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	err := s.store.SetChargeStationTriggerMessage(r.Context(), csId, &store.ChargeStationTriggerMessage{
		TriggerMessage: store.TriggerMessage(req.Trigger),
		TriggerStatus:  store.TriggerStatusPending,
	})
	if err != nil {
		_ = render.Render(w, r, ErrInternalError(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) SetToken(w http.ResponseWriter, r *http.Request) {
	req := new(Token)
	if err := render.Bind(r, req); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	normContractId, err := ocpp.NormalizeEmaid(req.ContractId)
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	err = s.store.SetToken(r.Context(), &store.Token{
		CountryCode:  req.CountryCode,
		PartyId:      req.PartyId,
		Type:         string(req.Type),
		Uid:          req.Uid,
		ContractId:   normContractId,
		VisualNumber: req.VisualNumber,
		Issuer:       req.Issuer,
		GroupId:      req.GroupId,
		Valid:        req.Valid,
		LanguageCode: req.LanguageCode,
		CacheMode:    string(req.CacheMode),
		LastUpdated:  s.clock.Now().Format(time.RFC3339),
	})
	if err != nil {
		_ = render.Render(w, r, ErrInternalError(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func newToken(tok *store.Token) (*Token, error) {
	lastUpdated, err := time.Parse(time.RFC3339, tok.LastUpdated)
	if err != nil {
		return nil, err
	}

	return &Token{
		CountryCode:  tok.CountryCode,
		PartyId:      tok.PartyId,
		Type:         TokenType(tok.Type),
		Uid:          tok.Uid,
		ContractId:   tok.ContractId,
		VisualNumber: tok.VisualNumber,
		Issuer:       tok.Issuer,
		GroupId:      tok.GroupId,
		Valid:        tok.Valid,
		LanguageCode: tok.LanguageCode,
		CacheMode:    TokenCacheMode(tok.CacheMode),
		LastUpdated:  &lastUpdated,
	}, nil
}

func (s *Server) LookupToken(w http.ResponseWriter, r *http.Request, tokenUid string) {
	tok, err := s.store.LookupToken(r.Context(), tokenUid)
	if err != nil {
		_ = render.Render(w, r, ErrInternalError(err))
		return
	}
	if tok == nil {
		_ = render.Render(w, r, ErrNotFound)
		return
	}

	resp, err := newToken(tok)
	if err != nil {
		_ = render.Render(w, r, ErrInternalError(err))
		return
	}

	_ = render.Render(w, r, resp)
}

func (s *Server) ListTokens(w http.ResponseWriter, r *http.Request, params ListTokensParams) {
	offset := 0
	limit := 20

	if params.Offset != nil {
		offset = *params.Offset
	}
	if params.Limit != nil {
		limit = *params.Limit
	}
	if limit > 100 {
		limit = 100
	}

	tokens, err := s.store.ListTokens(r.Context(), offset, limit)
	if err != nil {
		_ = render.Render(w, r, ErrInternalError(err))
		return
	}

	var resp = make([]render.Renderer, len(tokens))
	for i, tok := range tokens {
		resp[i], err = newToken(tok)
		if err != nil {
			_ = render.Render(w, r, ErrInternalError(err))
			return
		}
	}
	_ = render.RenderList(w, r, resp)
}

func (s *Server) UploadCertificate(w http.ResponseWriter, r *http.Request) {
	req := new(Certificate)
	if err := render.Bind(r, req); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	err := s.store.SetCertificate(r.Context(), req.Certificate)
	if err != nil {
		_ = render.Render(w, r, ErrInternalError(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) DeleteCertificate(w http.ResponseWriter, r *http.Request, certificateHash string) {
	err := s.store.DeleteCertificate(r.Context(), certificateHash)
	if err != nil {
		_ = render.Render(w, r, ErrInternalError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) LookupCertificate(w http.ResponseWriter, r *http.Request, certificateHash string) {
	cert, err := s.store.LookupCertificate(r.Context(), certificateHash)
	if err != nil {
		_ = render.Render(w, r, ErrInternalError(err))
		return
	}
	if cert == "" {
		_ = render.Render(w, r, ErrNotFound)
		return
	}

	resp := &Certificate{
		Certificate: cert,
	}
	_ = render.Render(w, r, resp)
}

func (s *Server) RegisterParty(w http.ResponseWriter, r *http.Request) {
	if s.ocpi == nil {
		_ = render.Render(w, r, ErrNotFound)
		return
	}

	req := new(Registration)
	if err := render.Bind(r, req); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if req.Url != nil {
		err := s.ocpi.RegisterNewParty(r.Context(), *req.Url, req.Token)
		if err != nil {
			_ = render.Render(w, r, ErrInternalError(err))
			return
		}
	} else {
		// store credentials in database
		status := store.OcpiRegistrationStatusPending
		if req.Status != nil && *req.Status == "REGISTERED" {
			status = store.OcpiRegistrationStatusRegistered
		}

		err := s.store.SetRegistrationDetails(r.Context(), req.Token, &store.OcpiRegistration{
			Status: status,
		})
		if err != nil {
			_ = render.Render(w, r, ErrInternalError(err))
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) RegisterLocation(w http.ResponseWriter, r *http.Request, locationId string) {
	if s.ocpi == nil {
		_ = render.Render(w, r, ErrNotFound)
		return
	}

	req := new(Location)
	if err := render.Bind(r, req); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	now := s.clock.Now()

	var numEvses int
	if req.Evses != nil {
		numEvses = len(*req.Evses)
	}
	storeEvses := make([]store.Evse, numEvses)
	if numEvses != 0 {
		for i, reqEvse := range *req.Evses {
			storeConnectors := make([]store.Connector, len(reqEvse.Connectors))
			for j, reqConnector := range reqEvse.Connectors {
				storeConnectors[j] = store.Connector{
					Id:          reqConnector.Id,
					Format:      string(reqConnector.Format),
					PowerType:   string(reqConnector.PowerType),
					Standard:    string(reqConnector.Standard),
					MaxVoltage:  reqConnector.MaxVoltage,
					MaxAmperage: reqConnector.MaxAmperage,
					LastUpdated: now.Format(time.RFC3339),
				}
				storeEvses[i] = store.Evse{
					Connectors:  storeConnectors,
					EvseId:      reqEvse.EvseId,
					Status:      string(ocpi.EvseStatusUNKNOWN),
					Uid:         reqEvse.Uid,
					LastUpdated: now.Format(time.RFC3339),
				}
			}
		}
	}
	err := s.store.SetLocation(r.Context(), &store.Location{
		Address: req.Address,
		City:    req.City,
		Coordinates: store.GeoLocation{
			Latitude:  req.Coordinates.Latitude,
			Longitude: req.Coordinates.Longitude,
		},
		Country:     req.Country,
		Evses:       &storeEvses,
		Id:          locationId,
		Name:        *req.Name,
		ParkingType: string(*req.ParkingType),
		PostalCode:  *req.PostalCode,
	})
	if err != nil {
		_ = render.Render(w, r, ErrInternalError(err))
		return
	}

	ocpiEvses := make([]ocpi.Evse, numEvses)
	if numEvses != 0 {
		for i, reqEvse := range *req.Evses {
			ocpiConnectors := make([]ocpi.Connector, len(reqEvse.Connectors))
			for j, reqConnector := range reqEvse.Connectors {
				ocpiConnectors[j] = ocpi.Connector{
					Id:          reqConnector.Id,
					Format:      ocpi.ConnectorFormat(reqConnector.Format),
					PowerType:   ocpi.ConnectorPowerType(reqConnector.PowerType),
					Standard:    ocpi.ConnectorStandard(reqConnector.Standard),
					MaxVoltage:  reqConnector.MaxVoltage,
					MaxAmperage: reqConnector.MaxAmperage,
					LastUpdated: now.Format(time.RFC3339),
				}
				ocpiEvses[i] = ocpi.Evse{
					Connectors:  ocpiConnectors,
					EvseId:      reqEvse.EvseId,
					Status:      ocpi.EvseStatusUNKNOWN,
					Uid:         reqEvse.Uid,
					LastUpdated: now.Format(time.RFC3339),
				}
			}
		}
	}
	err = s.ocpi.PushLocation(r.Context(), ocpi.Location{
		Address: req.Address,
		City:    req.City,
		Coordinates: ocpi.GeoLocation{
			Latitude:  req.Coordinates.Latitude,
			Longitude: req.Coordinates.Longitude,
		},
		Country:     req.Country,
		CountryCode: req.CountryCode,
		Evses:       &ocpiEvses,
		Id:          locationId,
		LastUpdated: now.Format(time.RFC3339),
		Name:        req.Name,
		ParkingType: (*ocpi.LocationParkingType)(req.ParkingType),
		PartyId:     req.PartyId,
		PostalCode:  req.PostalCode,
		Publish:     true,
	})
	if err != nil {
		_ = render.Render(w, r, ErrInternalError(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
}
