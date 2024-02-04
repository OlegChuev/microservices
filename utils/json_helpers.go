package utils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// Config represents the configuration for the utility functions in this package.
type Config struct{}

// JsonResponse defines the structure of the JSON response used in the utility functions.
type JsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// ReadJson reads and decodes JSON data from the request body, ensuring a single JSON value.
func (app *Config) ReadJson(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1 * 1024 * 1024 // 1 MB

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}

	return nil
}

// WriteJson marshals data to JSON and writes it to the response writer with the specified status and optional headers.
func (app *Config) WriteJson(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

// ErrorJson writes a JSON response with an error message and status code (default is http.StatusBadRequest).
func (app *Config) ErrorJson(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	payload := JsonResponse{
		Error:   true,
		Message: err.Error(),
	}

	return app.WriteJson(w, statusCode, payload)
}

// SuccessJson creates a JsonResponse with a success message and no error.
func (app *Config) SuccessJson(msg string) JsonResponse {
	return JsonResponse{
		Error:   false,
		Message: msg,
	}
}
