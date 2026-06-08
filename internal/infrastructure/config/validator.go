package config

import "github.com/go-playground/validator/v10"

var validate = validator.New()

// ValidateStruct validates a struct using the validator tags.
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}
