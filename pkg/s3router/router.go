package s3router

import (
	"errors"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/lvjp/s3impl/pkg/s3errors"
	"github.com/rs/zerolog"
)

func New(logger *zerolog.Logger, hosts []string) http.Handler {
	return &handler{
		logger: logger,
		hosts:  hosts,
	}
}

type handler struct {
	logger *zerolog.Logger
	hosts  []string
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.NewString()

	resp := &s3errors.S3Error{
		HTTPStatusCode: http.StatusNotImplemented,
		Code:           "NotImplemented",
		Message:        "A header that you provided implies functionality that is not implemented.",
		RequestID:      requestID,
		Resource:       r.URL.String(),
	}

	route, err := DetermineRoute(r, h.hosts)
	if err != nil {
		resp.HTTPStatusCode = http.StatusBadRequest
		resp.Code = "Badrequest"
		resp.Message = err.Error()
	} else {
		h.logger.Trace().
			Interface("route", route).
			Msg("Route determinated")
	}

	writer := s3errors.APIWriter{}

	if err := writer.Write(resp, w); err != nil && !errors.Is(err, io.ErrClosedPipe) {
		h.logger.Warn().Err(err).Msg("Cannot write response")
	}
}
