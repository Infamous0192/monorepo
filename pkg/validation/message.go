package validation

// messages is a map of validation error messages.
// The keys represent different types of validation errors,
// and the values are the corresponding error messages in English.
var messages = map[string]string{
	"default":   "{{.Attribute}} is not valid",
	"required":  "{{.Attribute}} is required",
	"min":       "{{.Attribute}} must be greater than {{.Param}}",
	"max":       "{{.Attribute}} must be less than {{.Param}}",
	"oneof":     "{{.Attribute}} must be one of {{.Param}}",
	"exist":     "{{.Attribute}} not found",
	"not_exist": "{{.Attribute}} has already been used",
}

// attributes is a map of attribute names used in validation.
var attributes = map[string]string{
	"Note": "Note",
}
