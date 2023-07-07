package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type RemoteRegistry struct {
	ManagerApiAddr string
}

type ChargeStationAuthDetailsResponse struct {
	SecurityProfile      int    `json:"securityProfile"`
	Base64SHA256Password string `json:"base64SHA256Password,omitempty"`
}

func (r RemoteRegistry) LookupChargeStation(clientId string) (*ChargeStation, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v0/cs/%s/auth", r.ManagerApiAddr, clientId), nil)
	if err != nil {
		return nil, fmt.Errorf("creating http request: %w", err)
	}
	req.Header.Set("accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making http request: %w", err)
	}

	if resp.StatusCode == http.StatusOK {
		var b []byte
		b, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("reading http body: %w", err)
		}
		defer func() {
			_ = resp.Body.Close()
		}()
		var chargeStationAuthDetails ChargeStationAuthDetailsResponse
		err = json.Unmarshal(b, &chargeStationAuthDetails)
		if err != nil {
			return nil, fmt.Errorf("unmarshaling data: %w", err)
		}
		return &ChargeStation{
			ClientId:             clientId,
			SecurityProfile:      SecurityProfile(chargeStationAuthDetails.SecurityProfile),
			Base64SHA256Password: chargeStationAuthDetails.Base64SHA256Password,
		}, nil
	}

	return nil, nil
}
