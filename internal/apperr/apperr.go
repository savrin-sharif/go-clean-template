// Package apperr defines typed application errors that handlers map to HTTP.
package apperr

import (
	"errors"
	"net/http"
)

// Kind classifies an application error for HTTP and client handling.
type Kind string

const (
	KindValidation   Kind = "validation"
	KindUnauthorized Kind = "unauthorized"
	KindForbidden    Kind = "forbidden"
	KindNotFound     Kind = "not_found"
	KindConflict     Kind = "conflict"
	KindData         Kind = "data"
	KindInternal     Kind = "internal"
)

const defaultMessage = "Something went wrong. Please try again."

// Error is a categorized, user-safe application error.
type Error struct {
	Kind        Kind
	Message     string
	FieldErrors map[string]string
	Cause       error
}

func (e *Error) Error() string {
	if e == nil {
		return defaultMessage
	}
	if e.Message != "" {
		return e.Message
	}
	return defaultMessage
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

// HTTPStatus returns the status code for this error kind.
func (e *Error) HTTPStatus() int {
	if e == nil {
		return http.StatusInternalServerError
	}
	switch e.Kind {
	case KindValidation:
		return http.StatusUnprocessableEntity
	case KindUnauthorized:
		return http.StatusUnauthorized
	case KindForbidden:
		return http.StatusForbidden
	case KindNotFound:
		return http.StatusNotFound
	case KindConflict:
		return http.StatusConflict
	case KindData:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

// Validation returns a validation error with optional field errors.
func Validation(message string, fieldErrors map[string]string) *Error {
	if message == "" {
		message = "Validation failed."
	}
	return &Error{Kind: KindValidation, Message: message, FieldErrors: fieldErrors}
}

// Unauthorized returns a 401-class error.
func Unauthorized(message string) *Error {
	if message == "" {
		message = "Authentication required."
	}
	return &Error{Kind: KindUnauthorized, Message: message}
}

// Forbidden returns a 403-class error.
func Forbidden(message string) *Error {
	if message == "" {
		message = "You do not have permission to perform this action."
	}
	return &Error{Kind: KindForbidden, Message: message}
}

// NotFound returns a 404-class error.
func NotFound(message string) *Error {
	if message == "" {
		message = "Resource not found."
	}
	return &Error{Kind: KindNotFound, Message: message}
}

// Conflict returns a 409-class error.
func Conflict(message string) *Error {
	if message == "" {
		message = "Resource conflict."
	}
	return &Error{Kind: KindConflict, Message: message}
}

// Data returns an error for malformed or unreadable data.
func Data(message string) *Error {
	if message == "" {
		message = "Invalid data."
	}
	return &Error{Kind: KindData, Message: message}
}

// Internal returns a 500-class error. Prefer wrapping with WrapInternal.
func Internal(message string) *Error {
	if message == "" {
		message = defaultMessage
	}
	return &Error{Kind: KindInternal, Message: message}
}

// WrapInternal wraps an unexpected cause as an internal error.
func WrapInternal(cause error, message string) *Error {
	if message == "" {
		message = defaultMessage
	}
	return &Error{Kind: KindInternal, Message: message, Cause: cause}
}

// As extracts an *Error from err, if present.
func As(err error) (*Error, bool) {
	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

// From returns err as *Error when possible; otherwise wraps it as internal.
func From(err error) *Error {
	if err == nil {
		return nil
	}
	if appErr, ok := As(err); ok {
		return appErr
	}
	return WrapInternal(err, defaultMessage)
}
