package ocpi

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/thoughtworks/maeve-csms/manager/store"
	"io"
	"math/big"
	"net/http"
)

func (o *OCPI) RegisterNewParty(ctx context.Context, url, token string) error {
	versions, err := o.getVersions(ctx, url, token)
	if err != nil {
		return err
	}

	endpointUrl, err := getEndpointUrl(versions)
	if err != nil {
		return err
	}
	endpoints, err := o.getEndpoints(ctx, endpointUrl, token)
	if err != nil {
		return err
	}

	newToken, err := generateRandomString()
	if err != nil {
		return err
	}

	err = o.store.SetRegistrationDetails(ctx, newToken, &store.OcpiRegistration{
		Status: store.OcpiRegistrationStatusRegistered,
	})
	if err != nil {
		return err
	}

	credentialsUrl, err := getCredentialsUrl(endpoints)
	if err != nil {
		return err
	}

	err = o.postCredentials(ctx, credentialsUrl, token, newToken)
	if err != nil {
		return err
	}

	return nil
}

func (o *OCPI) getVersions(ctx context.Context, url, token string) ([]Version, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var versionList OcpiResponseListVersion
	err = json.Unmarshal(b, &versionList)
	if err != nil {
		return nil, err
	}
	if versionList.StatusCode != StatusSuccess {
		return nil, fmt.Errorf("status code: %d", versionList.StatusCode)
	}
	return *versionList.Data, nil
}

func getEndpointUrl(versions []Version) (string, error) {
	var endpointUrl string

	for _, version := range versions {
		if version.Version == "2.2" {
			endpointUrl = version.Url
			break
		}
	}

	if endpointUrl == "" {
		return "", fmt.Errorf("no version 2.2 endpoint found")
	}

	return endpointUrl, nil
}

func (o *OCPI) getEndpoints(ctx context.Context, url, token string) ([]Endpoint, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var versionDetail OcpiResponseVersionDetail
	err = json.Unmarshal(b, &versionDetail)
	if err != nil {
		return nil, err
	}
	if versionDetail.StatusCode != StatusSuccess {
		return nil, fmt.Errorf("status code: %d", versionDetail.StatusCode)
	}

	return versionDetail.Data.Endpoints, nil
}

func generateRandomString() (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, 64)
	for i := 0; i < 64; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}

func getCredentialsUrl(endpoints []Endpoint) (string, error) {
	for _, endpoint := range endpoints {
		if endpoint.Identifier == "credentials" && endpoint.Role == RECEIVER {
			return endpoint.Url, nil
		}
	}
	return "", errors.New("no credentials endpoint for receiver found")
}

func (o *OCPI) postCredentials(ctx context.Context, url, token, newToken string) error {
	creds := Credentials{
		Roles: []CredentialsRole{
			{
				CountryCode: o.countryCode,
				PartyId:     o.partyId,
				Role:        "CPO",
			},
		},
		Token: newToken,
		Url:   o.externalUrl + "/ocpi/versions",
	}

	b, err := json.Marshal(creds)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	return nil
}
