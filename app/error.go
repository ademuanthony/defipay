package app

import (
	"deficonnect/defipayapi/web"
	"net/http"
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

func (m module) handleError(w http.ResponseWriter, err error, tag ...string) {
	msg := "Cannot update currency. Something went wrong"
	if messenger, ok := err.(ErrorMessenger); ok {
		msg = messenger.ErrorMessage()
	}
	web.SendErrorfJSON(w, msg)
	log.Error(tag, err)
}
