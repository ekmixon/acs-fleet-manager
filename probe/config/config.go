package config

import (
	"time"

	"github.com/stackrox/rox/pkg/errorhelpers"

	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
)

// Config contains this application's runtime configuration.
type Config struct {
	DataCloudProvider    string        `env:"DATA_PLANE_CLOUD_PROVIDER" envDefault:"aws"`
	DataPlaneRegion      string        `env:"DATA_PLANE_REGION" envDefault:"us-east-1"`
	FleetManagerEndpoint string        `env:"FLEET_MANAGER_ENDPOINT" envDefault:"http://127.0.0.1:8000"`
	MetricsAddress       string        `env:"FLEETSHARD_METRICS_ADDRESS" envDefault:":8080"`
	RuntimePollPeriod    time.Duration `env:"RUNTIME_POLL_PERIOD" envDefault:"5s"`
}

// GetConfig retrieves the current runtime configuration from the environment and returns it.
func GetConfig() (*Config, error) {
	c := Config{}
	var configErrors errorhelpers.ErrorList

	if err := env.Parse(&c); err != nil {
		return nil, errors.Wrapf(err, "Unable to parse runtime configuration from environment.")
	}
	cfgErr := configErrors.ToError()
	if cfgErr != nil {
		return nil, errors.Wrap(cfgErr, "Unexpected configuration settings.")
	}
	return &c, nil
}
