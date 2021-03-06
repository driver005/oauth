package helpers

import (
	"net/http"

	"github.com/ory/fosite"
	"github.com/ory/x/logrusx"
)

var (
	ErrNotFound = &fosite.RFC6749Error{
		CodeField:        http.StatusNotFound,
		ErrorField:       http.StatusText(http.StatusNotFound),
		DescriptionField: "Unable to located the requested resource",
	}
	ErrConflict = &fosite.RFC6749Error{
		CodeField:        http.StatusConflict,
		ErrorField:       http.StatusText(http.StatusConflict),
		DescriptionField: "Unable to process the requested resource because of conflict in the current state",
	}
)

func LogError(r *http.Request, err error, logger *logrusx.Logger) {
	if logger == nil {
		logger = logrusx.New("", "")
	}

	logger.WithRequest(r).
		WithError(err).Errorln("An error occurred")
}
