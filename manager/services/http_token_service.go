// SPDX-License-Identifier: Apache-2.0

package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"k8s.io/utils/clock"
	"net/http"
	"os"
	"regexp"
	"sync"
	"time"
)

type HttpTokenService interface {
	GetToken(ctx context.Context, refresh bool) (string, error)
}

type CachingHttpTokenService struct {
	sync.Mutex
	tokenService HttpTokenService
	ttl          time.Duration
	clock        clock.PassiveClock
	expiry       time.Time
	token        string
}

func (h *CachingHttpTokenService) GetToken(ctx context.Context, refresh bool) (string, error) {
	h.Lock()
	defer h.Unlock()

	if refresh || h.clock.Now().After(h.expiry) {
		var err error
		h.token, err = h.tokenService.GetToken(ctx, true)
		if err != nil {
			return "", err
		}
		h.expiry = h.clock.Now().Add(h.ttl)
	}

	return h.token, nil
}

func NewCachingHttpTokenService(tokenService HttpTokenService, ttl time.Duration, clock clock.PassiveClock) *CachingHttpTokenService {
	return &CachingHttpTokenService{
		tokenService: tokenService,
		ttl:          ttl,
		clock:        clock,
		expiry:       clock.Now(),
	}
}

func NewEnvHttpTokenService(envVar string) (*FixedHttpTokenService, error) {
	value, ok := os.LookupEnv(envVar)
	if !ok {
		return nil, fmt.Errorf("environment variable %s not set", envVar)
	}
	return NewFixedHttpTokenService(value), nil
}

type FixedHttpTokenService struct {
	token string
}

func NewFixedHttpTokenService(token string) *FixedHttpTokenService {
	return &FixedHttpTokenService{
		token: token,
	}
}

func (f *FixedHttpTokenService) GetToken(_ context.Context, _ bool) (string, error) {
	return f.token, nil
}

type HubjectTestHttpTokenService struct {
	url    string
	client *http.Client
}

func NewHubjectTestHttpTokenService(url string, httpClient *http.Client) *HubjectTestHttpTokenService {
	return &HubjectTestHttpTokenService{
		url:    url,
		client: httpClient,
	}
}

type HubjectTestTokenResponse struct {
	Title string `json:"title"`
	Data  string `json:"data"`
}

var HubjectTestTokenRegexp = regexp.MustCompile(`Bearer (.+)\n`)

func (h *HubjectTestHttpTokenService) GetToken(ctx context.Context, _ bool) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, h.url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response body failed: %w", err)
	}

	var tokenResponse HubjectTestTokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return "", fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	matches := HubjectTestTokenRegexp.FindStringSubmatch(tokenResponse.Data)
	if len(matches) != 2 {
		return "", fmt.Errorf("failed to extract token from response body: %w", err)
	}

	return matches[1], nil
}
