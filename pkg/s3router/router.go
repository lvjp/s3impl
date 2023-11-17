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
		logger:  logger,
		factory: ContextFactory{Hosts: hosts},
	}
}

type handler struct {
	logger  *zerolog.Logger
	factory ContextFactory
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.NewString()

	context := h.factory.Build(r)
	h.logger.Trace().
		Interface("context", context).
		Msg("Context decoded")

	resp := &s3errors.S3Error{
		HTTPStatusCode: http.StatusNotImplemented,
		Code:           "NotImplemented",
		Message:        "A header that you provided implies functionality that is not implemented.",
		RequestID:      requestID,
		Resource:       r.URL.String(),
	}

	writer := s3errors.APIWriter{}

	if err := writer.Write(resp, w); err != nil && !errors.Is(err, io.ErrClosedPipe) {
		h.logger.Warn().Err(err).Msg("Cannot write response")
	}
}
