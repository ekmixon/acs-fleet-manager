package server

import (
	"context"
	"fmt"
	"net"

	"github.com/golang/glog"
	private "github.com/stackrox/acs-fleet-manager/generated/privateapi"
	"github.com/stackrox/acs-fleet-manager/internal/dinosaur/pkg/api/dbapi"
	"github.com/stackrox/acs-fleet-manager/internal/dinosaur/pkg/presenters"
	"github.com/stackrox/acs-fleet-manager/pkg/environments"
	"github.com/stackrox/acs-fleet-manager/pkg/errors"
	"google.golang.org/grpc"
)

// GRPCServer exposes the gRPC API of fleet-manager
type GRPCServer struct {
	private.UnimplementedFleetManagerPrivateServer
	listener       net.Listener
	server         *grpc.Server
	presenter      *presenters.ManagedCentralPresenter
	centralService CentralService
}

var _ Server = &GRPCServer{}
var _ environments.BootService = &GRPCServer{}

// CentralService is the interface to DB operations, it was required to wrap the central service
// defined in the services package since otherwise a cyclyc dependency would be introduced.
// The reason for the cyclyc dependency is that some config values from the server package are used in the services
// package, but it was easier to break the dependency here
type CentralService interface {
	ListByClusterID(clusterID string) ([]*dbapi.CentralRequest, *errors.ServiceError)
}

// NewGRPCServer returns a new instance of GRPCServer
func NewGRPCServer(p *presenters.ManagedCentralPresenter, centralService CentralService) *GRPCServer {
	s := &GRPCServer{}
	s.presenter = p
	s.centralService = centralService
	return s
}

// UpdateCentralStatus handles requests to Update the status of centrals
func (gsv *GRPCServer) UpdateCentralStatus(context.Context, *private.UpdateCentralStatusRequest) (*private.UpdateCentralStatusResponse, error) {
	return nil, fmt.Errorf("not yet implemented")
}

// GetCentrals handles gRPC requests to get all central requests in the DB that are in correct state to be polled
// for provisioning
func (gsv *GRPCServer) GetCentrals(ctx context.Context, req *private.GetCentralsRequest) (*private.GetCentralsResponse, error) {
	dbCentrals, err := gsv.centralService.ListByClusterID(req.AgentCluster)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	centrals := make([]*private.ManagedCentral, 0, len(dbCentrals))
	for _, from := range dbCentrals {
		central := gsv.presenter.PresentManagedCentral(from)
		centrals = append(centrals, &central)
	}

	return &private.GetCentralsResponse{Centrals: centrals}, nil
}

// Run starts running the GRPCServer sequentially
func (gsv *GRPCServer) Run() {
	listener, err := gsv.Listen()
	if err != nil {
		glog.Fatalf("Unable to start gRPC server: %s", err)
	}

	gsv.listener = listener
	gsv.server = grpc.NewServer()
	private.RegisterFleetManagerPrivateServer(gsv.server, gsv)

	gsv.Serve(listener)
}

// Stop gracefully stops the GRPCServer
func (gsv *GRPCServer) Stop() {
	gsv.server.GracefulStop()
}

// Listen starts to listen to network connections
func (gsv *GRPCServer) Listen() (net.Listener, error) {
	l, err := net.Listen("tcp", ":8088")
	if err != nil {
		return l, fmt.Errorf("starting the listener: %w", err)
	}

	return l, nil
}

// Serve starts serving gRPC requests
func (gsv *GRPCServer) Serve(l net.Listener) {
	if err := gsv.server.Serve(l); err != nil {
		glog.Errorf("Unable to serve gRPC: %v", err)
	}
}

// Start runs the GRPCServer in a seperate go routine
func (gsv *GRPCServer) Start() {
	go gsv.Run()
}
