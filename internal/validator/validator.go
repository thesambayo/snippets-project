package validator

import (
	"strings"
	"unicode/utf8"
)

// Define a new Validator type which contains a map of validation errors for our form fields.
type Validator struct {
	FieldErrors map[string]string
}

// Valid() returns true if the FieldErrors map doesn't contain any entries.
func (validator *Validator) Valid() bool {
	return len(validator.FieldErrors) == 0
}

// AddFieldError() adds an error message to the FieldErrors map (so long as no entry already exists for the given key).
func (validator *Validator) AddFieldError(key, message string) {
	// Note: We need to initialize the map first, if it isn't already initialized.
	if validator.FieldErrors == nil {
		validator.FieldErrors = make(map[string]string)
	}

	if _, exists := validator.FieldErrors[key]; !exists {
		validator.FieldErrors[key] = message
	}
}

// CheckField() adds an error message to the FieldErrors map only if a validation check is not 'ok'.
func (validator *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		validator.AddFieldError(key, message)
	}
}

// NotBlank() returns true if a value is not an empty string.
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxChars() returns true if a value contains no more than n characters.
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// PermittedInt() returns true if a value is in a list of permitted integers.
func PermittedInt(value int, permittedValues ...int) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}
