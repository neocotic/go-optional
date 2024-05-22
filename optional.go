// Copyright (C) 2024 neocotic
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package optional enables the ability to differentiate a value that has its zero value due to not being set from
// having a zero value that was explicitly set.
package optional

import (
	"cmp"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"gopkg.in/yaml.v3"
	"reflect"
	"strings"
)

// Optional contains an immutable value as well as an indication whether it was explicitly set. This can be especially
// useful when needing to differentiate the source of a zero value.
//
// For the best experience when marshaling a struct with Optional struct field types, the following information may be
// useful;
//
//   - json: it's recommended to include the "omitempty" tag option and have the Optional field type declared as a
//     pointer, otherwise the "omitempty" tag option is ignored
//   - xml: seems to work perfectly as expected
//   - yaml: it's recommended to include the "omitempty" tag option
//
// That said; Optional is intended more for reading input rather than writing output. An important note for
// unmarshalling is that yaml, unlike json, will skip an Optional struct field that has been given an explicit null
// value, resulting in an empty Optional.
type Optional[T any] struct {
	// present is whether value was explicitly set.
	present bool
	// value is the value.
	value T
}

var (
	_ fmt.Stringer     = (*Optional[any])(nil)
	_ json.Marshaler   = (*Optional[any])(nil)
	_ json.Unmarshaler = (*Optional[any])(nil)
	_ xml.Marshaler    = (*Optional[any])(nil)
	_ xml.Unmarshaler  = (*Optional[any])(nil)
	_ yaml.IsZeroer    = (*Optional[any])(nil)
	_ yaml.Marshaler   = (*Optional[any])(nil)
	_ yaml.Unmarshaler = (*Optional[any])(nil)
)

// emptyString is returned by Optional.String when no value is present.
const emptyString = "<empty>"

// errNotPresent is used when panicking.
var errNotPresent = fmt.Errorf("go-optional: value not present")

// Filter returns the Optional if it has a value present that the given function returns true for, otherwise an empty
// Optional.
//
// Warning: While fn will only be called if Optional has a value present, that value may still be nil or the zero value
// for T.
func (o Optional[T]) Filter(fn func(value T) bool) Optional[T] {
	if o.present && fn(o.value) {
		return o
	}
	return Optional[T]{}
}

// Get returns the value of the Optional and whether it is present.
func (o Optional[T]) Get() (T, bool) {
	return o.value, o.present
}

// IfPresent calls the given function only the Optional has a value present, passing the value to the function.
//
// Warning: While fn will only be called if Optional has a value present, that value may still be nil or the zero value
// for T.
func (o Optional[T]) IfPresent(fn func(value T)) {
	if o.present {
		fn(o.value)
	}
}

// IsEmpty returns whether the value of the Optional is absent. That is; it has NOT been explicitly set.
//
// IsEmpty is effectively the inverse of IsPresent. It's important to note that IsEmpty will not return true if the
// underlying value of the Optional is equal to the zero value for T and in no way checks the length of the underlying
// value but instead only if the value is absent.
func (o Optional[T]) IsEmpty() bool {
	return !o.present
}

// IsPresent returns whether the value of the Optional is present. That is; it has been explicitly set.
func (o Optional[T]) IsPresent() bool {
	return o.present
}

// IsZero returns whether the value of the Optional is absent. That is; it has NOT been explicitly set.
//
// IsZero is effectively the inverse of IsPresent and an alternative for IsEmpty that conforms to the yaml.IsZeroer
// interface. It's important to note that IsZero will not return true if the underlying value of the Optional is equal
// to the zero value for T but instead only if the value is absent.
func (o Optional[T]) IsZero() bool {
	return !o.present
}

// MarshalJSON marshals the value of the Optional into JSON, if present, otherwise returns a null-like value.
//
// An error is returned if unable to marshal the value.
func (o Optional[T]) MarshalJSON() ([]byte, error) {
	if !o.present {
		return []byte("null"), nil
	}
	return json.Marshal(o.value)
}

// MarshalXML marshals the encoded value of the Optional into XML, if present, otherwise nothing is written to the given
// encoder.
//
// An error is returned if unable to write the value to the given encoder.
func (o Optional[T]) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	// In cases where an Optional is passed directly to xml.Marshal the start element should be ignored
	if start.Name.Space == "" && strings.HasPrefix(start.Name.Local, "Optional") {
		if !o.present {
			return e.Encode(nil)
		}
		return e.Encode(o.value)
	}
	if !o.present {
		return e.EncodeElement(nil, start)
	}
	return e.EncodeElement(o.value, start)
}

// MarshalYAML marshals the value of the Optional into YAML, if present, otherwise returns a null-like value.
//
// An error is returned if unable to marshal the value.
func (o Optional[T]) MarshalYAML() (any, error) {
	if !o.present {
		return nil, nil
	}
	return o.value, nil
}

// OrElse returns the value of the Optional if present, otherwise other.
func (o Optional[T]) OrElse(other T) T {
	if o.present {
		return o.value
	}
	return other
}

// OrElseGet returns the value of the Optional if present, otherwise calls other and returns its return value. This is
// recommended over OrElse in cases where a default value is expensive to initialize so lazy-initializes it.
func (o Optional[T]) OrElseGet(other func() T) T {
	if o.present {
		return o.value
	}
	return other()
}

// OrElseTryGet returns the value of the Optional if present, otherwise calls other and returns its return value. This
// is recommended over OrElse in cases where a default value is expensive to initialize so lazy-initializes it. The
// difference from OrElseGet is that the given function may return an error which, if not nil, will be returned by
// OrElseTryGet.
func (o Optional[T]) OrElseTryGet(other func() (T, error)) (T, error) {
	if o.present {
		return o.value, nil
	}
	return other()
}

// Require returns the value of the Optional only if present, otherwise panics.
func (o Optional[T]) Require() T {
	if o.present {
		return o.value
	}
	panic(errNotPresent)
}

// String returns a string representation of the underlying value, if any.
func (o Optional[T]) String() string {
	if o.present {
		return fmt.Sprint(o.value)
	}
	return emptyString
}

// UnmarshalJSON unmarshalls the JSON data provided as the value for the Optional. Anytime UnmarshalJSON is called, it
// treats the Optional as having a value even though that value may still be nil or the zero value for T.
//
// An error is returned if unable to unmarshal data.
func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &o.value); err != nil {
		return err
	}
	o.present = true
	return nil
}

// UnmarshalXML unmarshalls the decoded XML element provided as the value for the Optional. Anytime UnmarshalXML is
// called, it treats the Optional as having a value even though that value may still be nil or the zero value for T.
//
// An error is returned if unable to unmarshal the given element.
func (o *Optional[T]) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if err := d.DecodeElement(&o.value, &start); err != nil {
		return err
	}
	o.present = true
	return nil
}

// UnmarshalYAML unmarshalls the decoded YAML node provided as the value for the Optional. Anytime UnmarshalYAML is
// called, it treats the Optional as having a value even though that value may still be nil or the zero value for T.
// However, unlike UnmarshalJSON and UnmarshalXML, the YAML unmarshaler will not call UnmarshalYAML for an empty or
// null-like value.
//
// An error is returned if unable to unmarshal the given node.
func (o *Optional[T]) UnmarshalYAML(value *yaml.Node) error {
	if err := value.Decode(&o.value); err != nil {
		return err
	}
	o.present = true
	return nil
}

// Compare returns the following:
//
//   - -1 if x has not value present and y does; or if both have a value present and the value of x is less than that of
//     y
//   - 0 if neither x nor y have a value present; or if both have a value present that are equal
//   - +1 if x has a value present and y does not; or if both have a value present and the value of x is greater than
//     that of y
//
// For floating-point types, a NaN is considered less than any non-NaN, a NaN is considered equal to a NaN, and -0.0 is
// equal to 0.0.
func Compare[T cmp.Ordered](x, y Optional[T]) int {
	switch {
	case x.present && y.present:
		return cmp.Compare(x.value, y.value)
	case x.present:
		return 1
	case y.present:
		return -1
	default:
		return 0
	}
}

// Empty returns an Optional with no value. It's the equivalent of using a zero value Optional.
func Empty[T any]() Optional[T] {
	return Optional[T]{}
}

// Find returns the first given Optional that has a value present, otherwise an empty Optional.
func Find[T any](opts ...Optional[T]) Optional[T] {
	for _, opt := range opts {
		if opt.present {
			return opt
		}
	}
	return Optional[T]{}
}

// FlatMap calls the given function and returns the Optional returned by it if the Optional provided has a value
// present, otherwise an empty Optional is returned.
//
// Warning: While fn will only be called if opt has a value present, that value may still be nil or the zero value for
// T.
func FlatMap[T, M any](opt Optional[T], fn func(value T) Optional[M]) Optional[M] {
	if !opt.present {
		return Optional[M]{}
	}
	return fn(opt.value)
}

// GetAny returns a slice containing only the values of any given Optional that has a value present, where possible.
func GetAny[T any](opts ...Optional[T]) []T {
	var filtered []T
	for _, opt := range opts {
		if opt.present {
			filtered = append(filtered, opt.value)
		}
	}
	return filtered
}

// Map returns an Optional whose value is mapped from the Optional provided using the given function, if present,
// otherwise an empty Optional.
//
// Warning: While fn will only be called if opt has a value present, that value may still be nil or the zero value for
// T.
func Map[T, M any](opt Optional[T], fn func(value T) M) Optional[M] {
	if !opt.present {
		return Optional[M]{}
	}
	return Optional[M]{
		present: true,
		value:   fn(opt.value),
	}
}

// MustFind returns the value of the first given Optional that has a value present, otherwise panics.
func MustFind[T any](opts ...Optional[T]) T {
	for _, opt := range opts {
		if opt.present {
			return opt.value
		}
	}
	panic(errNotPresent)
}

// Of returns an Optional with the given value present.
func Of[T any](value T) Optional[T] {
	return Optional[T]{
		present: true,
		value:   value,
	}
}

// OfNillable returns an Optional with the given value present only if value is nil. That is; unlike Of, OfNillable
// treats a nil value as absent and so the returned Optional will be empty.
//
// Since T can be any type, whether value is nil is checked reflectively.
func OfNillable[T any](value T) Optional[T] {
	if isNil(reflect.ValueOf(value)) {
		return Optional[T]{}
	}
	return Optional[T]{
		present: true,
		value:   value,
	}
}

// OfPointer returns an Optional with the given value present as a pointer.
func OfPointer[T any](value T) Optional[*T] {
	return Optional[*T]{
		present: true,
		value:   &value,
	}
}

// OfZeroable returns an Optional with the given value present only if value does not equal the zero value for T. That
// is; unlike Of, OfZeroable treats a value of zero as absent and so the returned Optional will be empty.
//
// Since T can be any type, whether value is equal to the zero value of T is checked reflectively.
func OfZeroable[T any](value T) Optional[T] {
	if isZero(reflect.ValueOf(value)) {
		return Optional[T]{}
	}
	return Optional[T]{
		present: true,
		value:   value,
	}
}

// RequireAny returns a slice containing only the values of any given Optional that has a value present, panicking only
// if no Optional could be found with a value present.
func RequireAny[T any](opts ...Optional[T]) []T {
	var filtered []T
	for _, opt := range opts {
		if opt.present {
			filtered = append(filtered, opt.value)
		}
	}
	if len(filtered) == 0 {
		panic(errNotPresent)
	}
	return filtered
}

// TryFlatMap calls the given function and returns the Optional returned by it if the Optional provided has a value
// present, otherwise an empty Optional is returned. The difference from FlatMap is that the given function may return
// an error which, if not nil, will be returned by TryFlatMap.
//
// Warning: While fn will only be called if opt has a value present, that value may still be nil or the zero value for
// T.
func TryFlatMap[T, M any](opt Optional[T], fn func(value T) (Optional[M], error)) (Optional[M], error) {
	if !opt.present {
		return Optional[M]{}, nil
	}
	return fn(opt.value)
}

// TryMap returns an Optional whose value is mapped from the Optional provided using the given function, if present,
// otherwise an empty Optional. The difference from Map is that the given function may return an error which, if not
// nil, will be returned by TryMap.
//
// Warning: While fn will only be called if opt has a value present, that value may still be nil or the zero value for
// T.
func TryMap[T, M any](opt Optional[T], fn func(value T) (M, error)) (Optional[M], error) {
	if !opt.present {
		return Optional[M]{}, nil
	}
	mapped, err := fn(opt.value)
	if err != nil {
		return Optional[M]{}, err
	}
	return Optional[M]{
		present: true,
		value:   mapped,
	}, nil
}

// isNil returns whether the given value is nil using reflection.
func isNil(value any) bool {
	rv := reflect.ValueOf(value)
// isNil returns whether the given reflect.Value is nil using reflection.
func isNil(rv reflect.Value) bool {
	switch rv.Kind() {
	case reflect.Invalid:
		return true
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return rv.IsNil()
	default:
		return false
	}
}

// isZero returns whether the given reflect.Value is zero for its type using reflection.
func isZero(rv reflect.Value) bool {
	return !rv.IsValid() || rv.IsZero()
}
