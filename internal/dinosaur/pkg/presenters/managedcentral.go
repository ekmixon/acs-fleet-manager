package presenters

import (
	"time"

	private "github.com/stackrox/acs-fleet-manager/generated/privateapi"
	"github.com/stackrox/acs-fleet-manager/internal/dinosaur/pkg/api/dbapi"
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
		Metadata: &private.ManagedCentralMetadata{
			Name:      from.Name,
			Namespace: from.Namespace,
			Annotations: &private.ManagedCentralMetadataAnnotations{
				MasId:          from.ID,
				MasPlacementId: from.PlacementID,
			},
		},
		Spec: &private.ManagedCentralSpec{
			Owners: []string{
				from.Owner,
			},
			Auth: &private.ManagedCentralSpecAuth{
				ClientSecret: c.centralConfig.RhSsoClientSecret, // pragma: allowlist secret
				// TODO(ROX-11593): make part of centralConfig
				ClientId:    "rhacs-ms-dev",
				OwnerOrgId:  from.OrganisationID,
				OwnerUserId: from.OwnerUserID,
				Issuer:      c.centralConfig.RhSsoIssuer,
			},
			UiEndpoint: &private.ManagedCentralSpecUiEndpoint{
				Host: from.GetUIHost(),
				Tls: &private.ManagedCentralSpecUiEndpointTLS{
					Cert: c.centralConfig.CentralTLSCert,
					Key:  c.centralConfig.CentralTLSKey,
				},
			},
			DataEndpoint: &private.ManagedCentralSpecDataEndpoint{
				Host: from.GetDataHost(),
			},
			Versions: &private.CentralVersions{
				Central:         from.DesiredCentralVersion,
				CentralOperator: from.DesiredCentralOperatorVersion,
			},
			// TODO(create-ticket): add additional CAs to public create/get centrals api and internal models
			Central: &private.ManagedCentralSpecCentral{
				Resources: &private.ResourceRequirements{
					Requests: &private.ResourceList{
						Cpu:    defaultCentralRequestCPU,
						Memory: defaultCentralRequestMemory,
					},
					Limits: &private.ResourceList{
						Cpu:    defaultCentralLimitCPU,
						Memory: defaultCentralLimitMemory,
					},
				},
			},
			Scanner: &private.ManagedCentralSpecScanner{
				Analyzer: &private.ManagedCentralSpecScannerAnalyzer{
					Scaling: &private.ManagedCentralSpecScannerAnalyzerScaling{
						AutoScaling: defaultScannerAnalyzerAutoScaling,
						Replicas:    defaultScannerAnalyzerScalingReplicas,
						MinReplicas: defaultScannerAnalyzerScalingMinReplicas,
						MaxReplicas: defaultScannerAnalyzerScalingMaxReplicas,
					},
					Resources: &private.ResourceRequirements{
						Requests: &private.ResourceList{
							Cpu:    defaultScannerAnalyzerRequestCPU,
							Memory: defaultScannerAnalyzerRequestMemory,
						},
						Limits: &private.ResourceList{
							Cpu:    defaultScannerAnalyzerLimitCPU,
							Memory: defaultScannerAnalyzerLimitMemory,
						},
					},
				},
				Db: &private.ManagedCentralSpecScannerDb{
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
