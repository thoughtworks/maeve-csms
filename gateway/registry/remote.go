package registry

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type RemoteRegistry struct {
	ManagerApiAddr string
}

type ChargeStationAuthDetailsResponse struct {
	SecurityProfile        int    `json:"securityProfile"`
	Base64SHA256Password   string `json:"base64SHA256Password,omitempty"`
	InvalidUsernameAllowed bool   `json:"invalidUsernameAllowed,omitempty"`
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
			ClientId:               clientId,
			SecurityProfile:        SecurityProfile(chargeStationAuthDetails.SecurityProfile),
			Base64SHA256Password:   chargeStationAuthDetails.Base64SHA256Password,
			InvalidUsernameAllowed: chargeStationAuthDetails.InvalidUsernameAllowed,
		}, nil
	}

	return nil, nil
}

type CertificateResponse struct {
	Certificate string `json:"certificate"`
}

func (r RemoteRegistry) LookupCertificate(certHash string) (*x509.Certificate, error) {
	certHash = strings.Replace(certHash, "/", "_", -1)
	certHash = strings.Replace(certHash, "+", "-", -1)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v0/certificate/%s", r.ManagerApiAddr, certHash), nil)
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

		var certificateResponse CertificateResponse
		err = json.Unmarshal(b, &certificateResponse)
		if err != nil {
			return nil, fmt.Errorf("unmarshaling data: %w", err)
		}

		block, _ := pem.Decode([]byte(certificateResponse.Certificate))
		if block == nil {
			return nil, fmt.Errorf("no pem data found")
		}
		if block.Type == "CERTIFICATE" {
			return x509.ParseCertificate(block.Bytes)
		} else {
			return nil, fmt.Errorf("no certificate found in PEM data")
		}
	}

	return nil, nil
}
