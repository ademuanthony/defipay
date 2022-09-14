package app

import (
	"strings"
)

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

func (m Module) handleError(err error, tag ...string) (Response, error) {
	msg := "Cannot update currency. Something went wrong"
	if messenger, ok := err.(ErrorMessenger); ok {
		msg = messenger.ErrorMessage()
	}
	log.Error(tag, err)
	return SendErrorfJSON(msg)
}
