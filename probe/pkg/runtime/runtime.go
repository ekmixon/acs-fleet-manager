package runtime

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/stackrox/acs-fleet-manager/fleetshard/pkg/fleetmanager"
	"github.com/stackrox/acs-fleet-manager/internal/dinosaur/constants"
	"github.com/stackrox/acs-fleet-manager/internal/dinosaur/pkg/api/public"
	"github.com/stackrox/acs-fleet-manager/probe/config"
	"github.com/stackrox/acs-fleet-manager/probe/pkg/metrics"
)

// Runtime performs a probe run against fleet manager.
type Runtime struct {
	client   *fleetmanager.Client
	config   *config.Config
	request  *public.CentralRequestPayload
	response *public.CentralRequest
}

// NewRuntime creates a new runtime.
func NewRuntime(config *config.Config) (*Runtime, error) {
	// TODO: Authenticate with RH SSO service account. Refactor out the client from fleetshard sync.
	authType := "OCM"
	auth, err := fleetmanager.NewAuth(authType)
	if err != nil {
		return nil, fmt.Errorf("Auth creation failed: %w\n", err)
	}

	client, err := fleetmanager.NewClient(config.FleetManagerEndpoint, "cluster-id", auth)
	if err != nil {
		return nil, fmt.Errorf("Fleetmanager client creation failed: %w\n", err)
	}

	centralName, err := newCentralName()
	if err != nil {
		return nil, err
	}
	request := &public.CentralRequestPayload{
		Name:          centralName,
		MultiAz:       true,
		CloudProvider: config.DataCloudProvider,
		Region:        config.DataPlaneRegion,
	}

	return &Runtime{
		client:  client,
		config:  config,
		request: request,
	}, nil
}

// Run executes a probe run.
func (r *Runtime) Run() error {
	// TODO: Add total timeout
	// TODO: Measure total runtime and expose as metric
	metrics.MetricsInstance().IncStartedRuns()

	// Create Central
	response, err := r.client.CreateCentral(*r.request)
	if err != nil {
		metrics.MetricsInstance().IncFailedRuns()
		return fmt.Errorf("Central request failed: %w\n", err)
	}
	r.response = response
	glog.Infof("Central creation requested: %+v\n", r.response)

	// Poll ready state
	err = r.ensureCentralState(constants.CentralRequestStatusReady.String())
	if err != nil {
		metrics.MetricsInstance().IncFailedRuns()
		return err
	}

	// TODO: Ping Central UI. Verify login.

	// Deprovision
	err = r.client.DeleteCentral(r.response.Id)
	if err != nil {
		metrics.MetricsInstance().IncFailedRuns()
		return fmt.Errorf("Central deletion failed: %w", err)
	}
	err = r.ensureCentralState(constants.CentralRequestStatusDeprovision.String())
	if err != nil {
		metrics.MetricsInstance().IncFailedRuns()
		return err
	}

	// Deleting
	err = r.ensureCentralState(constants.CentralRequestStatusDeleting.String())
	if err != nil {
		metrics.MetricsInstance().IncFailedRuns()
		return fmt.Errorf("Central did not reach deprovision state: %w\n", err)
	}

	metrics.MetricsInstance().IncSuccessfulRuns()

	return nil
}

// Stop the probe run.
// TODO: Add write ahead log to clean up after ungraceful restarts.
func (r *Runtime) Stop() error {
	err := r.client.DeleteCentral(r.response.Id)
	if err != nil {
		return fmt.Errorf("Central deletion failed: %w", err)
	}
	glog.Infof("Central deletion requested.")
	return nil
}

func newCentralName() (string, error) {
	rnd := make([]byte, 8)
	_, err := rand.Read(rnd)
	if err != nil {
		return "", fmt.Errorf("Reading random bytes for unique central name: %w", err)
	}
	rndString := hex.EncodeToString(rnd)

	return fmt.Sprintf("probe-%s", rndString), nil
}

func (r *Runtime) ensureCentralState(targetState string) error {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), r.config.RuntimePollTimeout)
	defer cancel()

	isSuccess := make(chan bool, 1)

	go r.pollCentral(ctxTimeout, isSuccess, targetState)

	select {
	case <-ctxTimeout.Done():
		return fmt.Errorf("Central did not reach %s state: %w", targetState, ctxTimeout.Err())
	case <-isSuccess:
		return nil
	}
}

func (r *Runtime) pollCentral(ctx context.Context, isSuccess chan bool, targetState string) error {
	for {
		central, err := r.client.GetCentral(r.response.Id)
		// TODO: Only retry when error is recoverable
		if err != nil {
			continue
		}
		r.response = central

		if r.response.Status == targetState {
			glog.Infof("Central is in `%s` state.", targetState)
			isSuccess <- true
			return nil
		}
		time.Sleep(r.config.RuntimePollPeriod)
	}
}
