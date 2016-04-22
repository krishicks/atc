package resourceserver

import (
	"encoding/json"
	"net/http"

	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/resource"
	"github.com/pivotal-golang/lager"
	"github.com/tedsuo/rata"
)

func (s *Server) CheckResource(pipelineDB db.PipelineDB) http.Handler {
	logger := s.logger.Session("check-resource")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resourceName := rata.Param(r, "resource_name")

		var reqBody atc.CheckRequestBody
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			logger.Info("malformed-request", lager.Data{"error": err.Error()})
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		scanner := s.scannerFactory.NewResourceScanner(pipelineDB)

		err = scanner.ScanFromVersion(logger, resourceName, reqBody.From)
		switch scanErr := err.(type) {
		case resource.ErrResourceScriptFailed:
			checkResponseBody := atc.CheckResponseBody{
				ExitStatus: scanErr.ExitStatus,
				Stderr:     scanErr.Stderr,
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(checkResponseBody)
		case error:
			w.WriteHeader(http.StatusInternalServerError)
		default:
			w.WriteHeader(http.StatusOK)
		}
	})
}