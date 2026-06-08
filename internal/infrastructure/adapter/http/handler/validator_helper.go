package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/cristiansrc/hv-go-ms-resume/internal/infrastructure/config"
)

// decodeAndValidate decodes JSON body and validates it, writing error response if invalid.
// Returns true if valid, false if error was written.
func decodeAndValidate(w http.ResponseWriter, r *http.Request, dst interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		WriteError(w, r, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return false
	}
	if err := config.ValidateStruct(dst); err != nil {
		details := parseValidationErrors(err)
		WriteValidationError(w, r, details)
		return false
	}
	return true
}

// parseValidationErrors converts validator.ValidationErrors to []ApiErrorDetail
func parseValidationErrors(err error) []ApiErrorDetail {
	var details []ApiErrorDetail
	if verr, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range verr {
			code := "FIELD_REQUIRED"
			if fe.ActualTag() != "required" {
				code = "INVALID_" + fe.ActualTag()
			}
			details = append(details, ApiErrorDetail{
				Field:   fe.Field(),
				Code:    code,
				Message: fe.Field() + " is " + fe.ActualTag(),
			})
		}
	}
	return details
}
