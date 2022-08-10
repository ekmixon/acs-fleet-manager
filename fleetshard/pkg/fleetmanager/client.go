package fleetmanager

import "github.com/stackrox/acs-fleet-manager/internal/dinosaur/pkg/api/private"

// Client interface for communicating with the fleet-manager
type Client interface {
	// GetManagedCentralList returns a list of centrals from fleet-manager which should be managed by this fleetshard.
	GetManagedCentralList() ([]private.ManagedCentral, error)
	// UpdateStatus updates the status of managed central with given ID
	UpdateStatus(id string, status private.DataPlaneCentralStatus) error
}
