package validator

import (
	"fmt"
	"go-api/app/shared/errors"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	
	// Register custom tag name function to use json tags
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	
	// Register custom validations
	registerCustomValidations()
}

// ValidateStruct validates a struct and returns a ValidationError with field-specific errors
func ValidateStruct(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		fields := make(map[string][]string)
		
		for _, err := range validationErrors {
			fieldName := err.Field()
			if fieldName == "" {
				fieldName = strings.ToLower(err.StructField())
			}
			
			errorMessage := formatValidationError(err)
			fields[fieldName] = append(fields[fieldName], errorMessage)
		}
		
		return errors.NewValidationError("The given data was invalid", fields)
	}
	return nil
}

// ValidateVar validates a single variable
func ValidateVar(field interface{}, tag string) error {
	return validate.Var(field, tag)
}

// formatValidationError formats individual validation errors into human-readable messages
func formatValidationError(fe validator.FieldError) string {
	field := fe.Field()
	if field == "" {
		field = strings.ToLower(fe.StructField())
	}

	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("The %s field is required", field)
	case "email":
		return fmt.Sprintf("The %s field must be a valid email address", field)
	case "min":
		if fe.Kind().String() == "string" {
			return fmt.Sprintf("The %s field must be at least %s characters", field, fe.Param())
		}
		return fmt.Sprintf("The %s field must be at least %s", field, fe.Param())
	case "max":
		if fe.Kind().String() == "string" {
			return fmt.Sprintf("The %s field must not exceed %s characters", field, fe.Param())
		}
		return fmt.Sprintf("The %s field must not exceed %s", field, fe.Param())
	case "len":
		return fmt.Sprintf("The %s field must be exactly %s characters", field, fe.Param())
	case "alpha":
		return fmt.Sprintf("The %s field must contain only alphabetic characters", field)
	case "alphanum":
		return fmt.Sprintf("The %s field must contain only alphanumeric characters", field)
	case "numeric":
		return fmt.Sprintf("The %s field must be numeric", field)
	case "oneof":
		return fmt.Sprintf("The %s field must be one of: %s", field, strings.ReplaceAll(fe.Param(), " ", ", "))
	case "url":
		return fmt.Sprintf("The %s field must be a valid URL", field)
	case "uri":
		return fmt.Sprintf("The %s field must be a valid URI", field)
	case "contains":
		return fmt.Sprintf("The %s field must contain '%s'", field, fe.Param())
	case "excludes":
		return fmt.Sprintf("The %s field cannot contain '%s'", field, fe.Param())
	case "gte":
		return fmt.Sprintf("The %s field must be greater than or equal to %s", field, fe.Param())
	case "lte":
		return fmt.Sprintf("The %s field must be less than or equal to %s", field, fe.Param())
	case "gt":
		return fmt.Sprintf("The %s field must be greater than %s", field, fe.Param())
	case "lt":
		return fmt.Sprintf("The %s field must be less than %s", field, fe.Param())
	case "eq":
		return fmt.Sprintf("The %s field must equal %s", field, fe.Param())
	case "ne":
		return fmt.Sprintf("The %s field must not equal %s", field, fe.Param())
	case "uuid":
		return fmt.Sprintf("The %s field must be a valid UUID", field)
	case "uuid4":
		return fmt.Sprintf("The %s field must be a valid UUID v4", field)
	case "strong_password":
		return fmt.Sprintf("The %s field must contain at least 8 characters with uppercase, lowercase, and number", field)
	case "phone":
		return fmt.Sprintf("The %s field must be a valid phone number", field)
	case "indonesian_phone":
		return fmt.Sprintf("The %s field must be a valid Indonesian phone number", field)
	default:
		return fmt.Sprintf("The %s field is invalid", field)
	}
}

// registerCustomValidations registers custom validation rules
func registerCustomValidations() {
	// Strong password validation
	validate.RegisterValidation("strong_password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		return isStrongPassword(password)
	})
	
	// Indonesian phone number validation
	validate.RegisterValidation("indonesian_phone", func(fl validator.FieldLevel) bool {
		phone := fl.Field().String()
		return isIndonesianPhone(phone)
	})
}

// isStrongPassword checks if password meets strong password requirements
func isStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	
	var hasUpper, hasLower, hasNumber bool
	
	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasNumber = true
		}
	}
	
	return hasUpper && hasLower && hasNumber
}

// isIndonesianPhone validates Indonesian phone number format
func isIndonesianPhone(phone string) bool {
	// Remove common formatting characters
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, "(", "")
	phone = strings.ReplaceAll(phone, ")", "")
	
	// Check if it starts with +62 or 62 or 08
	if strings.HasPrefix(phone, "+62") {
		phone = phone[3:]
	} else if strings.HasPrefix(phone, "62") {
		phone = phone[2:]
	} else if strings.HasPrefix(phone, "08") {
		phone = phone[1:]
	} else if strings.HasPrefix(phone, "8") {
		// Already in correct format
	} else {
		return false
	}
	
	// Check if remaining digits are numeric and proper length
	if len(phone) < 8 || len(phone) > 13 {
		return false
	}
	
	for _, char := range phone {
		if char < '0' || char > '9' {
			return false
		}
	}
		return true
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
