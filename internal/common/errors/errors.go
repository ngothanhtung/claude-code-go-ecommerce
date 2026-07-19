package errors

import (
	"errors"
	"net/http"
)

// AppError is the canonical error type returned by services.
type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
	HTTP    int    `json:"-"`
	Cause   error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error { return e.Cause }

func newAppError(code, httpCode int, message string, cause error) *AppError {
	return &AppError{Code: code, Message: message, HTTP: httpCode, Cause: cause}
}

func NewInternal(message string, cause error) *AppError {
	return newAppError(CodeInternal, http.StatusInternalServerError, message, cause)
}

func NewValidation(message string) *AppError {
	return newAppError(CodeValidation, http.StatusBadRequest, message, nil)
}

func NewUnauthorized(message string) *AppError {
	return newAppError(CodeUnauthorized, http.StatusUnauthorized, message, nil)
}

func NewForbidden(message string) *AppError {
	return newAppError(CodeForbidden, http.StatusForbidden, message, nil)
}

func NewNotFound(message string) *AppError {
	return newAppError(CodeNotFound, http.StatusNotFound, message, nil)
}

func NewConflict(message string) *AppError {
	return newAppError(CodeConflict, http.StatusConflict, message, nil)
}

func NewInvalidCredentials() *AppError {
	return newAppError(CodeInvalidCredentials, http.StatusUnauthorized, "invalid email or password", nil)
}

func NewTokenExpired() *AppError {
	return newAppError(CodeTokenExpired, http.StatusUnauthorized, "token expired", nil)
}

func NewTokenInvalid() *AppError {
	return newAppError(CodeTokenInvalid, http.StatusUnauthorized, "invalid token", nil)
}

func NewEmailExists() *AppError {
	return newAppError(CodeEmailExists, http.StatusConflict, "email already exists", nil)
}

func NewRateLimited() *AppError {
	return newAppError(CodeRateLimited, http.StatusTooManyRequests, "too many requests", nil)
}

func NewUploadTooLarge() *AppError {
	return newAppError(CodeUploadTooLarge, http.StatusBadRequest, "file too large", nil)
}

func NewUploadUnsupported() *AppError {
	return newAppError(CodeUploadUnsupported, http.StatusBadRequest, "unsupported file type", nil)
}

// AsAppError extracts an *AppError from err, mapping unknown errors to 500.
func AsAppError(err error) *AppError {
	if err == nil {
		return nil
	}
	var ae *AppError
	if errors.As(err, &ae) {
		return ae
	}
	return NewInternal("internal server error", err)
}
