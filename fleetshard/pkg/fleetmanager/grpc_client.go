package fleetmanager

import (
	"context"

	"github.com/pkg/errors"
	private "github.com/stackrox/acs-fleet-manager/generated/privateapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var _ Client = (*GRPCClient)(nil)

// GRPCClient represents the gRPC client for connecting to fleet-manager
type GRPCClient struct {
	client    private.FleetManagerPrivateClient
	conn      *grpc.ClientConn
	clusterID string
}

// NewGRPCClient creates a new instance of GRPCClient
func NewGRPCClient(endpoint string, clusterID string) (*GRPCClient, error) {
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials())) // TODO(create-ticket): Add authentication
	if err != nil {
		return nil, errors.Wrap(err, "error initializing grpc connection")
	}
	client := private.NewFleetManagerPrivateClient(conn)
	return &GRPCClient{
		client:    client,
		conn:      conn,
		clusterID: clusterID,
	}, nil
}

// GetManagedCentralList returns a list of centrals from fleet-manager which should be managed by this fleetshard.
func (c *GRPCClient) GetManagedCentralList(ctx context.Context) ([]*private.ManagedCentral, error) {
	centrals, err := c.client.GetCentrals(ctx, &private.GetCentralsRequest{AgentCluster: c.clusterID})
	if err != nil {
		return nil, errors.Wrapf(err, "get centrals")
	}
	return centrals.Centrals, nil
}

// UpdateStatus updates the status of managed central with given ID
func (c *GRPCClient) UpdateStatus(ctx context.Context, id string, status *private.CentralStatus) error {
	statusUpdate := &private.UpdateCentralStatusRequest{
		AgentCluster: c.clusterID,
		Updates: map[string]*private.CentralStatus{
			id: status,
		},
	}
	_, err := c.client.UpdateCentralStatus(ctx, statusUpdate) // TODO(create-ticket): handle error response?
	return errors.Wrapf(err, "update central status")
}

// Close closes gRPC connection
func (c *GRPCClient) Close() error {
	err := c.conn.Close()
	return errors.Wrapf(err, "close grpc connection")
}
