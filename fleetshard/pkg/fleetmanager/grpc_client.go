package fleetmanager

import (
	private "github.com/stackrox/acs-fleet-manager/generated/privateapi"
)

// GRPCClient represents the gRPC client for connecting to fleet-manager
type GRPCClient struct {
	client private.FleetManagerPrivateClient
}
