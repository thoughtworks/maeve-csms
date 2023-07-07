package api

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/twlabs/maeve-csms/manager/store"
	"net/http"
)

type ChargeStationManager struct {
	ChargeStationAuthStore store.ChargeStationAuthStore
}

type CreateChargeStationRequest struct {
	SecurityProfile      int    `json:"securityProfile"`
	Base64SHA256Password string `json:"base64SHA256Password"`
}

type ChargeStationAuthDetailsResponse struct {
	SecurityProfile      int    `json:"securityProfile"`
	Base64SHA256Password string `json:"base64SHA256Password,omitempty"`
}

func (ccs *CreateChargeStationRequest) Bind(r *http.Request) error {
	return nil
}

func (c ChargeStationAuthDetailsResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (csm *ChargeStationManager) CreateChargeStation(w http.ResponseWriter, r *http.Request) {
	csId := chi.URLParam(r, "csId")

	req := new(CreateChargeStationRequest)
	if err := render.Bind(r, req); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	profile, err := convertToSecurityProfile(req)
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	err = csm.ChargeStationAuthStore.SetChargeStationAuth(r.Context(), csId, &store.ChargeStationAuth{
		SecurityProfile:      profile,
		Base64SHA256Password: req.Base64SHA256Password,
	})
	if err != nil {
		_ = render.Render(w, r, ErrInternalError(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (csm *ChargeStationManager) RetrieveChargeStationAuthDetails(w http.ResponseWriter, r *http.Request) {
	csId := chi.URLParam(r, "csId")

	auth, err := csm.ChargeStationAuthStore.LookupChargeStationAuth(r.Context(), csId)
	if err != nil {
		_ = render.Render(w, r, ErrInternalError(err))
		return
	}
	if auth == nil {
		_ = render.Render(w, r, ErrNotFound)
		return
	}

	_ = render.Render(w, r, &ChargeStationAuthDetailsResponse{
		SecurityProfile:      int(auth.SecurityProfile),
		Base64SHA256Password: auth.Base64SHA256Password,
	})
}

func convertToSecurityProfile(req *CreateChargeStationRequest) (store.SecurityProfile, error) {
	var profile store.SecurityProfile
	if req.SecurityProfile < 0 || req.SecurityProfile > 2 {
		return profile, fmt.Errorf("unknown security profile %d", req.SecurityProfile)
	}
	return store.SecurityProfile(req.SecurityProfile), nil
}
