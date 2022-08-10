package fleetmanager

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	private "github.com/stackrox/acs-fleet-manager/generated/privateapi"
	"github.com/stackrox/acs-fleet-manager/internal/dinosaur/compat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type noAuth struct{}

func (n noAuth) AddAuth(_ *http.Request) error {
	return nil
}

func TestClientGetManagedCentralList(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		assert.Contains(t, request.RequestURI, "/api/rhacs/v1/agent-clusters/cluster-id/centrals")
		bytes, err := json.Marshal([]private.ManagedCentral{})
		require.NoError(t, err)
		_, err = writer.Write(bytes)
		require.NoError(t, err)
	}))
	defer ts.Close()

	client, err := NewHTTPClient(ts.URL, "cluster-id", &noAuth{})
	require.NoError(t, err)

	result, err := client.GetManagedCentralList(context.TODO())
	require.NoError(t, err)
	assert.Equal(t, []private.ManagedCentral{}, result)
}

func TestClientReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		assert.Contains(t, request.RequestURI, "/api/rhacs/v1/agent-clusters/cluster-id/centrals")
		bytes, err := json.Marshal(compat.Error{
			Kind:   "error",
			Reason: "some reason",
		})
		require.NoError(t, err)
		_, err = writer.Write(bytes)
		require.NoError(t, err)
	}))
	defer ts.Close()

	client, err := NewHTTPClient(ts.URL, "cluster-id", &noAuth{})
	require.NoError(t, err)

	_, err = client.GetManagedCentralList(context.TODO())
	require.Error(t, err)
	assert.ErrorContains(t, err, "some reason")
}

func TestClientUpdateStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		assert.Contains(t, request.RequestURI, "/api/rhacs/v1/agent-clusters/cluster-id/centrals")
	}))
	defer ts.Close()

	client, err := NewHTTPClient(ts.URL, "cluster-id", &noAuth{})
	require.NoError(t, err)

	err = client.UpdateStatus(context.TODO(), "123", &private.CentralStatus{})
	require.NoError(t, err)
}
