package fleetmanager

import (
	private "github.com/stackrox/acs-fleet-manager/generated/privateapi"
	"golang.org/x/net/context"
)

// Client interface for communicating with the fleet-manager
type Client interface {
	// GetManagedCentralList returns a list of centrals from fleet-manager which should be managed by this fleetshard.
	GetManagedCentralList(ctx context.Context) ([]*private.ManagedCentral, error)
	// UpdateStatus updates the status of managed central with given ID
	UpdateStatus(ctx context.Context, id string, status *private.CentralStatus) error
	// Close gracefully shutdowns client resources
	Close() error
}
