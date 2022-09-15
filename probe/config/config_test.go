package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSingleton_Success(t *testing.T) {
	t.Setenv("FLEET_MANAGER_ENDPOINT", "http://127.0.0.1:8888")
	t.Cleanup(func() {
		_ = os.Unsetenv("FLEET_MANAGER_ENDPOINT")
	})
	cfg, err := GetConfig()
	require.NoError(t, err)
	assert.Equal(t, cfg.FleetManagerEndpoint, "http://127.0.0.1:8888")
	assert.Equal(t, cfg.MetricsAddress, ":8080")
	assert.Equal(t, cfg.RuntimePollPeriod, 5*time.Second)
}

func TestSingleton_Failure(t *testing.T) {
	t.Cleanup(func() {
	})
	cfg, err := GetConfig()
	assert.Error(t, err)
	assert.Nil(t, cfg)
}
