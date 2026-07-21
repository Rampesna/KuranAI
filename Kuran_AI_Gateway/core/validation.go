package core

import (
	"net/http"

	"github.com/go-playground/validator/v10"
)

// ValidationErrorResponse binding/validation hatasını 422 + alan hataları olarak döner.
func ValidationErrorResponse(validationErr error) *ServiceResponse {
	if validationErr == nil {
		return nil
	}
	fields := make(map[string]string)
	if validationErrors, ok := validationErr.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {
			fields[fieldErr.Field()] = validationMessage(fieldErr.Tag(), fieldErr.Param())
		}
		return NewServiceResponse(false, "Validation failed", http.StatusUnprocessableEntity, fields)
	}
	return NewServiceResponse(false, validationErr.Error(), http.StatusBadRequest, nil)
}

func validationMessage(tag, param string) string {
	switch tag {
	case "required":
		return "This field is required."
	case "email":
		return "Must be a valid email."
	case "min":
		return "Minimum length is " + param + "."
	case "max":
		return "Maximum length is " + param + "."
	default:
		return "Invalid value."
	}
}
