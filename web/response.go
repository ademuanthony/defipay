package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

func RenderHTML(tmpl string, w http.ResponseWriter, r *http.Request, data interface{}, server *Server) {
	str, err := server.Templates.ExecTemplateToString(tmpl, data)

	if err != nil {
		log.Errorf("Template execute failure: %v", err)
		server.StatusPage(w, r, DefaultErrorCode, DefaultErrorMessage, "", ExpStatusError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	if _, err = io.WriteString(w, str); err != nil {
		log.Error(err)
	}
}

func SendErrorfJSON(res http.ResponseWriter, errorMessage string, args ...interface{}) {
	data := map[string]interface{}{
		"data":    nil,
		"success": false,
		"error":   fmt.Sprintf(errorMessage, args...),
	}

	renderJSON(res, data)
}

func sendAuthErrorfJSON(res http.ResponseWriter, errorMessage string, args ...interface{}) {
	data := map[string]interface{}{
		"data":           nil,
		"success":        false,
		"error":          fmt.Sprintf(errorMessage, args...),
		"_un_authorized": true,
	}

	renderJSON(res, data)
}

func SendJSON(res http.ResponseWriter, data interface{}) {
	d := map[string]interface{}{
		"data":    data,
		"success": true,
		"error":   nil,
	}

	renderJSON(res, d)
}

func SendPagedJSON(res http.ResponseWriter, data interface{}, totalCount int64) {
	d := map[string]interface{}{
		"data":    data,
		"total":   totalCount,
		"success": true,
		"error":   nil,
	}

	renderJSON(res, d)
}

func renderJSON(res http.ResponseWriter, data interface{}) {
	d, err := json.Marshal(data)
	if err != nil {
		log.Errorf("Error marshalling data: %s", err.Error())
	}

	res.Header().Set("Content-Type", "application/json")
	_, _ = res.Write(d)
}

// RenderJSONBytes prepares the headers for pre-encoded JSON and writes the JSON
// bytes.
func RenderJSONBytes(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err := w.Write(data)
	if err != nil {
		// Filter out broken pipe (user pressed "stop") errors
		if _, ok := err.(*net.OpError); ok {
			if strings.Contains(err.Error(), "broken pipe") {
				return
			}
		}
		log.Warnf("ResponseWriter.Write error: %v", err)
	}
}
