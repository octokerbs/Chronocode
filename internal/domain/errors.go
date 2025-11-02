package domain

import "errors"

var (
	ErrBadRequest      = errors.New("bad request")
	ErrInternalFailure = errors.New("internal failure")
	ErrNotFound        = errors.New("not found")
	ErrUnauthorized    = errors.New("unauthorized")
)

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
