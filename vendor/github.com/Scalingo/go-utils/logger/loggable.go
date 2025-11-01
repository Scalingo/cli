package logger

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"
)

type Loggable interface {
	LogFields() logrus.Fields
}

// FieldsFor extracts loggable fields from a struct based on the "log" tag.
// It returns a logrus.Fields map where the keys are the tag values prefixed
// with the provided prefix, and the values are the corresponding field values.
//
// If the struct implements the Loggable interface. The `log` tags are ignored
// and the LogFields method is used to extract the fields.
//
// If the struct has no fields with the "log" tag, it checks if the struct
// implements the fmt.Stringer interface. If it does, it adds a single field
// with the prefix as the key and the result of the String() method as the value.
// If the struct does not implement fmt.Stringer, it adds a single field with
// the prefix as the key and a default error message as the value.
//
// The "omitempty" option specifies that the field should be omitted if the field has an empty value, defined as false, 0, a nil pointer, a nil interface value, and any array, slice, map, or string of length zero.
//
// The "omitzero" option specifies that the field should be omitted if the field has a zero value, according to rules:
//
// 1) If the field type has an "IsZero() bool" method, that will be used to determine whether the value is zero.
//
// 2) Otherwise, the value is zero if it is the zero value for its type.
//
// If both "omitempty" and "omitzero" are specified, the field will be omitted if the value is either empty or zero (or both).
//
// Parameters:
// - value: The struct to extract fields from.
// - prefix: The prefix to add to each field key.
//
// Returns:
// - logrus.Fields: A map of loggable fields.
func FieldsFor(prefix string, value interface{}) logrus.Fields {
	fields := logrus.Fields{}

	if loggableValue, ok := value.(Loggable); ok {
		for k, v := range loggableValue.LogFields() {
			fields[fmt.Sprintf("%s_%s", prefix, k)] = v
		}
		return fields
	}

	val := reflect.Indirect(reflect.ValueOf(value))

	if val.Kind() == reflect.Struct {
		for i := 0; i < val.NumField(); i++ {
			structField := val.Type().Field(i)
			tag, found := structField.Tag.Lookup("log")
			if !found {
				// if the `log` tag has not been found, iterate to the next structure field
				continue
			}

			fieldValue := val.Field(i)
			tagName, tagOpts := parseTag(tag)
			if tagOpts.contains("omitempty") && isEmptyValue(fieldValue) {
				// do not keep the log field if `omitempty` is set and it is an empty value
				continue
			}

			if tagOpts.contains("omitzero") {
				fieldIsZero := determineIsZeroMethod(structField.Type)

				if (fieldIsZero == nil && fieldValue.IsZero()) ||
					(fieldIsZero != nil && fieldIsZero(fieldValue)) {
					// do not keep the log field if `omitzero` is set and it is a zero value
					continue
				}

				// the field do NOT have a zero value, add it to the logger
			}

			fields[fmt.Sprintf("%s_%s", prefix, tagName)] = fieldValue.Interface()
		}
	}

	if len(fields) != 0 {
		return fields
	}

	if valueStr, ok := value.(fmt.Stringer); ok {
		fields[prefix] = valueStr.String()
	} else {
		fields[prefix] = "failed to use FieldsFor on struct: invalid type"
	}

	return fields
}

func WithStructToCtx(ctx context.Context, prefix string, value interface{}) (context.Context, logrus.FieldLogger) {
	return WithFieldsToCtx(ctx, FieldsFor(prefix, value))
}

// tagOptions is the string following a comma in a struct field's "log"
// tag, or the empty string. It does not include the leading comma.
type tagOptions string

// parseTag splits a struct field's json tag into its name and
// comma-separated options.
func parseTag(tag string) (string, tagOptions) {
	tag, opt, _ := strings.Cut(tag, ",")
	return tag, tagOptions(opt)
}

// contains reports whether a comma-separated list of options
// contains a particular substr flag. substr must be surrounded by a
// string boundary or commas.
func (o tagOptions) contains(optionName string) bool {
	if len(o) == 0 {
		return false
	}
	s := string(o)
	for s != "" {
		var name string
		name, s, _ = strings.Cut(s, ",")
		if name == optionName {
			return true
		}
	}
	return false
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Interface, reflect.Pointer:
		return v.IsZero()
	}
	return false
}

// isZeroer is the interface that wraps the basic IsZero method.
//
// IsZero returns true if the value is zero.
type isZeroer interface {
	IsZero() bool
}

// isZeroerType is the type for the isZeroer interface.
var isZeroerType = reflect.TypeFor[isZeroer]()

// isZeroMethod is a function that uses a type's IsZero method
type isZeroMethod func(v reflect.Value) bool

// determineIsZeroMethod provides a function that uses a type's IsZero method.
func determineIsZeroMethod(t reflect.Type) isZeroMethod {
	if t.Kind() == reflect.Interface && t.Implements(isZeroerType) {
		return func(v reflect.Value) bool {
			// Avoid panics calling IsZero on a nil interface or
			// non-nil interface with nil pointer.
			return v.IsNil() ||
				(v.Elem().Kind() == reflect.Pointer && v.Elem().IsNil()) ||
				//nolint:errcheck
				// We don't need to check the error for this type assertion. v is an interface which implements isZeroer.
				v.Interface().(isZeroer).IsZero()
		}
	}

	if t.Kind() == reflect.Pointer && t.Implements(isZeroerType) {
		return func(v reflect.Value) bool {
			// Avoid panics calling IsZero on nil pointer.
			return v.IsNil() ||
				//nolint:errcheck
				// We don't need to check the error for this type assertion. v is a pointer which implements isZeroer.
				v.Interface().(isZeroer).IsZero()
		}
	}

	if t.Implements(isZeroerType) {
		return func(v reflect.Value) bool {
			//nolint:errcheck
			// We don't need to check the error for this type assertion. v implements isZeroer.
			return v.Interface().(isZeroer).IsZero()
		}
	}

	if reflect.PointerTo(t).Implements(isZeroerType) {
		return func(v reflect.Value) bool {
			if !v.CanAddr() {
				// Temporarily box v so we can take the address.
				v2 := reflect.New(v.Type()).Elem()
				v2.Set(v)
				v = v2
			}
			//nolint:errcheck
			// We don't need to check the error for this type assertion. v implements isZeroer.
			return v.Addr().Interface().(isZeroer).IsZero()
		}
	}

	return nil
}
