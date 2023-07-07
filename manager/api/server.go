package api

import (
	"fmt"
	"github.com/go-chi/render"
	"github.com/twlabs/maeve-csms/manager/store"
	"net/http"
)

type Server struct {
	Store store.Engine
}

func (s *Server) RegisterChargeStation(w http.ResponseWriter, r *http.Request, csId string) {
	req := new(ChargeStationAuth)
	if err := render.Bind(r, req); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	profile, err := convertToSecurityProfile(req.SecurityProfile)
	if err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	var pwd string
	if req.Base64SHA256Password != nil {
		pwd = *req.Base64SHA256Password
	}
	err = s.Store.SetChargeStationAuth(r.Context(), csId, &store.ChargeStationAuth{
		SecurityProfile:      profile,
		Base64SHA256Password: pwd,
	})
	if err != nil {
		_ = render.Render(w, r, ErrInternalError(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) LookupChargeStationAuth(w http.ResponseWriter, r *http.Request, csId string) {
	auth, err := s.Store.LookupChargeStationAuth(r.Context(), csId)
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

func convertToSecurityProfile(securityProfile int) (store.SecurityProfile, error) {
	var profile store.SecurityProfile
	if securityProfile < 0 || securityProfile > 2 {
		return profile, fmt.Errorf("unknown security profile %d", securityProfile)
	}
	return store.SecurityProfile(securityProfile), nil
}
