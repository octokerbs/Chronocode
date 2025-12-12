package errs

import (
	"errors"
	"net/http"
)

// Client Category Errors (4xx)

// ErrBadRequest indicates that the client request structure is incorrect
// (e.g., malformed JSON, missing required fields, or invalid data types).
// Typically maps to HTTP 400 Bad Request.
var ErrBadRequest = errors.New("bad request")

// ErrUnauthorized indicates that the client has not provided valid credentials
// or is not authenticated (e.g., missing JWT token or the token is invalid/expired).
// Typically maps to HTTP 401 Unauthorized.
var ErrUnauthorized = errors.New("unauthorized")

// ErrForbidden indicates that the client is authenticated but lacks the necessary
// permissions to access the resource or execute the requested action (ACL/RBAC).
// Typically maps to HTTP 403 Forbidden.
var ErrForbidden = errors.New("forbidden")

// ErrNotFound indicates that the resource requested by the client (e.g., a user,
// a note) was not found in the system.
// Typically maps to HTTP 404 Not Found.
var ErrNotFound = errors.New("not found")

// ErrConflict indicates that the attempt to create or update a resource violated a
// uniqueness constraint or a resource state (e.g., email already registered).
// Typically maps to HTTP 409 Conflict.
var ErrConflict = errors.New("conflict")

// ErrValidationFailed indicates that the input data is structurally valid but
// fails business rules or semantic validations (e.g., password too weak).
// Typically maps to HTTP 422 Unprocessable Entity.
var ErrValidationFailed = errors.New("validation failed")

// ErrRateLimitExceeded indicates that the client has exceeded the maximum number
// of allowed requests in a period of time.
// Typically maps to HTTP 429 Too Many Requests.
var ErrRateLimitExceeded = errors.New("rate limit exceeded")

// Server Category Errors (5xx)

// ErrInternalFailure indicates an unexpected, unclassified error on the server side
// (e.g., an unhandled logic failure or a panic).
// Typically maps to HTTP 500 Internal Server Error.
var ErrInternalFailure = errors.New("internal failure")

// ErrDependencyFailed indicates that a critical external dependency (e.g., another
// service, a payment system) failed or returned an unexpected error.
// Typically maps to HTTP 500 Internal Server Error or 503 Service Unavailable.
var ErrDependencyFailed = errors.New("dependency failed")

// ErrUnavailable indicates that the service or a critical dependency is
// temporarily offline or under maintenance.
// Typically maps to HTTP 503 Service Unavailable.
var ErrUnavailable = errors.New("service unavailable")

// The idea of the package is for it to be used at the infrastructure/domain level.
// Given the domain error (e.g., Cannot divide by zero) or the implementation error
// (e.g., ERROR: duplicate key value violates unique constraint "uni_user_username" (SQLSTATE 23505)),
// we can categorize the errors with a generalization. Then, the handler compares
// the category with an HTTP or gRPC status code and determines which one to use.
// For example, ErrNotFound maps to HTTP status code 404. This helps us map the API
// result from the deeper layers. The orchestrator's job is ONLY TO RETURN THE ERROR.
// The orchestrator does not know how the lower layers are implemented or what errors
// might be generated. "Something failed? Just propagate it upwards."

type Error struct {
	category error // One of the predefined error categories (ErrBadRequest, ErrInternalFailure, etc.)
	specific error // Original error context (from DB, validation, etc.)
}

func NewError(category, specific error) error {
	return Error{
		category: category,
		specific: specific,
	}
}

func (e Error) Category() error {
	return e.category
}

func (e Error) Specific() error {
	return e.specific
}

func (e Error) Error() string {
	return errors.Join(e.category, e.specific).Error()
}

type HTTPError struct {
	Status  int
	Message string
}

// ToHttpErr translates a service/domain layer error (Error struct)
// to a standardized HTTP response (HTTPError).
func ToHttpErr(err error) HTTPError {
	var appErr Error
	var httpErr HTTPError

	if errors.As(err, &appErr) {
		appErrCategory := appErr.Category()

		switch appErrCategory {

		// Client Category Errors (4xx)
		case ErrBadRequest:
			httpErr.Status = http.StatusBadRequest // 400
		case ErrUnauthorized:
			httpErr.Status = http.StatusUnauthorized // 401
		case ErrForbidden:
			httpErr.Status = http.StatusForbidden // 403
		case ErrNotFound:
			httpErr.Status = http.StatusNotFound // 404
		case ErrConflict:
			httpErr.Status = http.StatusConflict // 409
		case ErrValidationFailed:
			httpErr.Status = http.StatusUnprocessableEntity // 422
		case ErrRateLimitExceeded:
			httpErr.Status = http.StatusTooManyRequests // 429

		// Server Category Errors (5xx)
		case ErrInternalFailure:
			httpErr.Status = http.StatusInternalServerError // 500
		case ErrDependencyFailed:
			// Using 503 is often better for dependency failures to suggest temporary issues
			httpErr.Status = http.StatusServiceUnavailable // 503
		case ErrUnavailable:
			httpErr.Status = http.StatusServiceUnavailable // 503

		default:
			httpErr.Status = http.StatusInternalServerError
		}

		// Use the specific error message for the HTTP response body
		httpErr.Message = appErr.Specific().Error()

	} else {
		// Handle non-categorized Go errors as a generic internal server error
		httpErr.Status = http.StatusInternalServerError
		httpErr.Message = "An unexpected error occurred"
	}

	return httpErr
}
