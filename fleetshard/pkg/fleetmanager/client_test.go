package fleetmanager

import (
	"encoding/json"
	"github.com/stackrox/acs-fleet-manager/internal/dinosaur/pkg/api/private"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestClientGetManagedCentralList(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		assert.Contains(t, request.RequestURI, "/api/dinosaurs_mgmt/v1/agent-clusters/cluster-id/dinosaurs")
		bytes, err := json.Marshal(private.ManagedDinosaurList{})
		require.NoError(t, err)
		_, err = writer.Write(bytes)
		require.NoError(t, err)
	}))
	defer ts.Close()

	err := os.Setenv("OCM_TOKEN", "token")
	require.NoError(t, err)

	client, err := NewClient(ts.URL, "cluster-id")
	require.NoError(t, err)

	result, err := client.GetManagedCentralList()
	require.NoError(t, err)
	assert.Equal(t, &private.ManagedDinosaurList{}, result)
}

func TestClientUpdateStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		assert.Contains(t, request.RequestURI, "/api/dinosaurs_mgmt/v1/agent-clusters/cluster-id/dinosaurs")
		bytes, err := json.Marshal(private.ManagedDinosaurList{})
		require.NoError(t, err)
		_, err = writer.Write(bytes)
		require.NoError(t, err)
	}))
	defer ts.Close()
}