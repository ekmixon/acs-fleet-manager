package shared

import (
	"net/http"

	"github.com/stackrox/acs-fleet-manager/pkg/errors"
	"github.com/stackrox/acs-fleet-manager/pkg/logger"
)

// HandleError handles a service error by returning an appropriate HTTP response with error reason
func HandleError(r *http.Request, w http.ResponseWriter, err *errors.ServiceError) {
	ctx := r.Context()
	ulog := logger.NewUHCLogger(ctx)
	operationID := logger.GetOperationID(ctx)
	if err.HttpCode >= 400 && err.HttpCode <= 499 {
		ulog.Infof(err.Error())
	} else {
		ulog.Error(err)
	}

	WriteJSONResponse(w, err.HttpCode, err.AsOpenapiError(operationID, r.RequestURI))
}
