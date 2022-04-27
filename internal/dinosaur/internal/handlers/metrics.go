package handlers

import (
	// "github.com/stackrox/acs-fleet-manager/internal/dinosaur/internal/api/public"
	"github.com/stackrox/acs-fleet-manager/internal/dinosaur/internal/metrics"
	// "github.com/stackrox/acs-fleet-manager/internal/dinosaur/internal/presenters"
	"github.com/stackrox/acs-fleet-manager/internal/dinosaur/internal/services"
	// "github.com/stackrox/acs-fleet-manager/pkg/handlers"
	"github.com/stackrox/acs-fleet-manager/pkg/shared"
	"github.com/getsentry/sentry-go"
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strings"

	"github.com/stackrox/acs-fleet-manager/pkg/client/observatorium"
	"github.com/stackrox/acs-fleet-manager/pkg/errors"
	"github.com/gorilla/mux"
)

type metricsHandler struct {
	service services.ObservatoriumService
}

func NewMetricsHandler(service services.ObservatoriumService) *metricsHandler {
	return &metricsHandler{
		service: service,
	}
}

func (h metricsHandler) FederateMetrics(w http.ResponseWriter, r *http.Request) {
	dinosaurId := strings.TrimSpace(mux.Vars(r)["id"])
	if dinosaurId == "" {
		shared.HandleError(r, w, &errors.ServiceError{
			Code:     errors.ErrorBadRequest,
			Reason:   "missing path parameter: dinosaur id",
			HttpCode: http.StatusBadRequest,
		})
		return
	}

	dinosaurMetrics := &observatorium.DinosaurMetrics{}
	params := observatorium.MetricsReqParams{
		ResultType: observatorium.Query,
	}

	_, err := h.service.GetMetricsByDinosaurId(r.Context(), dinosaurMetrics, dinosaurId, params)
	if err != nil {
		if err.Code == errors.ErrorNotFound {
			shared.HandleError(r, w, err)
		} else {
			glog.Errorf("error getting metrics: %v", err)
			sentry.CaptureException(err)
			shared.HandleError(r, w, &errors.ServiceError{
				Code:     err.Code,
				Reason:   "error getting metrics",
				HttpCode: http.StatusInternalServerError,
			})
		}
		return
	}

	// Define metric collector
	collector := metrics.NewFederatedUserMetricsCollector(dinosaurMetrics)
	registry := prometheus.NewPedanticRegistry()
	registry.MustRegister(collector)

	promHandler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		ErrorHandling: promhttp.HTTPErrorOnError,
	})
	promHandler.ServeHTTP(w, r)
}

func (h metricsHandler) GetMetricsByRangeQuery(w http.ResponseWriter, r *http.Request) {
	/* FIXME
	id := mux.Vars(r)["id"]
	params := observatorium.MetricsReqParams{}
	query := r.URL.Query()
	cfg := &handlers.HandlerConfig{
		Validate: []handlers.Validate{
			handlers.ValidatQueryParam(query, "duration"),
			handlers.ValidatQueryParam(query, "interval"),
		},
		Action: func() (i interface{}, serviceError *errors.ServiceError) {
			ctx := r.Context()
			params.ResultType = observatorium.RangeQuery
			extractMetricsQueryParams(r, &params)
			dinosaurMetrics := &observatorium.DinosaurMetrics{}
			foundDinosaurId, err := h.service.GetMetricsByDinosaurId(ctx, dinosaurMetrics, id, params)
			if err != nil {
				return nil, err
			}
			metricList := public.MetricsRangeQueryList{
				Kind: "MetricsRangeQueryList",
				Id:   foundDinosaurId,
			}
			metricsResult, err := presenters.PresentMetricsByRangeQuery(dinosaurMetrics)
			if err != nil {
				return nil, err
			}
			metricList.Items = metricsResult

			return metricList, nil
		},
	}
	handlers.HandleGet(w, r, cfg)
	*/
}

func (h metricsHandler) GetMetricsByInstantQuery(w http.ResponseWriter, r *http.Request) {
	/* FIXME
	id := mux.Vars(r)["id"]
	params := observatorium.MetricsReqParams{}
	cfg := &handlers.HandlerConfig{
		Action: func() (i interface{}, serviceError *errors.ServiceError) {
			ctx := r.Context()
			params.ResultType = observatorium.Query
			extractMetricsQueryParams(r, &params)
			dinosaurMetrics := &observatorium.DinosaurMetrics{}
			foundDinosaurId, err := h.service.GetMetricsByDinosaurId(ctx, dinosaurMetrics, id, params)
			if err != nil {
				return nil, err
			}
			metricList := public.MetricsInstantQueryList{
				Kind: "MetricsInstantQueryList",
				Id:   foundDinosaurId,
			}
			metricsResult, err := presenters.PresentMetricsByInstantQuery(dinosaurMetrics)
			if err != nil {
				return nil, err
			}
			metricList.Items = metricsResult

			return metricList, nil
		},
	}
	handlers.HandleGet(w, r, cfg)
	*/
}
