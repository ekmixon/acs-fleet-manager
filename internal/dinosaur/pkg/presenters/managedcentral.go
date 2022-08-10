package presenters

import (
	"time"

	"github.com/stackrox/acs-fleet-manager/internal/dinosaur/pkg/api/dbapi"
	"github.com/stackrox/acs-fleet-manager/internal/dinosaur/pkg/api/private"
	fmgrpc "github.com/stackrox/acs-fleet-manager/internal/dinosaur/pkg/api/private/grpc"
	"github.com/stackrox/acs-fleet-manager/internal/dinosaur/pkg/config"
)

// TODO(create-ticket): implement configurable central and scanner resources
const (
	defaultCentralRequestMemory = "250Mi"
	defaultCentralRequestCPU    = "250m"
	defaultCentralLimitMemory   = "4Gi"
	defaultCentralLimitCPU      = "1000m"

	defaultScannerAnalyzerRequestMemory = "100Mi"
	defaultScannerAnalyzerRequestCPU    = "250m"
	defaultScannerAnalyzerLimitMemory   = "2500Mi"
	defaultScannerAnalyzerLimitCPU      = "2000m"

	defaultScannerAnalyzerAutoScaling        = "enabled"
	defaultScannerAnalyzerScalingReplicas    = 1
	defaultScannerAnalyzerScalingMinReplicas = 1
	defaultScannerAnalyzerScalingMaxReplicas = 3
)

// ManagedCentralPresenter helper service which converts Central DB representation to the private API representation
type ManagedCentralPresenter struct {
	centralConfig *config.CentralConfig
}

// NewManagedCentralPresenter creates a new instance of ManagedCentralPresenter
func NewManagedCentralPresenter(config *config.CentralConfig) *ManagedCentralPresenter {
	return &ManagedCentralPresenter{centralConfig: config}
}

// PresentManagedCentral converts DB representation of Central to the private API representation
func (c *ManagedCentralPresenter) PresentManagedCentral(from *dbapi.CentralRequest) private.ManagedCentral {
	res := private.ManagedCentral{
		Id:   from.ID,
		Kind: "ManagedCentral",
		Metadata: private.ManagedCentralAllOfMetadata{
			Name:      from.Name,
			Namespace: from.Namespace,
			Annotations: private.ManagedCentralAllOfMetadataAnnotations{
				MasId:          from.ID,
				MasPlacementId: from.PlacementID,
			},
		},
		Spec: private.ManagedCentralAllOfSpec{
			Owners: []string{
				from.Owner,
			},
			Auth: private.ManagedCentralAllOfSpecAuth{
				ClientSecret: c.centralConfig.RhSsoClientSecret, // pragma: allowlist secret
				// TODO(ROX-11593): make part of centralConfig
				ClientId:    "rhacs-ms-dev",
				OwnerOrgId:  from.OrganisationID,
				OwnerUserId: from.OwnerUserID,
				Issuer:      c.centralConfig.RhSsoIssuer,
			},
			UiEndpoint: private.ManagedCentralAllOfSpecUiEndpoint{
				Host: from.GetUIHost(),
				Tls: private.ManagedCentralAllOfSpecUiEndpointTls{
					Cert: c.centralConfig.CentralTLSCert,
					Key:  c.centralConfig.CentralTLSKey,
				},
			},
			DataEndpoint: private.ManagedCentralAllOfSpecDataEndpoint{
				Host: from.GetDataHost(),
			},
			Versions: private.ManagedCentralVersions{
				Central:         from.DesiredCentralVersion,
				CentralOperator: from.DesiredCentralOperatorVersion,
			},
			// TODO(create-ticket): add additional CAs to public create/get centrals api and internal models
			Central: private.ManagedCentralAllOfSpecCentral{
				Resources: private.ResourceRequirements{
					Requests: private.ResourceList{
						Cpu:    defaultCentralRequestCPU,
						Memory: defaultCentralRequestMemory,
					},
					Limits: private.ResourceList{
						Cpu:    defaultCentralLimitCPU,
						Memory: defaultCentralLimitMemory,
					},
				},
			},
			Scanner: private.ManagedCentralAllOfSpecScanner{
				Analyzer: private.ManagedCentralAllOfSpecScannerAnalyzer{
					Scaling: private.ManagedCentralAllOfSpecScannerAnalyzerScaling{
						AutoScaling: defaultScannerAnalyzerAutoScaling,
						Replicas:    defaultScannerAnalyzerScalingReplicas,
						MinReplicas: defaultScannerAnalyzerScalingMinReplicas,
						MaxReplicas: defaultScannerAnalyzerScalingMaxReplicas,
					},
					Resources: private.ResourceRequirements{
						Requests: private.ResourceList{
							Cpu:    defaultScannerAnalyzerRequestCPU,
							Memory: defaultScannerAnalyzerRequestMemory,
						},
						Limits: private.ResourceList{
							Cpu:    defaultScannerAnalyzerLimitCPU,
							Memory: defaultScannerAnalyzerLimitMemory,
						},
					},
				},
				Db: private.ManagedCentralAllOfSpecScannerDb{
					// TODO:(create-ticket): add DB configuration values to ManagedCentral Scanner
					Host: "dbhost.rhacs-psql-instance",
				},
			},
		},
		RequestStatus: from.Status,
	}

	if from.DeletionTimestamp != nil {
		res.Metadata.DeletionTimestamp = from.DeletionTimestamp.Format(time.RFC3339)
	}

	return res
}

// PresentManagedCentralGRPC converts DB representation of Central to the gRPC API representation
func (c *ManagedCentralPresenter) PresentManagedCentralGRPC(from *dbapi.CentralRequest) *fmgrpc.ManagedCentral {
	res := &fmgrpc.ManagedCentral{
		Id:   from.ID,
		Kind: "ManagedCentral",
		Metadata: &fmgrpc.ManagedCentralMetadata{
			Name:      from.Name,
			Namespace: from.Namespace,
			Annotations: &fmgrpc.ManagedCentralMetadataAnnotations{
				MasId:          from.ID,
				MasPlacementId: from.PlacementID,
			},
		},
		Spec: &fmgrpc.ManagedCentralSpec{
			Owners: []string{
				from.Owner,
			},
			Auth: &fmgrpc.ManagedCentralSpecAuth{
				ClientSecret: c.centralConfig.RhSsoClientSecret, // pragma: allowlist secret
				// TODO(ROX-11593): make part of centralConfig
				ClientId:    "rhacs-ms-dev",
				OwnerOrgId:  from.OrganisationID,
				OwnerUserId: from.OwnerUserID,
				Issuer:      c.centralConfig.RhSsoIssuer,
			},
			UiEndpoint: &fmgrpc.ManagedCentralSpecUiEndpoint{
				Host: from.GetUIHost(),
				Tls: &fmgrpc.ManagedCentralSpecUiEndpointTLS{
					Cert: c.centralConfig.CentralTLSCert,
					Key:  c.centralConfig.CentralTLSKey,
				},
			},
			DataEndpoint: &fmgrpc.ManagedCentralSpecDataEndpoint{
				Host: from.GetDataHost(),
			},
			Versions: &fmgrpc.CentralVersions{
				Central:         from.DesiredCentralVersion,
				CentralOperator: from.DesiredCentralOperatorVersion,
			},
			// TODO(create-ticket): add additional CAs to public create/get centrals api and internal models
			Central: &fmgrpc.ManagedCentralSpecCentral{
				Resources: &fmgrpc.ResourceRequirements{
					Requests: &fmgrpc.ResourceList{
						Cpu:    defaultCentralRequestCPU,
						Memory: defaultCentralRequestMemory,
					},
					Limits: &fmgrpc.ResourceList{
						Cpu:    defaultCentralLimitCPU,
						Memory: defaultCentralLimitMemory,
					},
				},
			},
			Scanner: &fmgrpc.ManagedCentralSpecScanner{
				Analyzer: &fmgrpc.ManagedCentralSpecScannerAnalyzer{
					Scaling: &fmgrpc.ManagedCentralSpecScannerAnalyzerScaling{
						AutoScaling: defaultScannerAnalyzerAutoScaling,
						Replicas:    defaultScannerAnalyzerScalingReplicas,
						MinReplicas: defaultScannerAnalyzerScalingMinReplicas,
						MaxReplicas: defaultScannerAnalyzerScalingMaxReplicas,
					},
					Resources: &fmgrpc.ResourceRequirements{
						Requests: &fmgrpc.ResourceList{
							Cpu:    defaultScannerAnalyzerRequestCPU,
							Memory: defaultScannerAnalyzerRequestMemory,
						},
						Limits: &fmgrpc.ResourceList{
							Cpu:    defaultScannerAnalyzerLimitCPU,
							Memory: defaultScannerAnalyzerLimitMemory,
						},
					},
				},
				Db: &fmgrpc.ManagedCentralSpecScannerDb{
					// TODO:(create-ticket): add DB configuration values to ManagedCentral Scanner
					Host: "dbhost.rhacs-psql-instance",
				},
			},
		},
		RequestStatus: from.Status,
	}

	if from.DeletionTimestamp != nil {
		res.Metadata.DeletionTimestamp = from.DeletionTimestamp.Format(time.RFC3339)
	}

	return res
}
