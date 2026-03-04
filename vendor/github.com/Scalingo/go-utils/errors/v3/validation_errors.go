package errors

import (
	"fmt"
	"strings"
)

// ValidationErrors store each errors associated to every fields of a model
type ValidationErrors struct {
	Errors map[string][]string `json:"errors"`
}

func (v *ValidationErrors) Error() string {
	var builder strings.Builder
	index := 0

	for field, errs := range v.Errors {
		index++
		builder.WriteString(fmt.Sprintf("%s=%s", field, strings.Join(errs, ", ")))
		if index < len(v.Errors) {
			builder.WriteString(" ")
		}
	}

	return builder.String()
}

// ValidationErrorsBuilder is used to provide a simple way to create a ValidationErrors struct. The typical usecase is:
//
//	func (m *MyModel) Validate(ctx context.Context) *ValidationErrors {
//		validations := document.NewValidationErrorsBuilder()
//
//		if m.Name == "" {
//			validations.Set("name", "should not be empty")
//		}
//
//		if m.Email == "" {
//			validations.Set("email", "should not be empty")
//		}
//
//		return validations.Build()
//	}
type ValidationErrorsBuilder struct {
	errors map[string][]string
}

// NewValidationErrors return an empty ValidationErrors struct
func NewValidationErrorsBuilder() *ValidationErrorsBuilder {
	return &ValidationErrorsBuilder{
		errors: make(map[string][]string),
	}
}

// Set will add an error on a specific field, if the field already contains an error, it will just add it to the current errors list
func (v *ValidationErrorsBuilder) Set(field, err string) *ValidationErrorsBuilder {
	v.errors[field] = append(v.errors[field], err)
	return v
}

// Get will return all errors set for a specific field
func (v *ValidationErrorsBuilder) Get(field string) []string {
	return v.errors[field]
}

// Merge ValidationErrors with another ValidationErrors
func (v *ValidationErrorsBuilder) Merge(verr *ValidationErrors) *ValidationErrorsBuilder {
	return v.MergeWithPrefix("", verr)
}

// MergeWithPrefix is merging ValidationErrors in another ValidationError
// adding a prefix for each error field
func (v *ValidationErrorsBuilder) MergeWithPrefix(prefix string, verr *ValidationErrors) *ValidationErrorsBuilder {
	if verr == nil {
		return v
	}
	if prefix != "" && prefix[len(prefix)-1] != '_' {
		prefix = prefix + "_"
	}

	for key, values := range verr.Errors {
		for _, value := range values {
			v.Set(prefix+key, value)
		}
	}
	return v
}

// Build will send a ValidationErrors struct if there is some errors or nil if no errors has been defined
func (v *ValidationErrorsBuilder) Build() error {
	if len(v.errors) == 0 {
		return nil
	}

	return &ValidationErrors{
		Errors: v.errors,
	}
}
