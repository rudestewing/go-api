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



// ValidateStruct validates a struct and returns validation errors
func ValidateStruct(s any) map[string][]string {
	if err := validate.Struct(s); err != nil {
		errors := make(map[string][]string)

		// Check if error is of type ValidationErrors
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldError := range validationErrors {
				field := getFieldName(fieldError)
				message := formatValidationError(fieldError)
				
				if _, exists := errors[field]; !exists {
					errors[field] = []string{}
				}
				errors[field] = append(errors[field], message)
			}
		} else {
			// Handle other types of validation errors
			errors["general"] = []string{err.Error()}
		}

		return errors
	}
	return nil
}

// getFieldName converts field name to dot notation for nested objects
func getFieldName(fe validator.FieldError) string {
	// Get the full namespace and convert to lowercase dot notation
	namespace := fe.Namespace()
	
	// Remove the root struct name (first part before the first dot)
	parts := strings.Split(namespace, ".")
	if len(parts) > 1 {
		// Skip the first part (struct name) and join the rest
		fieldParts := parts[1:]
		// Convert to lowercase and join with dots
		for i, part := range fieldParts {
			fieldParts[i] = strings.ToLower(part)
		}
		return strings.Join(fieldParts, ".")
	}
	
	// If no namespace, just return the lowercase field name
	return strings.ToLower(fe.Field())
}

// formatValidationError formats individual validation errors into human-readable messages
func formatValidationError(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Must be a valid email address"
	case "min":
		return fmt.Sprintf("Must be at least %s characters long", fe.Param())
	case "max":
		return fmt.Sprintf("Must be at most %s characters long", fe.Param())
	case "len":
		return fmt.Sprintf("Must be exactly %s characters long", fe.Param())
	case "alpha":
		return "Must contain only alphabetic characters"
	case "alphanum":
		return "Must contain only alphanumeric characters"
	case "numeric":
		return "Must be numeric"
	case "oneof":
		return fmt.Sprintf("Must be one of: %s", fe.Param())
	case "url":
		return "Must be a valid URL"
	case "contains":
		return fmt.Sprintf("Must contain '%s'", fe.Param())
	case "excludes":
		return fmt.Sprintf("Cannot contain '%s'", fe.Param())
	case "gte":
		return fmt.Sprintf("Must be greater than or equal to %s", fe.Param())
	case "lte":
		return fmt.Sprintf("Must be less than or equal to %s", fe.Param())
	case "gt":
		return fmt.Sprintf("Must be greater than %s", fe.Param())
	case "lt":
		return fmt.Sprintf("Must be less than %s", fe.Param())
	default:
		return "This field is invalid"
	}
}

// IsValidEmail checks if the given string is a valid email address
func IsValidEmail(email string) bool {
	err := validate.Var(email, "required,email")
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

// ValidatePasswordWithDetails validates password and returns structured error
func ValidatePasswordWithDetails(password string) map[string][]string {
	var messages []string

	if len(strings.TrimSpace(password)) == 0 {
		messages = append(messages, "This field is required")
	} else {
		if len(password) < 8 {
			messages = append(messages, "Must be at least 8 characters long")
		}

		if len(password) > 100 {
			messages = append(messages, "Must be at most 100 characters long")
		}

		// Check for at least one letter and one number
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
			messages = append(messages, "Must contain at least one letter")
		}

		if !hasNumber {
			messages = append(messages, "Must contain at least one number")
		}
	}

	if len(messages) > 0 {
		errors := make(map[string][]string)
		errors["password"] = messages
		return errors
	}

	return nil
}
