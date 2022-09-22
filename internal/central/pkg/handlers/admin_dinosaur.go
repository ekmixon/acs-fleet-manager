package handlers

import (
	"fmt"
	"net/http"

	"github.com/stackrox/acs-fleet-manager/pkg/services/account"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/gorilla/mux"
	"github.com/stackrox/acs-fleet-manager/internal/central/pkg/api/admin/private"
	"github.com/stackrox/acs-fleet-manager/internal/central/pkg/api/dbapi"
	"github.com/stackrox/acs-fleet-manager/internal/central/pkg/api/public"
	"github.com/stackrox/acs-fleet-manager/internal/central/pkg/config"
	"github.com/stackrox/acs-fleet-manager/internal/central/pkg/converters"
	"github.com/stackrox/acs-fleet-manager/internal/central/pkg/defaults"
	"github.com/stackrox/acs-fleet-manager/internal/central/pkg/presenters"
	"github.com/stackrox/acs-fleet-manager/internal/central/pkg/services"
	"github.com/stackrox/acs-fleet-manager/pkg/errors"
	"github.com/stackrox/acs-fleet-manager/pkg/handlers"
	coreServices "github.com/stackrox/acs-fleet-manager/pkg/services"
)

type adminCentralHandler struct {
	service        services.CentralService
	accountService account.AccountService
	providerConfig *config.ProviderConfig
}

// NewAdminCentralHandler ...
func NewAdminCentralHandler(service services.CentralService, accountService account.AccountService, providerConfig *config.ProviderConfig) *adminCentralHandler {
	return &adminCentralHandler{
		service:        service,
		accountService: accountService,
		providerConfig: providerConfig,
	}
}

// Create ...
func (h adminCentralHandler) Create(w http.ResponseWriter, r *http.Request) {
	centralRequest := public.CentralRequestPayload{
		Central: public.CentralSpec{
			Resources: converters.ConvertCoreV1ResourceRequirementsToPublic(&defaults.CentralResources),
		},
		Scanner: public.ScannerSpec{
			Analyzer: public.ScannerSpecAnalyzer{
				Resources: converters.ConvertCoreV1ResourceRequirementsToPublic(&defaults.ScannerAnalyzerResources),
				Scaling:   converters.ConvertScalingToPublic(&dbapi.DefaultScannerAnalyzerScaling),
			},
			Db: public.ScannerSpecDb{
				Resources: converters.ConvertCoreV1ResourceRequirementsToPublic(&defaults.ScannerDbResources),
			},
		},
	}
	ctx := r.Context()
	convCentral := dbapi.CentralRequest{}

	cfg := &handlers.HandlerConfig{
		MarshalInto: &centralRequest,
		Validate: []handlers.Validate{
			handlers.ValidateAsyncEnabled(r, "creating central requests"),
			handlers.ValidateLength(&centralRequest.Name, "name", &handlers.MinRequiredFieldLength, &MaxCentralNameLength),
			ValidCentralClusterName(&centralRequest.Name, "name"),
			ValidateCentralClusterNameIsUnique(r.Context(), &centralRequest.Name, h.service),
			ValidateCentralClaims(ctx, &centralRequest, &convCentral),
			ValidateCloudProvider(&h.service, &convCentral, h.providerConfig, "creating central requests"),
			handlers.ValidateMultiAZEnabled(&centralRequest.MultiAz, "creating central requests"),
			ValidateCentralSpec(ctx, &centralRequest, &convCentral),
			ValidateScannerSpec(ctx, &centralRequest, &convCentral),
		},
		Action: func() (interface{}, *errors.ServiceError) {
			svcErr := h.service.RegisterCentralJob(&convCentral)
			if svcErr != nil {
				return nil, svcErr
			}
			// TODO(mclasmeier): Do we need PresentDinosaurRequestAdminEndpoint?
			return presenters.PresentCentralRequest(&convCentral), nil
		},
	}

	// return 202 status accepted
	handlers.Handle(w, r, cfg, http.StatusAccepted)
}

// Get ...
func (h adminCentralHandler) Get(w http.ResponseWriter, r *http.Request) {
	cfg := &handlers.HandlerConfig{
		Action: func() (i interface{}, serviceError *errors.ServiceError) {
			id := mux.Vars(r)["id"]
			ctx := r.Context()
			centralRequest, err := h.service.Get(ctx, id)
			if err != nil {
				return nil, err
			}
			return presenters.PresentDinosaurRequestAdminEndpoint(centralRequest, h.accountService)
		},
	}
	handlers.HandleGet(w, r, cfg)
}

// List ...
func (h adminCentralHandler) List(w http.ResponseWriter, r *http.Request) {
	cfg := &handlers.HandlerConfig{
		Action: func() (interface{}, *errors.ServiceError) {
			ctx := r.Context()

			listArgs := coreServices.NewListArguments(r.URL.Query())

			if err := listArgs.Validate(); err != nil {
				return nil, errors.NewWithCause(errors.ErrorMalformedRequest, err, "Unable to list central requests: %s", err.Error())
			}

			centralRequests, paging, err := h.service.List(ctx, listArgs)
			if err != nil {
				return nil, err
			}

			centralRequestList := private.CentralList{
				Kind:  "CentralList",
				Page:  int32(paging.Page),
				Size:  int32(paging.Size),
				Total: int32(paging.Total),
				Items: []private.Central{},
			}

			for _, centralRequest := range centralRequests {
				converted, err := presenters.PresentDinosaurRequestAdminEndpoint(centralRequest, h.accountService)
				if err != nil {
					return nil, err
				}

				if converted != nil {
					centralRequestList.Items = append(centralRequestList.Items, *converted)
				}
			}

			return centralRequestList, nil
		},
	}

	handlers.HandleList(w, r, cfg)
}

// Delete ...
func (h adminCentralHandler) Delete(w http.ResponseWriter, r *http.Request) {
	cfg := &handlers.HandlerConfig{
		Validate: []handlers.Validate{
			handlers.ValidateAsyncEnabled(r, "deleting central requests"),
		},
		Action: func() (i interface{}, serviceError *errors.ServiceError) {
			id := mux.Vars(r)["id"]
			ctx := r.Context()

			err := h.service.RegisterCentralDeprovisionJob(ctx, id)
			return nil, err
		},
	}

	handlers.HandleDelete(w, r, cfg, http.StatusAccepted)
}

func updateResourcesList(to *corev1.ResourceList, from map[string]string) error {
	newResourceList := to.DeepCopy()
	for name, qty := range from {
		if qty == "" {
			continue
		}
		resourceName, isSupported := ValidateResourceName(name)
		if !isSupported {
			return fmt.Errorf("resource type %q is not supported", name)
		}
		resourceQty, err := resource.ParseQuantity(qty)
		if err != nil {
			return fmt.Errorf("parsing %s quantity %q: %w", resourceName, qty, err)
		}
		if newResourceList == nil {
			newResourceList = corev1.ResourceList(make(map[corev1.ResourceName]resource.Quantity))
		}
		newResourceList[resourceName] = resourceQty
	}
	*to = newResourceList
	return nil
}

func updateCoreV1Resources(to *corev1.ResourceRequirements, from private.ResourceRequirements) error {
	newResources := to.DeepCopy()

	err := updateResourcesList(&newResources.Limits, from.Limits)
	if err != nil {
		return err
	}
	err = updateResourcesList(&newResources.Requests, from.Requests)
	if err != nil {
		return err
	}

	*to = *newResources
	return nil
}

// updateCentralFromPrivateAPI updates the CentralSpec using the non-zero fields from the API's CentralSpec.
func updateCentralFromPrivateAPI(c *dbapi.CentralSpec, apiCentralSpec *private.CentralSpec) error {
	err := updateCoreV1Resources(&c.Resources, apiCentralSpec.Resources)
	if err != nil {
		return fmt.Errorf("updating resources within CentralSpec: %w", err)
	}
	return nil
}

// updateScannerFromPrivateAPI updates the ScannerSpec using the non-zero fields from the API's ScannerSpec.
func updateScannerFromPrivateAPI(s *dbapi.ScannerSpec, apiSpec *private.ScannerSpec) error {
	var err error
	new := *s

	err = updateCoreV1Resources(&new.Analyzer.Resources, apiSpec.Analyzer.Resources)
	if err != nil {
		return fmt.Errorf("updating resources within ScannerSpec Analyzer: %w", err)
	}
	err = updateScannerAnalyzerScaling(&new.Analyzer.Scaling, apiSpec.Analyzer.Scaling)
	if err != nil {
		return fmt.Errorf("updating scaling configuration within ScannerSpec Analyzer: %w", err)
	}
	err = updateCoreV1Resources(&new.Db.Resources, apiSpec.Db.Resources)
	if err != nil {
		return fmt.Errorf("updating resources within ScannerSpec DB: %w", err)
	}
	*s = new
	return nil
}

func updateScannerAnalyzerScaling(s *dbapi.ScannerAnalyzerScaling, apiScaling private.ScannerSpecAnalyzerScaling) error {
	if apiScaling.AutoScaling != "" {
		s.AutoScaling = apiScaling.AutoScaling
	}
	if apiScaling.MaxReplicas > 0 {
		s.MaxReplicas = apiScaling.MaxReplicas
	}
	if apiScaling.MinReplicas > 0 {
		s.MinReplicas = apiScaling.MinReplicas
	}
	if apiScaling.Replicas > 0 {
		s.Replicas = apiScaling.Replicas
	}
	return nil
}

func updateCentralRequest(request *dbapi.CentralRequest, updateRequest *private.CentralUpdateRequest) error {
	if updateRequest == nil {
		return nil
	}

	centralSpec, err := request.GetCentralSpec()
	if err != nil {
		return fmt.Errorf("retrieving CentralSpec from CentralRequest: %w", err)
	}
	scannerSpec, err := request.GetScannerSpec()
	if err != nil {
		return fmt.Errorf("retrieving ScannerSpec from CentralRequest: %w", err)
	}

	err = updateCentralFromPrivateAPI(centralSpec, &updateRequest.Central)
	if err != nil {
		return fmt.Errorf("updating CentralSpec from CentralUpdateRequest: %w", err)
	}
	err = updateScannerFromPrivateAPI(scannerSpec, &updateRequest.Scanner)
	if err != nil {
		return fmt.Errorf("updating ScannerSpec from CentralUpdateRequest: %w", err)
	}

	err = ValidateScannerAnalyzerScaling(&scannerSpec.Analyzer.Scaling)
	if err != nil {
		return err
	}

	// TODO(create-ticket): We should also validate the resource configuration here. If the configuration is invalid
	// the operator will not be able to create the Central instance and we could fail early here.

	new := *request

	err = new.SetCentralSpec(centralSpec)
	if err != nil {
		return fmt.Errorf("updating CentralSpec within CentralRequest: %w", err)
	}

	err = new.SetScannerSpec(scannerSpec)
	if err != nil {
		return fmt.Errorf("updating ScannerSpec within CentralRequest: %w", err)
	}

	*request = new
	return nil
}

// Update a Central instance.
func (h adminCentralHandler) Update(w http.ResponseWriter, r *http.Request) {

	var centralUpdateRequest private.CentralUpdateRequest
	cfg := &handlers.HandlerConfig{
		MarshalInto: &centralUpdateRequest,
		Validate:    []handlers.Validate{},
		Action: func() (i interface{}, serviceError *errors.ServiceError) {
			id := mux.Vars(r)["id"]
			ctx := r.Context()
			centralRequest, svcErr := h.service.Get(ctx, id)
			if svcErr != nil {
				return nil, svcErr
			}

			err := updateCentralRequest(centralRequest, &centralUpdateRequest)
			if err != nil {
				return nil, errors.NewWithCause(errors.ErrorBadRequest, err, "Updating CentralRequest: %s", err.Error())
			}

			svcErr = h.service.VerifyAndUpdateCentralAdmin(ctx, centralRequest)
			if svcErr != nil {
				return nil, svcErr
			}
			return presenters.PresentDinosaurRequestAdminEndpoint(centralRequest, h.accountService)
		},
	}
	handlers.Handle(w, r, cfg, http.StatusOK)
}
