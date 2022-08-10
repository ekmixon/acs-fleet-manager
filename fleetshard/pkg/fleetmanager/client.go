package fleetmanager

import "github.com/stackrox/acs-fleet-manager/internal/dinosaur/pkg/api/private"

// Client interface for communicating with the fleet-manager
type Client interface {
	// GetManagedCentralList returns a list of centrals from fleet-manager which should be managed by this fleetshard.
	GetManagedCentralList() (*private.ManagedCentralList, error)
	// UpdateStatus batch updates the status of managed centrals. The status param takes a map of DataPlaneCentralStatus indexed by
	// the Centrals ID.
	UpdateStatus(statuses map[string]private.DataPlaneCentralStatus) error
}
