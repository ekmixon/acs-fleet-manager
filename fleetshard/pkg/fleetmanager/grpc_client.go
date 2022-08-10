package fleetmanager

import (
	privateGrpc "github.com/stackrox/acs-fleet-manager/internal/dinosaur/pkg/api/private/grpc"
)

// GRPCClient represents the gRPC client for connecting to fleet-manager
type GRPCClient struct {
	client privateGrpc.FleetManagerPrivateClient
}
