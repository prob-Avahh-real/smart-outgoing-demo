package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// ValidationError represents validation errors
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   interface{} `json:"value,omitempty"`
}

// ValidationErrors represents multiple validation errors
type ValidationErrors []ValidationError

// Error implements the error interface
func (ve ValidationErrors) Error() string {
	var messages []string
	for _, err := range ve {
		messages = append(messages, fmt.Sprintf("%s: %s", err.Field, err.Message))
	}
	return strings.Join(messages, "; ")
}

// Validator provides validation functionality
type Validator struct {
	errors ValidationErrors
}

// NewValidator creates a new validator
func NewValidator() *Validator {
	return &Validator{errors: make(ValidationErrors, 0)}
}

// Required validates that a field is not empty
func (v *Validator) Required(field string, value interface{}) *Validator {
	if value == nil || value == "" {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: "field is required",
		})
	}
	return v
}

// MinLength validates minimum string length
func (v *Validator) MinLength(field string, value string, min int) *Validator {
	if len(value) < min {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: fmt.Sprintf("must be at least %d characters", min),
			Value:   value,
		})
	}
	return v
}

// MaxLength validates maximum string length
func (v *Validator) MaxLength(field string, value string, max int) *Validator {
	if len(value) > max {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: fmt.Sprintf("must be at most %d characters", max),
			Value:   value,
		})
	}
	return v
}

// Range validates numeric range
func (v *Validator) Range(field string, value float64, min, max float64) *Validator {
	if value < min || value > max {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: fmt.Sprintf("must be between %.2f and %.2f", min, max),
			Value:   value,
		})
	}
	return v
}

// Email validates email format
func (v *Validator) Email(field string, value string) *Validator {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(value) {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: "must be a valid email address",
			Value:   value,
		})
	}
	return v
}

// Coordinates validates latitude and longitude
func (v *Validator) Coordinates(latField, lngField string, lat, lng float64) *Validator {
	// Validate latitude (-90 to 90)
	if lat < -90 || lat > 90 {
		v.errors = append(v.errors, ValidationError{
			Field:   latField,
			Message: "latitude must be between -90 and 90",
			Value:   lat,
		})
	}
	
	// Validate longitude (-180 to 180)
	if lng < -180 || lng > 180 {
		v.errors = append(v.errors, ValidationError{
			Field:   lngField,
			Message: "longitude must be between -180 and 180",
			Value:   lng,
		})
	}
	
	return v
}

// OneOf validates that value is one of the allowed values
func (v *Validator) OneOf(field string, value interface{}, allowed []interface{}) *Validator {
	for _, allowedValue := range allowed {
		if reflect.DeepEqual(value, allowedValue) {
			return v
		}
	}
	
	v.errors = append(v.errors, ValidationError{
		Field:   field,
		Message: fmt.Sprintf("must be one of: %v", allowed),
		Value:   value,
	})
	return v
}

// Regex validates against a regular expression
func (v *Validator) Regex(field string, value string, pattern string) *Validator {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: "invalid regex pattern",
		})
		return v
	}
	
	if !regex.MatchString(value) {
		v.errors = append(v.errors, ValidationError{
			Field:   field,
			Message: fmt.Sprintf("must match pattern: %s", pattern),
			Value:   value,
		})
	}
	return v
}

// ValidateStruct validates a struct using field tags
func (v *Validator) ValidateStruct(s interface{}) error {
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("expected struct, got %s", val.Kind())
	}
	
	typ := val.Type()
	
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		fieldName := fieldType.Name
		tag := fieldType.Tag.Get("validate")
		
		if tag == "" {
			continue
		}
		
		v.validateField(fieldName, field.Interface(), tag)
	}
	
	return nil
}

// validateField validates a single field based on tags
func (v *Validator) validateField(fieldName string, value interface{}, tags string) {
	rules := strings.Split(tags, ",")
	
	for _, rule := range rules {
		parts := strings.SplitN(rule, "=", 2)
		ruleName := parts[0]
		
		switch ruleName {
		case "required":
			v.Required(fieldName, value)
		case "min":
			if len(parts) == 2 {
				if min, err := strconv.Atoi(parts[1]); err == nil {
					if str, ok := value.(string); ok {
						v.MinLength(fieldName, str, min)
					}
				}
			}
		case "max":
			if len(parts) == 2 {
				if max, err := strconv.Atoi(parts[1]); err == nil {
					if str, ok := value.(string); ok {
						v.MaxLength(fieldName, str, max)
					}
				}
			}
		case "email":
			if str, ok := value.(string); ok {
				v.Email(fieldName, str)
			}
		}
	}
}

// HasErrors returns true if there are validation errors
func (v *Validator) HasErrors() bool {
	return len(v.errors) > 0
}

// Errors returns all validation errors
func (v *Validator) Errors() ValidationErrors {
	return v.errors
}

// Error returns the error string
func (v *Validator) Error() string {
	return v.errors.Error()
}
