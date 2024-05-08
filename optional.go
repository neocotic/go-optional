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
// That said; Optional is intended more for reading input rather than writing output.
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

// errNotPresent is used when panicking.
var errNotPresent = fmt.Errorf("optional value not present")

// Filter returns the Optional if it has a value present that the given function returns true for, otherwise an empty
// Optional.
//
// Warning: While fn will only be called if Optional has a value present, that value may still be nil or the zero value
// for T.
//
// For example;
//
//	isPos := func(value int) bool {
//		return value >= 0
//	}
//	Empty[int]().Filter(isPos)  // Empty[int]()
//	Of(-123).Filter(isPos)      // Empty[int]()
//	Of(0).Filter(isPos)         // Of(0)
//	Of(123).Filter(isPos)       // Of(123)
//
//	isLower := func(value string) bool {
//		return !strings.ContainsFunc(value, unicode.IsUpper)
//	}
//	Empty[string]().Filter(isLower)  // Empty[string]()
//	Of("ABC").Filter(isLower)        // Empty[string]()
//	Of("").Filter(isLower)           // Of("")
//	Of("abc").Filter(isLower)        // Of("abc")
func (o Optional[T]) Filter(fn func(value T) bool) Optional[T] {
	if o.present && fn(o.value) {
		return o
	}
	return Optional[T]{}
}

// Get returns the value of the Optional and whether it is present.
//
// For example;
//
//	Empty[int]().Get()  // 0, false
//	Of(0).Get()         // 0, true
//	Of(123).Get()       // 123, true
//
//	Empty[string]().Get()  // "", false
//	Of("").Get()           // "", true
//	Of("abc").Get()        // "abc", true
func (o Optional[T]) Get() (T, bool) {
	return o.value, o.present
}

// IfPresent calls the given function only the Optional has a value present, passing the value to the function.
//
// Warning: While fn will only be called if Optional has a value present, that value may still be nil or the zero value
// for T.
//
// For example;
//
//	intFunc := func(value int) {
//		fmt.Println(value)
//	}
//	Empty[int]().IfPresent(intFunc)  // Does nothing
//	Of(0).IfPresent(intFunc)         // Prints "0"
//	Of(123).IfPresent(intFunc)       // Prints "123"
//
//	stringFunc := func(value string) {
//		fmt.Println(value)
//	}
//	Empty[string]().IfPresent(stringFunc)  // Does nothing
//	Of("").IfPresent(stringFunc)           // Prints ""
//	Of("abc").IfPresent(stringFunc)        // Prints "abc"
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
//
// For example;
//
//	Empty[int]().IsEmpty()  // true
//	Of(0).IsEmpty()         // false
//	Of(123).IsEmpty()       // false
//
//	Empty[string]().IsEmpty()  // true
//	Of("").IsEmpty()           // false
//	Of("abc").IsEmpty()        // false
func (o Optional[T]) IsEmpty() bool {
	return !o.present
}

// IsPresent returns whether the value of the Optional is present. That is; it has been explicitly set.
//
// For example;
//
//	Empty[int]().IsPresent()  // false
//	Of(0).IsPresent()         // true
//	Of(123).IsPresent()       // true
//
//	Empty[string]().IsPresent()  // false
//	Of("").IsPresent()           // true
//	Of("abc").IsPresent()        // true
func (o Optional[T]) IsPresent() bool {
	return o.present
}

// IsZero returns whether the value of the Optional is absent. That is; it has NOT been explicitly set.
//
// IsZero is effectively the inverse of IsPresent and an alternative for IsEmpty that conforms to the yaml.IsZeroer
// interface. It's important to note that IsZero will not return true if the underlying value of the Optional is equal
// to the zero value for T but instead only if the value is absent.
//
// For example;
//
//	Empty[int]().IsZero()  // true
//	Of(0).IsZero()         // false
//	Of(123).IsZero()       // false
//
//	Empty[string]().IsZero()  // true
//	Of("").IsZero()           // false
//	Of("abc").IsZero()        // false
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
//
// For example;
//
//	defaultInt := -1
//	Empty[int]().OrElse(defaultInt)  // -1
//	Of(0).OrElse(defaultInt)         // 0
//	Of(123).OrElse(defaultInt)       // 123
//
//	defaultString := "unknown"
//	Empty[string]().OrElse(defaultString)  // "unknown"
//	Of("").OrElse(defaultString)           // ""
//	Of("abc").OrElse(defaultString)        // "abc"
func (o Optional[T]) OrElse(other T) T {
	if o.present {
		return o.value
	}
	return other
}

// OrElseGet returns the value of the Optional if present, otherwise calls other and returns its return value. This is
// recommended over OrElse in cases where a default value is expensive to initialize so lazy-initializes it.
//
// For example;
//
//	defaultInt := func() int {
//		return -1
//	}
//	Empty[int]().OrElseGet(defaultFunc)  // -1
//	Of(0).OrElseGet(defaultFunc)         // 0
//	Of(123).OrElseGet(defaultFunc)       // 123
//
//	defaultString := func() string {
//		return "unknown"
//	}
//	Empty[string]().OrElseGet(defaultString)  // "unknown"
//	Of("").OrElseGet(defaultString)           // ""
//	Of("abc").OrElseGet(defaultString)        // "abc"
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
//
// For example;
//
//	defaultInt := func() (int, error) {
//		return -1, nil
//	}
//	Empty[int]().OrElseTryGet(defaultFunc)  // -1, nil
//	Of(0).OrElseTryGet(defaultFunc)         // 0, nil
//	Of(123).OrElseTryGet(defaultFunc)       // 123, nil
//
//	var defaultStringUsed bool
//	errDefaultStringUsed := errors.New("default string already used")
//	defaultString := func() (string, error) {
//		if defaultStringUsed {
//			return "", errDefaultStringUsed
//		}
//		defaultStringUsed = true
//		return "unknown", nil
//	}
//	Empty[string]().OrElseTryGet(defaultString)  // "unknown", nil
//	Of("").OrElseTryGet(defaultString)           // "", nil
//	Of("abc").OrElseTryGet(defaultString)        // "abc", nil
//	Empty[string]().OrElseTryGet(defaultString)  // "", errDefaultStringUsed
func (o Optional[T]) OrElseTryGet(other func() (T, error)) (T, error) {
	if o.present {
		return o.value, nil
	}
	return other()
}

// Require returns the value of the Optional only if present, otherwise panics.
//
// For example;
//
//	Empty[int]().Require()  // panics
//	Of(0).Require()         // 0
//	Of(123).Require()       // 123
//
//	Empty[string]().Require()  // panics
//	Of("").Require()           // ""
//	Of("abc").Require()        // "abc"
func (o Optional[T]) Require() T {
	if o.present {
		return o.value
	}
	panic(errNotPresent)
}

// String returns a string representation of the underlying value, if any.
//
// For example;
//
//	Empty[int]().String()  // "0"
//	Of(0).String()         // "0"
//	Of(123).String()       // "123"
//
//	Empty[string]().String()  // ""
//	Of("").String()           // ""
//	Of("abc").String()        // "abc"
func (o Optional[T]) String() string {
	return fmt.Sprint(o.value)
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
//
// For example;
//
//	Compare(Empty[int](), Of(0))  // -1
//	Compare(Of(0), Of(123))       // -1
//
//	Compare(Empty[int](), Empty[int]())  // 0
//	Compare(Of(0), Of(0))                // 0
//	Compare(Of(123), Of(123))            // 0
//
//	Compare(Of(0), Empty[int]())  // 1
//	Compare(Of(123), Of(0))       // 1
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
//
// For example;
//
//	Empty[int]().Get()  // 0, false
//
//	Empty[string]().Get()  // "", false
func Empty[T any]() Optional[T] {
	return Optional[T]{}
}

// Find returns the first given Optional that has a value present, otherwise an empty Optional.
//
// For example;
//
//	Find[int]()                         // Empty[int]()
//	Find(Empty[int]())                  // Empty[int]()
//	Find(Empty[int](), Of(0), Of(123))  // Of(0)
//
//	Find[string]()                            // Empty[string]()
//	Find(Empty[string]())                     // Empty[string]()
//	Find(Empty[string](), Of("abc"), Of(""))  // Of("abc")
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
//
// For example;
//
//	toString := func(value int) Optional[string] {
//		if value == 0 {
//			return Empty[string]()
//		}
//		return Of(strconv.FormatInt(int64(value), 10))
//	}
//	FlatMap(Empty[int](), toString)  // Empty[string]()
//	FlatMap(Of(0), toString)         // Empty[string]()
//	FlatMap(Of(123), toString)       // Of("123")
//
//	toInt := func(value string) Optional[int] {
//		if value == "" {
//			return Empty[int]()
//		}
//		i, err := strconv.ParseInt(value, 10, 0)
//		if err != nil {
//			panic(err)
//		}
//		return OfZeroable(int(i))
//	}
//	FlatMap(Empty[string](), toInt)  // Empty[int]()
//	FlatMap(Of(""), toInt)           // Empty[int]()
//	FlatMap(Of("0"), toInt)          // Empty[int]()
//	FlatMap(Of("123"), toInt)        // Of(123)
func FlatMap[T, M any](opt Optional[T], fn func(value T) Optional[M]) Optional[M] {
	if !opt.present {
		return Optional[M]{}
	}
	return fn(opt.value)
}

// GetAny returns a slice containing only the values of any given Optional that has a value present, where possible.
//
// For example;
//
//	GetAny[int]()                         // nil
//	GetAny(Empty[int]())                  // nil
//	GetAny(Empty[int](), Of(0), Of(123))  // []int{0, 123}
//
//	GetAny[string]()                            // nil
//	GetAny(Empty[string]())                     // nil
//	GetAny(Empty[string](), Of("abc"), Of(""))  // []string{"abc", ""}
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
//
// For example;
//
//	toString := func(value int) string {
//		return strconv.FormatInt(int64(value), 10)
//	}
//	Map(Empty[int](), toString)  // Empty[string]()
//	Map(Of(0), toString)         // Of("0")
//	Map(Of(123), toString)       // Of("123")
//
//	toInt := func(value string) int {
//		i, err := strconv.ParseInt(value, 10, 0)
//		if err != nil {
//			panic(err)
//		}
//		return int(i)
//	}
//	Map(Empty[string](), toInt)  // Empty[int]()
//	Map(Of("0"), toInt)          // Of(0)
//	Map(Of("123"), toInt)        // Of(123)
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
//
// For example;
//
//	MustFind[int]()                         // panics
//	MustFind(Empty[int]())                  // panics
//	MustFind(Empty[int](), Of(0), Of(123))  // 0
//
//	MustFind[string]()                            // panics
//	MustFind(Empty[string]())                     // panics
//	MustFind(Empty[string](), Of("abc"), Of(""))  // "abc"
func MustFind[T any](opts ...Optional[T]) T {
	for _, opt := range opts {
		if opt.present {
			return opt.value
		}
	}
	panic(errNotPresent)
}

// Of returns an Optional with the given value present.
//
// For example;
//
//	Of((*int)(nil)).Get()     // nil, true
//	Of(ptrs.ZeroInt()).Get()  // &0, true
//	Of(ptrs.Int(123)).Get()   // &123, true
//	Of(0).Get()               // 0, true
//	Of(123).Get()             // 123, true
//
//	Of((*string)(nil)).Get()      // nil, true
//	Of(ptrs.ZeroString()).Get()   // &"", true
//	Of(ptrs.String("abc")).Get()  // &"abc", true
//	Of("").Get()                  // "", true
//	Of("abc").Get()               // "abc", true
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
//
// For example;
//
//	OfNillable((*int)(nil)).Get()     // nil, false
//	OfNillable(ptrs.ZeroInt()).Get()  // &0, true
//	OfNillable(ptrs.Int(123)).Get()   // &123, true
//	OfNillable(0).Get()               // 0, true
//	OfNillable(123).Get()             // 123, true
//
//	OfNillable((*string)(nil)).Get()      // nil, false
//	OfNillable(ptrs.ZeroString()).Get()   // &"", true
//	OfNillable(ptrs.String("abc")).Get()  // &"abc", true
//	OfNillable("").Get()                  // "", true
//	OfNillable("abc").Get()               // "abc", true
func OfNillable[T any](value T) Optional[T] {
	if isNil(value) {
		return Optional[T]{}
	}
	return Optional[T]{
		present: true,
		value:   value,
	}
}

// OfPointer returns an Optional with the given value present as a pointer.
//
// For example;
//
//	OfPointer(0).Get()    // &0, true
//	OfPointer(123).Get()  // &123, true
//
//	OfPointer("").Get()     // &"", true
//	OfPointer("abc").Get()  // &"abc", true
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
//
// For example;
//
//	OfZeroable((*int)(nil)).Get()     // nil, false
//	OfZeroable(ptrs.ZeroInt()).Get()  // &0, true
//	OfZeroable(ptrs.Int(123)).Get()   // &123, true
//	OfZeroable(0).Get()               // 0, false
//	OfZeroable(123).Get()             // 123, true
//
//	OfZeroable((*string)(nil)).Get()      // nil, false
//	OfZeroable(ptrs.ZeroString()).Get()   // &"", true
//	OfZeroable(ptrs.String("abc")).Get()  // &"abc", true
//	OfZeroable("").Get()                  // "", false
//	OfZeroable("abc").Get()               // "abc", true
func OfZeroable[T any](value T) Optional[T] {
	if isZero(value) {
		return Optional[T]{}
	}
	return Optional[T]{
		present: true,
		value:   value,
	}
}

// RequireAny returns a slice containing only the values of any given Optional that has a value present, panicking only
// if no Optional could be found with a value present.
//
// For example;
//
//	RequireAny[int]()                         // panics
//	RequireAny(Empty[int]())                  // panics
//	RequireAny(Empty[int](), Of(0), Of(123))  // []int{0, 123}
//
//	RequireAny[string]()                            // panics
//	RequireAny(Empty[string]())                     // panics
//	RequireAny(Empty[string](), Of("abc"), Of(""))  // []string{"abc", ""}
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
//
// For example;
//
//	toString := func(value int) (Optional[string], error) {
//		if value == 0 {
//			return Empty[string](), nil
//		}
//		return Of(strconv.FormatInt(int64(value), 10)), nil
//	}
//	TryFlatMap(Empty[int](), toString)  // Empty[string](), nil
//	TryFlatMap(Of(0), toString)         // Empty[string](), nil
//	TryFlatMap(Of(123), toString)       // Of("123"), nil
//
//	toInt := func(value string) (Optional[int], error) {
//		if value == "" {
//			return Empty[int](), nil
//		}
//		i, err := strconv.ParseInt(value, 10, 0)
//		if err != nil {
//			return Empty[int](), err
//		}
//		return OfZeroable(int(i)), nil
//	}
//	TryFlatMap(Empty[string](), toInt)  // Empty[int](), nil
//	TryFlatMap(Of(""), toInt)           // Empty[int](), nil
//	TryFlatMap(Of("0"), toInt)          // Empty[int](), nil
//	TryFlatMap(Of("123"), toInt)        // Of(123), nil
//	TryFlatMap(Of("abc"), toInt)        // Empty[int](), strconv.NumError
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
//
// For example;
//
//	toString := func(value int) (string, error) {
//		return strconv.FormatInt(int64(value), 10), nil
//	}
//	TryMap(Empty[int](), toString)  // Empty[string]()
//	TryMap(Of(0), toString)         // Of("0")
//	TryMap(Of(123), toString)       // Of("123")
//
//	toInt := func(value string) (int, error) {
//		i, err := strconv.ParseInt(value, 10, 0)
//		return int(i), err
//	}
//	TryMap(Empty[string](), toInt)  // Empty[int]()
//	TryMap(Of("0"), toInt)          // Of(0)
//	TryMap(Of("123"), toInt)        // Of(123)
//	TryMap(Of("abc"), toInt)        // Empty[int](), strconv.NumError
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
	switch rv.Kind() {
	case reflect.Invalid:
		return true
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return rv.IsNil()
	default:
		return false
	}
}

// isZero returns whether the given value is zero for its type using reflection.
func isZero(value any) bool {
	rv := reflect.ValueOf(value)
	return !rv.IsValid() || rv.IsZero()
}
