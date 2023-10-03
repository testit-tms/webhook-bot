package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ErrorResponse represents an error response returned by the API.
type ErrorResponse struct {
	Message string `json:"message"`
}

// NewErrorResponse writes an error response to the provided http.ResponseWriter with the given status code and error message.
func NewErrorResponse(w http.ResponseWriter, status int, err string) {
	w.WriteHeader(status)
	// nolint:errcheck
	json.NewEncoder(w).Encode(ErrorResponse{
		Message: err,
	})
}

// ValidationError generates a string message from the provided validation errors.
func ValidationError(errs validator.ValidationErrors) string {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "parse-mode":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s must be empty or have following value: markdownv2, markdown or html", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return strings.Join(errMsgs, ", ")
}
