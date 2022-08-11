package presenters

import (
	private "github.com/stackrox/acs-fleet-manager/generated/privateapi"
	"github.com/stackrox/acs-fleet-manager/internal/dinosaur/pkg/api/dbapi"
)

// ConvertDataPlaneDinosaurStatus ...
func ConvertDataPlaneDinosaurStatus(status map[string]private.CentralStatus) []*dbapi.DataPlaneCentralStatus {
	res := make([]*dbapi.DataPlaneCentralStatus, 0, len(status))

	for k, v := range status {
		c := make([]dbapi.DataPlaneCentralStatusCondition, 0, len(v.Conditions))
		var routes []dbapi.DataPlaneCentralRoute
		for _, s := range v.Conditions {
			c = append(c, dbapi.DataPlaneCentralStatusCondition{
				Type:    s.Type,
				Reason:  s.Reason,
				Status:  s.Status,
				Message: s.Message,
			})
		}
		if v.Routes != nil {
			routes = make([]dbapi.DataPlaneCentralRoute, 0, len(v.Routes))
			for _, ro := range v.Routes {
				routes = append(routes, dbapi.DataPlaneCentralRoute{
					Domain: ro.Domain,
					Router: ro.Router,
				})
			}
		}

		dbStatus := &dbapi.DataPlaneCentralStatus{
			CentralClusterID: k,
			Conditions:       c,
			Routes:           routes,
		}
		if v.Versions != nil {
			dbStatus.CentralVersion = v.Versions.Central
			dbStatus.CentralOperatorVersion = v.Versions.CentralOperator
		}
		res = append(res, dbStatus)
	}

	return res
}

// ConvertDataPlaneDinosaurStatusFromGRPC converts GRPC requests to dbapi
func ConvertDataPlaneDinosaurStatusFromGRPC(from *private.UpdateCentralStatusRequest) []*dbapi.DataPlaneCentralStatus {
	res := make([]*dbapi.DataPlaneCentralStatus, 0, len(from.Updates))

	for k, v := range from.Updates {
		c := make([]dbapi.DataPlaneCentralStatusCondition, 0, len(v.Conditions))
		var routes []dbapi.DataPlaneCentralRoute
		for _, s := range v.Conditions {
			c = append(c, dbapi.DataPlaneCentralStatusCondition{
				Type:    s.Type,
				Reason:  s.Reason,
				Status:  s.Status,
				Message: s.Message,
			})
		}
		if v.Routes != nil {
			routes = make([]dbapi.DataPlaneCentralRoute, 0, len(v.Routes))
			for _, ro := range v.Routes {
				routes = append(routes, dbapi.DataPlaneCentralRoute{
					Domain: ro.Domain,
					Router: ro.Router,
				})
			}
		}

		dbStatus := &dbapi.DataPlaneCentralStatus{
			CentralClusterID: k,
			Conditions:       c,
			Routes:           routes,
		}

		if v.Versions != nil {
			dbStatus.CentralVersion = v.Versions.Central
			dbStatus.CentralOperatorVersion = v.Versions.CentralOperator
		}

		res = append(res, dbStatus)
	}

	return res
}
