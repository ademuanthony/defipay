package app

import (
	"database/sql"
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

func NewValidationError(errors []string) validationError {
	return validationError{Errors: errors}
}

func (m Module) handleError(err error, tag ...string) (Response, error) {
	msg := "Unable to complete your request. Something went wrong"
	if messenger, ok := err.(ErrorMessenger); ok {
		msg = messenger.ErrorMessage()
	} else if err.Error() == sql.ErrNoRows.Error() {
		msg = "Not Found"
	}
	log.Error(tag, err)
	return SendErrorfJSON(msg)
}
