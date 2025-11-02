package http

import (
	"errors"
	"net/http"

	"github.com/octokerbs/chronocode-backend/internal/domain"
)

type APIError struct {
	Status  int
	Message string
}

func FromError(err error) APIError {
	var apiError APIError
	var svcError domain.Error

	if errors.As(err, &svcError) {
		svcErr := svcError.Category()
		switch svcErr {
		case domain.ErrBadRequest:
			apiError.Status = http.StatusBadRequest
		case domain.ErrInternalFailure:
			apiError.Status = http.StatusInternalServerError
		}
		apiError.Message = svcError.Specific().Error()
	} else {
		apiError.Status = http.StatusInternalServerError
		apiError.Message = "An unexpected error occurred"
	}

	return apiError
}
