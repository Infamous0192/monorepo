package validation

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"text/template"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// Validation handles request validation
type Validation struct {
	validate *validator.Validate
}

// ValidationError represents validation errors
type ValidationError struct {
	Errors map[string]string `json:"errors"`
}

// Error implements the error interface
func (e ValidationError) Error() string {
	return "Validation failed"
}

// validate is a singleton instance of the validator
var validate = validator.New()

// Validate validates a struct using go-playground/validator
func Validate(s interface{}) error {
	err := validate.Struct(s)
	if err != nil {
		errors := make(map[string]string)

		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()
			tag := err.Tag()
			errors[field] = getErrorMessage(tag)
		}

		return ValidationError{
			Errors: errors,
		}
	}

	return nil
}

// getErrorMessage returns a human-readable error message for a validation tag
func getErrorMessage(tag string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Must be a valid email address"
	case "min":
		return "Value is too short"
	case "max":
		return "Value is too long"
	default:
		return "Invalid value"
	}
}

// NewValidation creates a new validation instance
func NewValidation() *Validation {
	return &Validation{
		validate: validator.New(),
	}
}

// Body parses the request body into a given struct and validates it.
func (v *Validation) Body(data interface{}, ctx *fiber.Ctx) error {
	if err := ctx.BodyParser(data); err != nil {
		return ValidationError{
			Errors: map[string]string{
				"body": err.Error(),
			},
		}
	}

	return v.ValidateStruct(data)
}

// Query parses the query parameters into a given struct and validates it.
func (v *Validation) Query(data interface{}, ctx *fiber.Ctx) error {
	if err := ctx.QueryParser(data); err != nil {
		return ValidationError{
			Errors: map[string]string{
				"query": err.Error(),
			},
		}
	}

	return v.ValidateStruct(data)
}

// Retrieves an integer value from the Fiber context using the specified key.
func (v *Validation) ParamsInt(ctx *fiber.Ctx, keys ...string) (int, error) {
	key := "id"
	if len(keys) > 0 && keys[0] != "" {
		key = keys[0]
	}

	value, err := ctx.ParamsInt(key)
	if err != nil {
		return -1, ValidationError{
			Errors: map[string]string{
				"id": "ID not valid",
			},
		}
	}

	return value, nil
}

func (v *Validation) Field(field interface{}, tag string) error {
	if err := v.validate.Var(field, tag); err != nil {
		return fmt.Errorf("%v", getMessage(err.(validator.FieldError)))
	}

	return nil
}

// ValidateStruct validates a given struct and returns a map of validation errors.
func (v *Validation) ValidateStruct(values interface{}) error {
	if err := v.validate.Struct(values); err != nil {
		st := reflect.Indirect(reflect.ValueOf(values)).Type() // Get indirect type
		messages := make(map[string]string)

		for _, err := range err.(validator.ValidationErrors) {
			field, _ := st.FieldByName(err.Field())

			key := field.Tag.Get("json")

			if key == "" {
				key = strings.ToLower(err.Field())
			}

			messages[key] = getMessage(err)
		}

		return ValidationError{
			Errors: messages,
		}
	}

	return nil
}

// getMessage returns a validation error message based on the given validator.FieldError.
func getMessage(err validator.FieldError) string {
	attribute, format := attributes[err.Field()], messages[err.Tag()]
	if format == "" {
		format = messages["default"]
	}

	if attribute == "" {
		attribute = err.Field()
	}

	param := strings.Join(strings.Split(err.Param(), " "), ", ")

	t := template.Must(template.New("").Parse(format))
	b := new(strings.Builder)
	if err := t.Execute(b, map[string]interface{}{
		"Attribute": attribute,
		"Param":     param,
	}); err != nil {
		return err.Error()
	}

	return b.String()
}

// FormValue validates and binds form field value to struct
func (v *Validation) FormValue(obj interface{}, field string, c *fiber.Ctx) error {
	// Get form field value
	value := c.FormValue(field)
	if value == "" {
		return fmt.Errorf("form field '%s' is required", field)
	}

	// Parse JSON value
	if err := json.Unmarshal([]byte(value), obj); err != nil {
		return fmt.Errorf("invalid JSON in form field '%s': %v", field, err)
	}

	// Validate struct
	if err := v.validate.Struct(obj); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errorMap := make(map[string]string)
			for _, e := range validationErrors {
				errorMap[e.Field()] = fmt.Sprintf("validation failed on '%s' tag", e.Tag())
			}
			return ValidationError{Errors: errorMap}
		}
		return err
	}

	return nil
}
