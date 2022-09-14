package app

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func convertToSlsResponse(data interface{}, statusCode int) (Response, error) {
	var buf bytes.Buffer
	body, err := json.Marshal(data)
	if err != nil {
		return Response{StatusCode: 500}, err
	}
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      statusCode,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "hello-handler",
		},
	}

	return resp, nil
}

func SendErrorfJSON(errorMessage string, args ...interface{}) (Response, error) {
	data := map[string]interface{}{
		"data":    nil,
		"success": false,
		"error":   fmt.Sprintf(errorMessage, args...),
	}

	return convertToSlsResponse(data, 200)
}

func SendAuthErrorfJSON(errorMessage string, args ...interface{}) (Response, error) {
	data := map[string]interface{}{
		"data":           nil,
		"success":        false,
		"error":          fmt.Sprintf(errorMessage, args...),
		"_un_authorized": true,
	}

	return convertToSlsResponse(data, 200)
}

func SendJSON(data interface{}) (Response, error) {
	d := map[string]interface{}{
		"data":    data,
		"success": true,
		"error":   nil,
	}

	return convertToSlsResponse(d, 200)
}

func SendPagedJSON(data interface{}, totalCount int64) (Response, error) {
	d := map[string]interface{}{
		"data":    data,
		"total":   totalCount,
		"success": true,
		"error":   nil,
	}

	return convertToSlsResponse(d, 200)
}
