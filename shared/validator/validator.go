package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateStruct validates a struct and returns formatted error messages
func ValidateStruct(s any) error {
	if err := validate.Struct(s); err != nil {
		var validationErrors []string

		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, formatValidationError(err))
		}

		return fmt.Errorf("%s", strings.Join(validationErrors, "; "))
	}
	return nil
}

// formatValidationError formats individual validation errors into human-readable messages
func formatValidationError(fe validator.FieldError) string {
	field := strings.ToLower(fe.Field())

	switch fe.Tag() {	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid mail address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", field, fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", field, fe.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters long", field, fe.Param())
	case "alpha":
		return fmt.Sprintf("%s must contain only alphabetic characters", field)
	case "alphanum":
		return fmt.Sprintf("%s must contain only alphanumeric characters", field)
	case "numeric":
		return fmt.Sprintf("%s must be numeric", field)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, fe.Param())
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "contains":
		return fmt.Sprintf("%s must contain '%s'", field, fe.Param())
	case "excludes":
		return fmt.Sprintf("%s cannot contain '%s'", field, fe.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, fe.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, fe.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, fe.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, fe.Param())
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

// IsValidMail checks if the given string is a valid mail address
func IsValidMail(mail string) bool {
	err := validate.Var(mail, "required,email")
	return err == nil
}

// IsValidPassword checks if the password meets minimum requirements
func IsValidPassword(password string) bool {
	err := validate.Var(password, "required,min=8")
	return err == nil
}

// ValidatePassword validates password with custom rules
func ValidatePassword(password string) error {
	if len(strings.TrimSpace(password)) == 0 {
		return fmt.Errorf("password is required")
	}

	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	if len(password) > 100 {
		return fmt.Errorf("password must be at most 100 characters long")
	}

	// Check for at least one letter and one number (optional enhancement)
	hasLetter := false
	hasNumber := false

	for _, char := range password {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
			hasLetter = true
		}
		if char >= '0' && char <= '9' {
			hasNumber = true
		}
	}

	if !hasLetter {
		return fmt.Errorf("password must contain at least one letter")
	}

	if !hasNumber {
		return fmt.Errorf("password must contain at least one number")
	}

	return nil
}
