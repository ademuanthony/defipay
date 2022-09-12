package app

import "strings"

type ErrorMessenger interface {
	ErrorMessage() string
}

type validationError struct {
	Errors []string `json:"errors"`
}

func (v validationError) Error() string {
	return strings.Join(v.Errors, "|")
}

func (v validationError) ErrorMessage() string {
	return strings.Join(v.Errors, "|")
}

func newValidationError(errors []string) validationError {
	return validationError{Errors: errors}
}
