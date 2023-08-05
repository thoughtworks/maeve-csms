package services

import (
	"fmt"
	"net/http"
	"os"
	"sync"
)

type HttpAuthService interface {
	Authenticate(r *http.Request) error
}

type EnvTokenHttpAuthService struct {
	sync.Mutex
	envVar   string
	envValue *string
}

func NewEnvTokenHttpAuthService(envVar string) *EnvTokenHttpAuthService {
	return &EnvTokenHttpAuthService{
		envVar: envVar,
	}
}

func (e *EnvTokenHttpAuthService) Authenticate(r *http.Request) error {
	e.Lock()
	defer e.Unlock()

	if e.envValue == nil {
		envValue := os.Getenv(e.envVar)
		e.envValue = &envValue
	}

	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *e.envValue))
	r.Header.Set("X-API-Key", *e.envValue)

	return nil
}

type FixedTokenHttpAuthService struct {
	token string
}

func NewFixedTokenHttpAuthService(token string) *FixedTokenHttpAuthService {
	return &FixedTokenHttpAuthService{
		token: token,
	}
}

func (f *FixedTokenHttpAuthService) Authenticate(r *http.Request) error {
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", f.token))
	r.Header.Set("X-API-Key", f.token)

	return nil
}
