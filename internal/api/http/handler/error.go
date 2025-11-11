package handler

import (
	"errors"
	"net/http"

	pkg_errors "github.com/octokerbs/chronocode-backend/internal/errors"
)

type APIError struct {
	Status  int
	Message string
}

func FromError(err error) APIError {
	var apiError APIError
	var svcError pkg_errors.Error

	if errors.As(err, &svcError) {
		svcErr := svcError.Category()
		switch svcErr {
		case pkg_errors.ErrBadRequest:
			apiError.Status = http.StatusBadRequest
		case pkg_errors.ErrInternalFailure:
			apiError.Status = http.StatusInternalServerError
		}
		apiError.Message = svcError.Specific().Error()
	} else {
		apiError.Status = http.StatusInternalServerError
		apiError.Message = "An unexpected error occurred"
	}

	return apiError
}
