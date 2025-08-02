package utils

// This file contains related things need to deal with when implementing request handler
// Currently only encodejson and decodejson

import (
	"encoding/json"
	"errors"
	"net/http"
)

var ErrNotMatchContentTypeJSON = errors.New("content type must be 'application/json'")

// EncodeJSON writes the given data as JSON to the provided http.ResponseWriter.
// If there is an error encoding the data, a 500 Internal Server Error response is written to the ResponseWriter.
func EncodeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// DecodeJSON reads the request body as JSON and decodes it into the provided form interface.
// If the request Content-Type header is not "application/json", it returns ErrNotMatchContentTypeJSON.
// If there is an error decoding the JSON, it returns the error.
func DecodeJSON(r *http.Request, form interface{}) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return ErrNotMatchContentTypeJSON
	}

	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		return err
	}

	return nil
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

// JSONErrorResponse writes an HTTP error response with the given status code and message as a JSON-encoded ErrorResponse.
func JSONErrorResponse(w http.ResponseWriter, statusCode int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ErrorResponse{
		Code:    statusCode,
		Error:   http.StatusText(statusCode),
		Message: msg,
	}

	EncodeJSON(w, response)
}
