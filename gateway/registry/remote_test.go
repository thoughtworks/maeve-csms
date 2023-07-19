package registry_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thoughtworks/maeve-csms/gateway/registry"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLookupChargeStation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"securityProfile":1,"base64SHA256Password":"DEADBEEF"}`))
	}))
	defer server.Close()

	reg := registry.RemoteRegistry{
		ManagerApiAddr: server.URL,
	}

	want := &registry.ChargeStation{
		ClientId:             "cs001",
		SecurityProfile:      1,
		Base64SHA256Password: "DEADBEEF",
	}

	got, _ := reg.LookupChargeStation("cs001")
	require.NotNil(t, got)

	assert.Equal(t, want, got)
}
