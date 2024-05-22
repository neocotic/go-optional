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

package optional

import (
	"cmp"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/neocotic/go-optional/internal/example"
	"github.com/neocotic/go-optional/internal/test"
	ptrs "github.com/neocotic/go-pointers"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"strconv"
	"strings"
	"testing"
	"unicode"
)

func BenchmarkOptional_Filter(b *testing.B) {
	isPos := func(value int) bool {
		return value >= 0
	}
	opt := Of(123)
	for i := 0; i < b.N; i++ {
		_ = opt.Filter(isPos)
	}
}

func ExampleOptional_Filter_int() {
	isPos := func(value int) bool {
		return value >= 0
	}

	example.Print(Empty[int]().Filter(isPos))
	example.Print(Of(-123).Filter(isPos))
	example.Print(Of(0).Filter(isPos))
	example.Print(Of(123).Filter(isPos))

	// Output:
	// <empty>
	// <empty>
	// 0
	// 123
}

func ExampleOptional_Filter_string() {
	isLower := func(value string) bool {
		return !strings.ContainsFunc(value, unicode.IsUpper)
	}

	example.Print(Empty[string]().Filter(isLower))
	example.Print(Of("ABC").Filter(isLower))
	example.Print(Of("").Filter(isLower))
	example.Print(Of("abc").Filter(isLower))

	// Output:
	// <empty>
	// <empty>
	// ""
	// "abc"
}

type optionalFilterTC[T any] struct {
	opt    Optional[T]
	fn     func(value T) bool
	expect Optional[T]
	test.Control
}

func (tc optionalFilterTC[T]) Test(t *testing.T) {
	actual := tc.opt.Filter(tc.fn)
	assert.Equal(t, tc.expect, actual, "unexpected optional")
}

func TestOptional_Filter(t *testing.T) {
	isPos := func(value int) bool {
		return value >= 0
	}
	isLower := func(value string) bool {
		return !strings.ContainsFunc(value, unicode.IsUpper)
	}

	test.RunCases(t, test.Cases{
		// Test cases for documented examples
		"on empty int Optional": optionalFilterTC[int]{
			opt:    Empty[int](),
			fn:     isPos,
			expect: Empty[int](),
		},
		"on non-empty int Optional with non-zero non-matching value": optionalFilterTC[int]{
			opt:    Of(-123),
			fn:     isPos,
			expect: Empty[int](),
		},
		"on non-empty int Optional with zero matching value": optionalFilterTC[int]{
			opt:    Of(0),
			fn:     isPos,
			expect: Of(0),
		},
		"on non-empty int Optional with non-zero matching value": optionalFilterTC[int]{
			opt:    Of(123),
			fn:     isPos,
			expect: Of(123),
		},
		"on empty string Optional": optionalFilterTC[string]{
			opt:    Empty[string](),
			fn:     isLower,
			expect: Empty[string](),
		},
		"on non-empty string Optional with non-zero non-matching value": optionalFilterTC[string]{
			opt:    Of("ABC"),
			fn:     isLower,
			expect: Empty[string](),
		},
		"on non-empty string Optional with zero value": optionalFilterTC[string]{
			opt:    Of(""),
			fn:     isLower,
			expect: Of(""),
		},
		"on non-empty string Optional with non-zero value": optionalFilterTC[string]{
			opt:    Of("abc"),
			fn:     isLower,
			expect: Of("abc"),
		},
		// Other test cases...
	})
}

func BenchmarkOptional_Get(b *testing.B) {
	opt := Of(123)
	for i := 0; i < b.N; i++ {
		_, _ = opt.Get()
	}
}

func ExampleOptional_Get_int() {
	example.PrintGet(Empty[int]().Get())
	example.PrintGet(Of(0).Get())
	example.PrintGet(Of(123).Get())

	// Output:
	// 0 false
	// 0 true
	// 123 true
}

func ExampleOptional_Get_string() {
	example.PrintGet(Empty[string]().Get())
	example.PrintGet(Of("").Get())
	example.PrintGet(Of("abc").Get())

	// Output:
	// "" false
	// "" true
	// "abc" true
}

type optionalGetTC[T any] struct {
	opt           Optional[T]
	expectPresent bool
	expectValue   T
	test.Control
}

func (tc optionalGetTC[T]) Test(t *testing.T) {
	value, present := tc.opt.Get()
	assert.Equal(t, tc.expectValue, value, "unexpected value")
	assert.Equal(t, tc.expectPresent, present, "unexpected value presence")
}

func TestOptional_Get(t *testing.T) {
	test.RunCases(t, test.Cases{
		// Test cases for documented examples
		"on empty int Optional": optionalGetTC[int]{
			opt:           Empty[int](),
			expectPresent: false,
			expectValue:   0,
		},
		"on non-empty int Optional with zero value": optionalGetTC[int]{
			opt:           Of(0),
			expectPresent: true,
			expectValue:   0,
		},
		"on non-empty int Optional with non-zero value": optionalGetTC[int]{
			opt:           Of(123),
			expectPresent: true,
			expectValue:   123,
		},
		"on empty string Optional": optionalGetTC[string]{
			opt:           Empty[string](),
			expectPresent: false,
			expectValue:   "",
		},
		"on non-empty string Optional with zero value": optionalGetTC[string]{
			opt:           Of(""),
			expectPresent: true,
			expectValue:   "",
		},
		"on non-empty string Optional with non-zero value": optionalGetTC[string]{
			opt:           Of("abc"),
			expectPresent: true,
			expectValue:   "abc",
		},
		// Other test cases...
	})
}

func BenchmarkOptional_IfPresent(b *testing.B) {
	opt := Of(123)
	for i := 0; i < b.N; i++ {
		opt.IfPresent(func(_ int) {})
	}
}

func ExampleOptional_IfPresent_int() {
	Empty[int]().IfPresent(example.PrintValue[int]) // Does nothing
	Of(0).IfPresent(example.PrintValue[int])
	Of(123).IfPresent(example.PrintValue[int])

	// Output:
	// 0
	// 123
}

func ExampleOptional_IfPresent_string() {
	Empty[string]().IfPresent(example.PrintValue[string]) // Does nothing
	Of("").IfPresent(example.PrintValue[string])
	Of("abc").IfPresent(example.PrintValue[string])

	// Output:
	// ""
	// "abc"
}

type optionalIfPresentTC[T any] struct {
	opt             Optional[T]
	expectCallCount uint
	test.Control
}

func (tc optionalIfPresentTC[T]) Test(t *testing.T) {
	var callCount uint
	tc.opt.IfPresent(func(value T) {
		callCount++
		assert.Equal(t, tc.opt.value, value)
	})
	assert.Equalf(t, tc.expectCallCount, callCount, "expected function to be called %v times", tc.expectCallCount)
}

func TestOptional_IfPresent(t *testing.T) {
	test.RunCases(t, test.Cases{
		// Test cases for documented examples
		"on empty int Optional": optionalIfPresentTC[int]{
			opt:             Empty[int](),
			expectCallCount: 0,
		},
		"on non-empty int Optional with zero value": optionalIfPresentTC[int]{
			opt:             Of(0),
			expectCallCount: 1,
		},
		"on non-empty int Optional with non-zero value": optionalIfPresentTC[int]{
			opt:             Of(123),
			expectCallCount: 1,
		},
		"on empty string Optional": optionalIfPresentTC[string]{
			opt:             Empty[string](),
			expectCallCount: 0,
		},
		"on non-empty string Optional with zero value": optionalIfPresentTC[string]{
			opt:             Of(""),
			expectCallCount: 1,
		},
		"on non-empty string Optional with non-zero value": optionalIfPresentTC[string]{
			opt:             Of("abc"),
			expectCallCount: 1,
		},
		// Other test cases...
	})
}

func BenchmarkOptional_IsEmpty(b *testing.B) {
	opt := Of(123)
	for i := 0; i < b.N; i++ {
		_ = opt.IsEmpty()
	}
}

func ExampleOptional_IsEmpty_int() {
	fmt.Println(Empty[int]().IsEmpty())
	fmt.Println(Of(0).IsEmpty())
	fmt.Println(Of(123).IsEmpty())

	// Output:
	// true
	// false
	// false
}

func ExampleOptional_IsEmpty_string() {
	fmt.Println(Empty[string]().IsEmpty())
	fmt.Println(Of("").IsEmpty())
	fmt.Println(Of("abc").IsEmpty())

	// Output:
	// true
	// false
	// false
}

type optionalIsEmptyTC[T any] struct {
	opt    Optional[T]
	expect bool
	test.Control
}

func (tc optionalIsEmptyTC[T]) Test(t *testing.T) {
	absent := tc.opt.IsEmpty()
	assert.Equal(t, tc.expect, absent, "unexpected value absence")
}

func TestOptional_IsEmpty(t *testing.T) {
	test.RunCases(t, test.Cases{
		// Test cases for documented examples
		"on empty int Optional": optionalIsEmptyTC[int]{
			opt:    Empty[int](),
			expect: true,
		},
		"on non-empty int Optional with zero value": optionalIsEmptyTC[int]{
			opt:    Of(0),
			expect: false,
		},
		"on non-empty int Optional with non-zero value": optionalIsEmptyTC[int]{
			opt:    Of(123),
			expect: false,
		},
		"on empty string Optional": optionalIsEmptyTC[string]{
			opt:    Empty[string](),
			expect: true,
		},
		"on non-empty string Optional with zero value": optionalIsEmptyTC[string]{
			opt:    Of(""),
			expect: false,
		},
		"on non-empty string Optional with non-zero value": optionalIsEmptyTC[string]{
			opt:    Of("abc"),
			expect: false,
		},
		// Other test cases...
	})
}

func BenchmarkOptional_IsPresent(b *testing.B) {
	opt := Of(123)
	for i := 0; i < b.N; i++ {
		_ = opt.IsPresent()
	}
}

func ExampleOptional_IsPresent_int() {
	fmt.Println(Empty[int]().IsPresent())
	fmt.Println(Of(0).IsPresent())
	fmt.Println(Of(123).IsPresent())

	// Output:
	// false
	// true
	// true
}

func ExampleOptional_IsPresent_string() {
	fmt.Println(Empty[string]().IsPresent())
	fmt.Println(Of("").IsPresent())
	fmt.Println(Of("abc").IsPresent())

	// Output:
	// false
	// true
	// true
}

type optionalIsPresentTC[T any] struct {
	opt    Optional[T]
	expect bool
	test.Control
}

func (tc optionalIsPresentTC[T]) Test(t *testing.T) {
	present := tc.opt.IsPresent()
	assert.Equal(t, tc.expect, present, "unexpected value presence")
}

func TestOptional_IsPresent(t *testing.T) {
	test.RunCases(t, test.Cases{
		// Test cases for documented examples
		"on empty int Optional": optionalIsPresentTC[int]{
			opt:    Empty[int](),
			expect: false,
		},
		"on non-empty int Optional with zero value": optionalIsPresentTC[int]{
			opt:    Of(0),
			expect: true,
		},
		"on non-empty int Optional with non-zero value": optionalIsPresentTC[int]{
			opt:    Of(123),
			expect: true,
		},
		"on empty string Optional": optionalIsPresentTC[string]{
			opt:    Empty[string](),
			expect: false,
		},
		"on non-empty string Optional with zero value": optionalIsPresentTC[string]{
			opt:    Of(""),
			expect: true,
		},
		"on non-empty string Optional with non-zero value": optionalIsPresentTC[string]{
			opt:    Of("abc"),
			expect: true,
		},
		// Other test cases...
	})
}

func BenchmarkOptional_IsZero(b *testing.B) {
	opt := Of(123)
	for i := 0; i < b.N; i++ {
		_ = opt.IsZero()
	}
}

func ExampleOptional_IsZero_int() {
	fmt.Println(Empty[int]().IsZero())
	fmt.Println(Of(0).IsZero())
	fmt.Println(Of(123).IsZero())

	// Output:
	// true
	// false
	// false
}

func ExampleOptional_IsZero_string() {
	fmt.Println(Empty[string]().IsZero())
	fmt.Println(Of("").IsZero())
	fmt.Println(Of("abc").IsZero())

	// Output:
	// true
	// false
	// false
}

type optionalIsZeroTC[T any] struct {
	opt    Optional[T]
	expect bool
	test.Control
}

func (tc optionalIsZeroTC[T]) Test(t *testing.T) {
	absent := tc.opt.IsZero()
	assert.Equal(t, tc.expect, absent, "unexpected value absence")
}

func TestOptional_IsZero(t *testing.T) {
	test.RunCases(t, test.Cases{
		// Test cases for documented examples
		"on empty int Optional": optionalIsZeroTC[int]{
			opt:    Empty[int](),
			expect: true,
		},
		"on non-empty int Optional with zero value": optionalIsZeroTC[int]{
			opt:    Of(0),
			expect: false,
		},
		"on non-empty int Optional with non-zero value": optionalIsZeroTC[int]{
			opt:    Of(123),
			expect: false,
		},
		"on empty string Optional": optionalIsZeroTC[string]{
			opt:    Empty[string](),
			expect: true,
		},
		"on non-empty string Optional with zero value": optionalIsZeroTC[string]{
			opt:    Of(""),
			expect: false,
		},
		"on non-empty string Optional with non-zero value": optionalIsZeroTC[string]{
			opt:    Of("abc"),
			expect: false,
		},
		// Other test cases...
	})
}

func BenchmarkOptional_MarshalJSON(b *testing.B) {
	opt := Of(123)
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(opt)
	}
}

func ExampleOptional_MarshalJSON() {
	// json omitempty option does not apply to zero value structs
	type MyStruct struct {
		Number Optional[int]    `json:"number,omitempty"`
		Text   Optional[string] `json:"text,omitempty"`
	}

	example.PrintMarshalled(json.Marshal(Of(MyStruct{})))
	example.PrintMarshalled(json.Marshal(Of(MyStruct{Number: Of(0), Text: Of("")})))
	example.PrintMarshalled(json.Marshal(Of(MyStruct{Number: Of(123), Text: Of("abc")})))

	// Output:
	// {"number":null,"text":null} <nil>
	// {"number":0,"text":""} <nil>
	// {"number":123,"text":"abc"} <nil>
}

func ExampleOptional_MarshalJSON_pointers() {
	type MyStruct struct {
		Number *Optional[int]    `json:"number,omitempty"`
		Text   *Optional[string] `json:"text,omitempty"`
	}

	example.PrintMarshalled(json.Marshal(Of(MyStruct{})))
	example.PrintMarshalled(json.Marshal(Of(MyStruct{Number: ptrs.Value(Of(0)), Text: ptrs.Value(Of(""))})))
	example.PrintMarshalled(json.Marshal(Of(MyStruct{Number: ptrs.Value(Of(123)), Text: ptrs.Value(Of("abc"))})))

	// Output:
	// {} <nil>
	// {"number":0,"text":""} <nil>
	// {"number":123,"text":"abc"} <nil>
}

type optionalMarshalJSONTC struct {
	value      any
	expectJSON string
	test.Control
}

func (tc optionalMarshalJSONTC) Test(t *testing.T) {
	b, err := json.Marshal(tc.value)
	assert.NoError(t, err, "unexpected error")
	assert.Equal(t, tc.expectJSON, string(b), "unexpected JSON")
}

func TestOptional_MarshalJSON(t *testing.T) {
	type Example struct {
		Int           Optional[int]     `json:"int"`
		String        Optional[string]  `json:"string"`
		IntOmit       Optional[int]     `json:"intOmit,omitempty"`
		StringOmit    Optional[string]  `json:"stringOmit,omitempty"`
		IntOmitPtr    *Optional[int]    `json:"intOmitPtr,omitempty"`
		StringOmitPtr *Optional[string] `json:"stringOmitPtr,omitempty"`
	}

	test.RunCases(t, test.Cases{
		"on empty int Optional": optionalMarshalJSONTC{
			value:      Empty[int](),
			expectJSON: `null`,
		},
		"on non-empty int Optional with zero value": optionalMarshalJSONTC{
			value:      Of(0),
			expectJSON: `0`,
		},
		"on non-empty int Optional with non-zero value": optionalMarshalJSONTC{
			value:      Of(123),
			expectJSON: `123`,
		},
		"on empty string Optional": optionalMarshalJSONTC{
			value:      Empty[string](),
			expectJSON: `null`,
		},
		"on non-empty string Optional with zero value": optionalMarshalJSONTC{
			value:      Of(""),
			expectJSON: `""`,
		},
		"on non-empty string Optional with non-zero value": optionalMarshalJSONTC{
			value:      Of("abc"),
			expectJSON: `"abc"`,
		},
		"on struct with empty Optionals": optionalMarshalJSONTC{
			value:      Example{},
			expectJSON: `{"int":null,"string":null,"intOmit":null,"stringOmit":null}`,
			// json omitempty option does not apply to zero value structs
		},
		"on struct with non-empty Optionals and zero field values": optionalMarshalJSONTC{
			value: Example{
				Int:           Of(0),
				String:        Of(""),
				IntOmit:       Of(0),
				StringOmit:    Of(""),
				IntOmitPtr:    ptrs.Value(Of(0)),
				StringOmitPtr: ptrs.Value(Of("")),
			},
			expectJSON: `{"int":0,"string":"","intOmit":0,"stringOmit":"","intOmitPtr":0,"stringOmitPtr":""}`,
		},
		"on struct with non-empty Optionals and non-zero field values": optionalMarshalJSONTC{
			value: Example{
				Int:           Of(123),
				String:        Of("abc"),
				IntOmit:       Of(123),
				StringOmit:    Of("abc"),
				IntOmitPtr:    ptrs.Value(Of(123)),
				StringOmitPtr: ptrs.Value(Of("abc")),
			},
			expectJSON: `{"int":123,"string":"abc","intOmit":123,"stringOmit":"abc","intOmitPtr":123,"stringOmitPtr":"abc"}`,
		},
	})
}

func BenchmarkOptional_MarshalXML(b *testing.B) {
	opt := Of(123)
	for i := 0; i < b.N; i++ {
		_, _ = xml.Marshal(opt)
	}
}

func ExampleOptional_MarshalXML() {
	type MyStruct struct {
		Number Optional[int]    `xml:"number,omitempty"`
		Text   Optional[string] `xml:"text,omitempty"`
	}

	example.PrintMarshalled(xml.Marshal(Of(MyStruct{})))
	example.PrintMarshalled(xml.Marshal(Of(MyStruct{Number: Of(0), Text: Of("")})))
	example.PrintMarshalled(xml.Marshal(Of(MyStruct{Number: Of(123), Text: Of("abc")})))

	// Output:
	// <MyStruct></MyStruct> <nil>
	// <MyStruct><number>0</number><text></text></MyStruct> <nil>
	// <MyStruct><number>123</number><text>abc</text></MyStruct> <nil>
}

type optionalMarshalXMLTC struct {
	value     any
	expectXML string
	test.Control
}

func (tc optionalMarshalXMLTC) Test(t *testing.T) {
	b, err := xml.Marshal(tc.value)
	assert.NoError(t, err, "unexpected error")
	assert.Equal(t, tc.expectXML, string(b), "unexpected XML")
}

func TestOptional_MarshalXML(t *testing.T) {
	type Example struct {
		Int           Optional[int]     `xml:"int"`
		String        Optional[string]  `xml:"string"`
		IntOmit       Optional[int]     `xml:"intOmit,omitempty"`
		StringOmit    Optional[string]  `xml:"stringOmit,omitempty"`
		IntOmitPtr    *Optional[int]    `xml:"intOmitPtr,omitempty"`
		StringOmitPtr *Optional[string] `xml:"stringOmitPtr,omitempty"`
	}

	test.RunCases(t, test.Cases{
		"on empty int Optional": optionalMarshalXMLTC{
			value:     Empty[int](),
			expectXML: ``,
		},
		"on non-empty int Optional with zero value": optionalMarshalXMLTC{
			value:     Of(0),
			expectXML: `<int>0</int>`,
		},
		"on non-empty int Optional with non-zero value": optionalMarshalXMLTC{
			value:     Of(123),
			expectXML: `<int>123</int>`,
		},
		"on empty string Optional": optionalMarshalXMLTC{
			value:     Empty[string](),
			expectXML: ``,
		},
		"on non-empty string Optional with zero value": optionalMarshalXMLTC{
			value:     Of(""),
			expectXML: `<string></string>`,
		},
		"on non-empty string Optional with non-zero value": optionalMarshalXMLTC{
			value:     Of("abc"),
			expectXML: `<string>abc</string>`,
		},
		"on struct with empty Optionals": optionalMarshalXMLTC{
			value:     Example{},
			expectXML: `<Example></Example>`,
		},
		"on struct with non-empty Optionals and zero field values": optionalMarshalXMLTC{
			value: Example{
				Int:           Of(0),
				String:        Of(""),
				IntOmit:       Of(0),
				StringOmit:    Of(""),
				IntOmitPtr:    ptrs.Value(Of(0)),
				StringOmitPtr: ptrs.Value(Of("")),
			},
			expectXML: `<Example><int>0</int><string></string><intOmit>0</intOmit><stringOmit></stringOmit><intOmitPtr>0</intOmitPtr><stringOmitPtr></stringOmitPtr></Example>`,
		},
		"on struct with non-empty Optionals and non-zero field values": optionalMarshalXMLTC{
			value: Example{
				Int:           Of(123),
				String:        Of("abc"),
				IntOmit:       Of(123),
				StringOmit:    Of("abc"),
				IntOmitPtr:    ptrs.Value(Of(123)),
				StringOmitPtr: ptrs.Value(Of("abc")),
			},
			expectXML: `<Example><int>123</int><string>abc</string><intOmit>123</intOmit><stringOmit>abc</stringOmit><intOmitPtr>123</intOmitPtr><stringOmitPtr>abc</stringOmitPtr></Example>`,
		},
	})
}

func BenchmarkOptional_MarshalYAML(b *testing.B) {
	opt := Of(123)
	for i := 0; i < b.N; i++ {
		_, _ = yaml.Marshal(opt)
	}
}

func ExampleOptional_MarshalYAML() {
	type MyStruct struct {
		Number Optional[int]    `yaml:"number,omitempty"`
		Text   Optional[string] `yaml:"text,omitempty"`
	}

	example.PrintMarshalled(yaml.Marshal(Of(MyStruct{})))
	example.PrintMarshalled(yaml.Marshal(Of(MyStruct{Number: Of(0), Text: Of("")})))
	example.PrintMarshalled(yaml.Marshal(Of(MyStruct{Number: Of(123), Text: Of("abc")})))

	// Output:
	// {} <nil>
	// number: 0
	// text: "" <nil>
	// number: 123
	// text: abc <nil>
}

type optionalMarshalYAMLTC struct {
	value      any
	expectYAML string
	test.Control
}

func (tc optionalMarshalYAMLTC) Test(t *testing.T) {
	b, err := yaml.Marshal(tc.value)
	assert.NoError(t, err, "unexpected error")
	assert.Equal(t, tc.expectYAML, strings.TrimSpace(string(b)), "unexpected YAML")
}

func TestOptional_MarshalYAML(t *testing.T) {
	type Example struct {
		Int           Optional[int]     `yaml:"int"`
		String        Optional[string]  `yaml:"string"`
		IntOmit       Optional[int]     `yaml:"intOmit,omitempty"`
		StringOmit    Optional[string]  `yaml:"stringOmit,omitempty"`
		IntOmitPtr    *Optional[int]    `yaml:"intOmitPtr,omitempty"`
		StringOmitPtr *Optional[string] `yaml:"stringOmitPtr,omitempty"`
	}

	test.RunCases(t, test.Cases{
		"on empty int Optional": optionalMarshalYAMLTC{
			value:      Empty[int](),
			expectYAML: `null`,
		},
		"on non-empty int Optional with zero value": optionalMarshalYAMLTC{
			value:      Of(0),
			expectYAML: `0`,
		},
		"on non-empty int Optional with non-zero value": optionalMarshalYAMLTC{
			value:      Of(123),
			expectYAML: `123`,
		},
		"on empty string Optional": optionalMarshalYAMLTC{
			value:      Empty[string](),
			expectYAML: `null`,
		},
		"on non-empty string Optional with zero value": optionalMarshalYAMLTC{
			value:      Of(""),
			expectYAML: `""`,
		},
		"on non-empty string Optional with non-zero value": optionalMarshalYAMLTC{
			value:      Of("abc"),
			expectYAML: `abc`,
		},
		"on struct with empty Optionals": optionalMarshalYAMLTC{
			value: Example{},
			expectYAML: `int: null
string: null`,
		},
		"on struct with non-empty Optionals and zero field values": optionalMarshalYAMLTC{
			value: Example{
				Int:           Of(0),
				String:        Of(""),
				IntOmit:       Of(0),
				StringOmit:    Of(""),
				IntOmitPtr:    ptrs.Value(Of(0)),
				StringOmitPtr: ptrs.Value(Of("")),
			},
			expectYAML: `int: 0
string: ""
intOmit: 0
stringOmit: ""
intOmitPtr: 0
stringOmitPtr: ""`,
		},
		"on struct with non-empty Optionals and non-zero field values": optionalMarshalYAMLTC{
			value: Example{
				Int:           Of(123),
				String:        Of("abc"),
				IntOmit:       Of(123),
				StringOmit:    Of("abc"),
				IntOmitPtr:    ptrs.Value(Of(123)),
				StringOmitPtr: ptrs.Value(Of("abc")),
			},
			expectYAML: `int: 123
string: abc
intOmit: 123
stringOmit: abc
intOmitPtr: 123
stringOmitPtr: abc`,
		},
	})
}

func BenchmarkOptional_OrElse(b *testing.B) {
	opt := Of(123)
	for i := 0; i < b.N; i++ {
		_ = opt.OrElse(-1)
	}
}

func ExampleOptional_OrElse_int() {
	defaultVal := -1

	example.PrintValue(Empty[int]().OrElse(defaultVal))
	example.PrintValue(Of(0).OrElse(defaultVal))
	example.PrintValue(Of(123).OrElse(defaultVal))

	// Output:
	// -1
	// 0
	// 123
}

func ExampleOptional_OrElse_string() {
	defaultVal := "unknown"

	example.PrintValue(Empty[string]().OrElse(defaultVal))
	example.PrintValue(Of("").OrElse(defaultVal))
	example.PrintValue(Of("abc").OrElse(defaultVal))

	// Output:
	// "unknown"
	// ""
	// "abc"
}

type optionalOrElseTC[T any] struct {
	opt    Optional[T]
	other  T
	expect T
	test.Control
}

func (tc optionalOrElseTC[T]) Test(t *testing.T) {
	value := tc.opt.OrElse(tc.other)
	assert.Equal(t, tc.expect, value, "unexpected value")
}

func TestOptional_OrElse(t *testing.T) {
	defaultInt := -1
	defaultString := "unknown"

	test.RunCases(t, test.Cases{
		// Test cases for documented examples
		"on empty int Optional": optionalOrElseTC[int]{
			opt:    Empty[int](),
			other:  defaultInt,
			expect: defaultInt,
		},
		"on non-empty int Optional with zero value": optionalOrElseTC[int]{
			opt:    Of(0),
			other:  defaultInt,
			expect: 0,
		},
		"on non-empty int Optional with non-zero value": optionalOrElseTC[int]{
			opt:    Of(123),
			other:  defaultInt,
			expect: 123,
		},
		"on empty string Optional": optionalOrElseTC[string]{
			opt:    Empty[string](),
			other:  defaultString,
			expect: defaultString,
		},
		"on non-empty string Optional with zero value": optionalOrElseTC[string]{
			opt:    Of(""),
			other:  defaultString,
			expect: "",
		},
		"on non-empty string Optional with non-zero value": optionalOrElseTC[string]{
			opt:    Of("abc"),
			other:  defaultString,
			expect: "abc",
		},
		// Other test cases...
	})
}

func BenchmarkOptional_OrElseGet(b *testing.B) {
	opt := Of(123)
	for i := 0; i < b.N; i++ {
		_ = opt.OrElseGet(func() int {
			return -1
		})
	}
}

func ExampleOptional_OrElseGet_int() {
	defaultFunc := func() int {
		return -1
	}

	example.PrintValue(Empty[int]().OrElseGet(defaultFunc))
	example.PrintValue(Of(0).OrElseGet(defaultFunc))
	example.PrintValue(Of(123).OrElseGet(defaultFunc))

	// Output:
	// -1
	// 0
	// 123
}

func ExampleOptional_OrElseGet_string() {
	defaultFunc := func() string {
		return "unknown"
	}

	example.PrintValue(Empty[string]().OrElseGet(defaultFunc))
	example.PrintValue(Of("").OrElseGet(defaultFunc))
	example.PrintValue(Of("abc").OrElseGet(defaultFunc))

	// Output:
	// "unknown"
	// ""
	// "abc"
}

type optionalOrElseGetTC[T any] struct {
	opt    Optional[T]
	other  func() T
	expect T
	test.Control
}

func (tc optionalOrElseGetTC[T]) Test(t *testing.T) {
	value := tc.opt.OrElseGet(tc.other)
	assert.Equal(t, tc.expect, value, "unexpected value")
}

func TestOptional_OrElseGet(t *testing.T) {
	defaultInt := -1
	defaultIntFunc := func() int {
		return defaultInt
	}
	defaultString := "unknown"
	defaultStringFunc := func() string {
		return defaultString
	}

	test.RunCases(t, test.Cases{
		// Test cases for documented examples
		"on empty int Optional": optionalOrElseGetTC[int]{
			opt:    Empty[int](),
			other:  defaultIntFunc,
			expect: defaultInt,
		},
		"on non-empty int Optional with zero value": optionalOrElseGetTC[int]{
			opt:    Of(0),
			other:  defaultIntFunc,
			expect: 0,
		},
		"on non-empty int Optional with non-zero value": optionalOrElseGetTC[int]{
			opt:    Of(123),
			other:  defaultIntFunc,
			expect: 123,
		},
		"on empty string Optional": optionalOrElseGetTC[string]{
			opt:    Empty[string](),
			other:  defaultStringFunc,
			expect: defaultString,
		},
		"on non-empty string Optional with zero value": optionalOrElseGetTC[string]{
			opt:    Of(""),
			other:  defaultStringFunc,
			expect: "",
		},
		"on non-empty string Optional with non-zero value": optionalOrElseGetTC[string]{
			opt:    Of("abc"),
			other:  defaultStringFunc,
			expect: "abc",
		},
		// Other test cases...
	})
}

func BenchmarkOptional_OrElseTryGet(b *testing.B) {
	opt := Of(123)
	for i := 0; i < b.N; i++ {
		_, _ = opt.OrElseTryGet(func() (int, error) {
			return -1, nil
		})
	}
}

func ExampleOptional_OrElseTryGet_int() {
	defaultFunc := func() (int, error) {
		return -1, nil
	}

	example.PrintTryValue(Empty[int]().OrElseTryGet(defaultFunc))
	example.PrintTryValue(Of(0).OrElseTryGet(defaultFunc))
	example.PrintTryValue(Of(123).OrElseTryGet(defaultFunc))

	// Output:
	// -1 <nil>
	// 0 <nil>
	// 123 <nil>
}

func ExampleOptional_OrElseTryGet_string() {
	var defaultStringUsed bool
	defaultFunc := func() (string, error) {
		if defaultStringUsed {
			return "", errors.New("default string already used")
		}
		defaultStringUsed = true
		return "unknown", nil
	}

	example.PrintTryValue(Empty[string]().OrElseTryGet(defaultFunc))
	example.PrintTryValue(Of("").OrElseTryGet(defaultFunc))
	example.PrintTryValue(Of("abc").OrElseTryGet(defaultFunc))
	example.PrintTryValue(Empty[string]().OrElseTryGet(defaultFunc))

	// Output:
	// "unknown" <nil>
	// "" <nil>
	// "abc" <nil>
	// "" "default string already used"
}

type optionalOrElseTryGetTC[T any] struct {
	opt         Optional[T]
	other       func() (T, error)
	expectError bool
	expectValue T
	test.Control
}

func (tc optionalOrElseTryGetTC[T]) Test(t *testing.T) {
	value, err := tc.opt.OrElseTryGet(tc.other)
	assert.Equalf(t, tc.expectError, err != nil, "unexpected error: %v", err)
	assert.Equal(t, tc.expectValue, value, "unexpected value")
}

func TestOptional_OrElseTryGet(t *testing.T) {
	defaultInt := -1
	defaultIntFunc := func() (int, error) {
		return defaultInt, nil
	}
	defaultString := "unknown"
	defaultStringFunc := func(err error) func() (string, error) {
		return func() (string, error) {
			if err != nil {
				return "", err
			}
			return defaultString, nil
		}
	}

	test.RunCases(t, test.Cases{
		// Test cases for documented examples
		"on empty int Optional": optionalOrElseTryGetTC[int]{
			opt:         Empty[int](),
			other:       defaultIntFunc,
			expectValue: defaultInt,
		},
		"on non-empty int Optional with zero value": optionalOrElseTryGetTC[int]{
			opt:         Of(0),
			other:       defaultIntFunc,
			expectValue: 0,
		},
		"on non-empty int Optional with non-zero value": optionalOrElseTryGetTC[int]{
			opt:         Of(123),
			other:       defaultIntFunc,
			expectValue: 123,
		},
		"on empty string Optional": optionalOrElseTryGetTC[string]{
			opt:         Empty[string](),
			other:       defaultStringFunc(nil),
			expectValue: defaultString,
		},
		"on non-empty string Optional with zero value": optionalOrElseTryGetTC[string]{
			opt:         Of(""),
			other:       defaultStringFunc(nil),
			expectValue: "",
		},
		"on non-empty string Optional with non-zero value": optionalOrElseTryGetTC[string]{
			opt:         Of("abc"),
			other:       defaultStringFunc(nil),
			expectValue: "abc",
		},
		"on empty string Optional triggering erroneous default call": optionalOrElseTryGetTC[string]{
			opt:         Empty[string](),
			other:       defaultStringFunc(errors.New("default string already used")),
			expectError: true,
		},
		// Other test cases...
	})
}

func BenchmarkOptional_Require(b *testing.B) {
	opt := Of(123)
	for i := 0; i < b.N; i++ {
		_ = opt.Require()
	}
}

func ExampleOptional_Require_int() {
	example.PrintValue(Of(0).Require())
	example.PrintValue(Of(123).Require())

	// Output:
	// 0
	// 123
}

func ExampleOptional_Require_panic() {
	defer func() {
		fmt.Println(recover())
	}()

	Empty[int]().Require()

	// Output: go-optional: value not present
}

func ExampleOptional_Require_string() {
	example.PrintValue(Of("").Require())
	example.PrintValue(Of("abc").Require())

	// Output:
	// ""
	// "abc"
}

type optionalRequireTC[T any] struct {
	opt         Optional[T]
	expectPanic bool
	expectValue T
	test.Control
}

func (tc optionalRequireTC[T]) Test(t *testing.T) {
	if tc.expectPanic {
		assert.Panics(t, func() {
			tc.opt.Require()
		}, "expected panic")
	} else {
		var value T
		assert.NotPanics(t, func() {
			value = tc.opt.Require()
		}, "unexpected panic")
		assert.Equal(t, tc.expectValue, value, "unexpected value")
	}
}

func TestOptional_Require(t *testing.T) {
	test.RunCases(t, test.Cases{
		// Test cases for documented examples
		"on empty int Optional": optionalRequireTC[int]{
			opt:         Empty[int](),
			expectPanic: true,
		},
		"on non-empty int Optional with zero value": optionalRequireTC[int]{
			opt:         Of(0),
			expectValue: 0,
		},
		"on non-empty int Optional with non-zero value": optionalRequireTC[int]{
			opt:         Of(123),
			expectValue: 123,
		},
		"on empty string Optional": optionalRequireTC[string]{
			opt:         Empty[string](),
			expectPanic: true,
		},
		"on non-empty string Optional with zero value": optionalRequireTC[string]{
			opt:         Of(""),
			expectValue: "",
		},
		"on non-empty string Optional with non-zero value": optionalRequireTC[string]{
			opt:         Of("abc"),
			expectValue: "abc",
		},
		// Other test cases...
	})
}

func BenchmarkOptional_String(b *testing.B) {
	opt := Of(123)
	for i := 0; i < b.N; i++ {
		_ = opt.String()
	}
}

func ExampleOptional_String_int() {
	fmt.Printf("%q\n", Empty[int]().String())
	fmt.Printf("%q\n", Of(0).String())
	fmt.Printf("%q\n", Of(123).String())

	// Output:
	// "<empty>"
	// "0"
	// "123"
}

func ExampleOptional_String_string() {
	fmt.Printf("%q\n", Empty[string]().String())
	fmt.Printf("%q\n", Of("").String())
	fmt.Printf("%q\n", Of("abc").String())

	// Output:
	// "<empty>"
	// ""
	// "abc"
}

type optionalStringTC[T any] struct {
	opt    Optional[T]
	expect string
	test.Control
}

func (tc optionalStringTC[T]) Test(t *testing.T) {
	value := tc.opt.String()
	assert.Equal(t, tc.expect, value, "unexpected string representation")
}

func TestOptional_String(t *testing.T) {
	test.RunCases(t, test.Cases{
		// Test cases for documented examples
		"on empty int Optional": optionalStringTC[int]{
			opt:    Empty[int](),
			expect: "<empty>",
		},
		"on non-empty int Optional with zero value": optionalStringTC[int]{
			opt:    Of(0),
			expect: "0",
		},
		"on non-empty int Optional with non-zero value": optionalStringTC[int]{
			opt:    Of(123),
			expect: "123",
		},
		"on empty string Optional": optionalStringTC[string]{
			opt:    Empty[string](),
			expect: "<empty>",
		},
		"on non-empty string Optional with zero value": optionalStringTC[string]{
			opt:    Of(""),
			expect: "",
		},
		"on non-empty string Optional with non-zero value": optionalStringTC[string]{
			opt:    Of("abc"),
			expect: "abc",
		},
		// Other test cases...
	})
}

func BenchmarkOptional_UnmarshalJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var opt Optional[int]
		_ = json.Unmarshal([]byte(`123`), &opt)
	}
}

func ExampleOptional_UnmarshalJSON() {
	type MyStruct struct {
		Number Optional[int]    `json:"number"`
		Text   Optional[string] `json:"text"`
	}

	inputs := []string{
		`{}`,
		`{"number":null,"text":null}`,
		`{"number":0,"text":""}`,
		`{"number":123,"text":"abc"}`,
	}

	for _, input := range inputs {
		var output MyStruct
		if err := json.Unmarshal([]byte(input), &output); err != nil {
			panic(err)
		}

		example.Print(output.Number)
		example.Print(output.Text)
	}

	// Output:
	// <empty>
	// <empty>
	// 0
	// ""
	// 0
	// ""
	// 123
	// "abc"
}

type optionalUnmarshalJSONTC[T any] struct {
	json   string
	expect T
	test.Control
}

func (tc optionalUnmarshalJSONTC[T]) Test(t *testing.T) {
	var value T
	err := json.Unmarshal([]byte(tc.json), &value)
	assert.NoError(t, err, "unexpected error")
	assert.Equal(t, tc.expect, value, "unexpected value")
}

func TestOptional_UnmarshalJSON(t *testing.T) {
	type Example struct {
		Int       Optional[int]     `json:"int"`
		String    Optional[string]  `json:"string"`
		IntPtr    *Optional[int]    `json:"intPtr"`
		StringPtr *Optional[string] `json:"stringPtr"`
	}

	test.RunCases(t, test.Cases{
		"on empty int Optional": optionalUnmarshalJSONTC[Optional[int]]{
			json:   `null`,
			expect: Of(0),
		},
		"on non-empty int Optional with zero value": optionalUnmarshalJSONTC[Optional[int]]{
			json:   `0`,
			expect: Of(0),
		},
		"on non-empty int Optional with non-zero value": optionalUnmarshalJSONTC[Optional[int]]{
			json:   `123`,
			expect: Of(123),
		},
		"on empty string Optional": optionalUnmarshalJSONTC[Optional[string]]{
			json:   `null`,
			expect: Of(""),
		},
		"on non-empty string Optional with zero value": optionalUnmarshalJSONTC[Optional[string]]{
			json:   `""`,
			expect: Of(""),
		},
		"on non-empty string Optional with non-zero value": optionalUnmarshalJSONTC[Optional[string]]{
			json:   `"abc"`,
			expect: Of("abc"),
		},
		"on struct with empty Optionals": optionalUnmarshalJSONTC[Example]{
			json:   `{}`,
			expect: Example{},
		},
		"on struct with non-empty Optionals and zero field values": optionalUnmarshalJSONTC[Example]{
			json: `{"int":0,"string":"","intPtr":0,"stringPtr":""}`,
			expect: Example{
				Int:       Of(0),
				String:    Of(""),
				IntPtr:    ptrs.Value(Of(0)),
				StringPtr: ptrs.Value(Of("")),
			},
		},
		"on struct with non-empty Optionals and non-zero field values": optionalUnmarshalJSONTC[Example]{
			json: `{"int":123,"string":"abc","intPtr":123,"stringPtr":"abc"}`,
			expect: Example{
				Int:       Of(123),
				String:    Of("abc"),
				IntPtr:    ptrs.Value(Of(123)),
				StringPtr: ptrs.Value(Of("abc")),
			},
		},
	})
}

func BenchmarkOptional_UnmarshalXML(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var opt Optional[int]
		_ = xml.Unmarshal([]byte(`<int>123</int>`), &opt)
	}
}

func ExampleOptional_UnmarshalXML() {
	type MyStruct struct {
		Number Optional[int]    `xml:"number"`
		Text   Optional[string] `xml:"text"`
	}

	inputs := []string{
		`<MyStruct></MyStruct>`,
		`<MyStruct><number>0</number><text></text></MyStruct>`,
		`<MyStruct><number>123</number><text>abc</text></MyStruct>`,
	}

	for _, input := range inputs {
		var output MyStruct
		if err := xml.Unmarshal([]byte(input), &output); err != nil {
			panic(err)
		}

		example.Print(output.Number)
		example.Print(output.Text)
	}

	// Output:
	// <empty>
	// <empty>
	// 0
	// ""
	// 123
	// "abc"
}

type optionalUnmarshalXMLTC[T any] struct {
	xml    string
	expect T
	test.Control
}

func (tc optionalUnmarshalXMLTC[T]) Test(t *testing.T) {
	var value T
	err := xml.Unmarshal([]byte(tc.xml), &value)
	assert.NoError(t, err, "unexpected error")
	assert.Equal(t, tc.expect, value, "unexpected value")
}

func TestOptional_UnmarshalXML(t *testing.T) {
	type Example struct {
		Int       Optional[int]     `xml:"int"`
		String    Optional[string]  `xml:"string"`
		IntPtr    *Optional[int]    `xml:"intPtr"`
		StringPtr *Optional[string] `xml:"stringPtr"`
	}

	test.RunCases(t, test.Cases{
		"on empty int Optional": optionalUnmarshalXMLTC[Optional[int]]{
			xml:    `<int/>`,
			expect: Of(0),
		},
		"on non-empty int Optional with zero value": optionalUnmarshalXMLTC[Optional[int]]{
			xml:    `<int>0</int>`,
			expect: Of(0),
		},
		"on non-empty int Optional with non-zero value": optionalUnmarshalXMLTC[Optional[int]]{
			xml:    `<int>123</int>`,
			expect: Of(123),
		},
		"on empty string Optional": optionalUnmarshalXMLTC[Optional[string]]{
			xml:    `<string/>`,
			expect: Of(""),
		},
		"on non-empty string Optional with zero value": optionalUnmarshalXMLTC[Optional[string]]{
			xml:    `<string></string>`,
			expect: Of(""),
		},
		"on non-empty string Optional with non-zero value": optionalUnmarshalXMLTC[Optional[string]]{
			xml:    `<string>abc</string>`,
			expect: Of("abc"),
		},
		"on struct with empty Optionals": optionalUnmarshalXMLTC[Example]{
			xml:    `<Example></Example>`,
			expect: Example{},
		},
		"on struct with non-empty Optionals and zero field values": optionalUnmarshalXMLTC[Example]{
			xml: `<Example><int>0</int><string></string><intPtr>0</intPtr><stringPtr></stringPtr></Example>`,
			expect: Example{
				Int:       Of(0),
				String:    Of(""),
				IntPtr:    ptrs.Value(Of(0)),
				StringPtr: ptrs.Value(Of("")),
			},
		},
		"on struct with non-empty Optionals and non-zero field values": optionalUnmarshalXMLTC[Example]{
			xml: `<Example><int>123</int><string>abc</string><intPtr>123</intPtr><stringPtr>abc</stringPtr></Example>`,
			expect: Example{
				Int:       Of(123),
				String:    Of("abc"),
				IntPtr:    ptrs.Value(Of(123)),
				StringPtr: ptrs.Value(Of("abc")),
			},
		},
	})
}

func BenchmarkOptional_UnmarshalYAML(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var opt Optional[int]
		_ = yaml.Unmarshal([]byte(`123`), &opt)
	}
}

func ExampleOptional_UnmarshalYAML() {
	type MyStruct struct {
		Number Optional[int]    `json:"number"`
		Text   Optional[string] `json:"text"`
	}

	inputs := []string{
		`{}`,
		`number: null
text: null`,
		`number: 0
text: ""`,
		`number: 123
text: abc`,
	}

	for _, input := range inputs {
		var output MyStruct
		if err := yaml.Unmarshal([]byte(input), &output); err != nil {
			panic(err)
		}

		example.Print(output.Number)
		example.Print(output.Text)
	}

	// Output:
	// <empty>
	// <empty>
	// <empty>
	// <empty>
	// 0
	// ""
	// 123
	// "abc"
}

type optionalUnmarshalYAMLTC[T any] struct {
	yaml   string
	expect T
	test.Control
}

func (tc optionalUnmarshalYAMLTC[T]) Test(t *testing.T) {
	var value T
	err := yaml.Unmarshal([]byte(tc.yaml), &value)
	assert.NoError(t, err, "unexpected error")
	assert.Equal(t, tc.expect, value, "unexpected value")
}

func TestOptional_UnmarshalYAML(t *testing.T) {
	type Example struct {
		Int       Optional[int]     `yaml:"int"`
		String    Optional[string]  `yaml:"string"`
		IntPtr    *Optional[int]    `yaml:"intPtr"`
		StringPtr *Optional[string] `yaml:"stringPtr"`
	}

	test.RunCases(t, test.Cases{
		"on empty int Optional": optionalUnmarshalYAMLTC[Optional[int]]{
			yaml:   `null`,
			expect: Empty[int](),
		},
		"on non-empty int Optional with zero value": optionalUnmarshalYAMLTC[Optional[int]]{
			yaml:   `0`,
			expect: Of(0),
		},
		"on non-empty int Optional with non-zero value": optionalUnmarshalYAMLTC[Optional[int]]{
			yaml:   `123`,
			expect: Of(123),
		},
		"on empty string Optional": optionalUnmarshalYAMLTC[Optional[string]]{
			yaml:   `null`,
			expect: Empty[string](),
		},
		"on non-empty string Optional with zero value": optionalUnmarshalYAMLTC[Optional[string]]{
			yaml:   `""`,
			expect: Of(""),
		},
		"on non-empty string Optional with non-zero value": optionalUnmarshalYAMLTC[Optional[string]]{
			yaml:   `"abc"`,
			expect: Of("abc"),
		},
		"on struct with empty Optionals": optionalUnmarshalYAMLTC[Example]{
			yaml:   `{}`,
			expect: Example{},
		},
		"on struct with non-empty Optionals and zero field values": optionalUnmarshalYAMLTC[Example]{
			yaml: `int: 0
string: ""
intPtr: 0
stringPtr: ""`,
			expect: Example{
				Int:       Of(0),
				String:    Of(""),
				IntPtr:    ptrs.Value(Of(0)),
				StringPtr: ptrs.Value(Of("")),
			},
		},
		"on struct with non-empty Optionals and non-zero field values": optionalUnmarshalYAMLTC[Example]{
			yaml: `int: 123
string: abc
intPtr: 123
stringPtr: abc`,
			expect: Example{
				Int:       Of(123),
				String:    Of("abc"),
				IntPtr:    ptrs.Value(Of(123)),
				StringPtr: ptrs.Value(Of("abc")),
			},
		},
	})
}

func BenchmarkCompare(b *testing.B) {
	x := Of(123)
	y := Of(-123)
	for i := 0; i < b.N; i++ {
		Compare(x, y)
	}
}

func ExampleCompare_int() {
	fmt.Println(Compare(Empty[int](), Of(0)))
	fmt.Println(Compare(Of(0), Of(123)))

	fmt.Println(Compare(Empty[int](), Empty[int]()))
	fmt.Println(Compare(Of(0), Of(0)))
	fmt.Println(Compare(Of(123), Of(123)))

	fmt.Println(Compare(Of(0), Empty[int]()))
	fmt.Println(Compare(Of(123), Of(0)))

	// Output:
	// -1
	// -1
	// 0
	// 0
	// 0
	// 1
	// 1
}

func ExampleCompare_string() {
	fmt.Println(Compare(Empty[string](), Of("")))
	fmt.Println(Compare(Of(""), Of("abc")))

	fmt.Println(Compare(Empty[string](), Empty[string]()))
	fmt.Println(Compare(Of(""), Of("")))
	fmt.Println(Compare(Of("abc"), Of("abc")))

	fmt.Println(Compare(Of(""), Empty[string]()))
	fmt.Println(Compare(Of("abc"), Of("")))

	// Output:
	// -1
	// -1
	// 0
	// 0
	// 0
	// 1
	// 1
}

type compareTC[T cmp.Ordered] struct {
	x      Optional[T]
	y      Optional[T]
	expect int
	test.Control
}

func (tc compareTC[T]) Test(t *testing.T) {
	actual := Compare(tc.x, tc.y)
	assert.Equal(t, tc.expect, actual, "unexpected comparison result")
}

func TestCompare(t *testing.T) {
	test.RunCases(t, test.Cases{
		// Test cases for documented examples
		"with empty int Optional and non-empty int Optional with zero value": compareTC[int]{
			x:      Empty[int](),
			y:      Of(0),
			expect: -1,
		},
		"with non-empty int Optional with zero value and non-empty int Optional with positive non-zero value": compareTC[int]{
			x:      Of(0),
			y:      Of(123),
			expect: -1,
		},
		"with two empty int Optionals": compareTC[int]{
			x:      Empty[int](),
			y:      Empty[int](),
			expect: 0,
		},
		"with two non-empty int Optionals with zero values": compareTC[int]{
			x:      Of(0),
			y:      Of(0),
			expect: 0,
		},
		"with two non-empty int Optionals with same non-zero values": compareTC[int]{
			x:      Of(123),
			y:      Of(123),
			expect: 0,
		},
		"with non-empty int Optional with zero value and empty int Optional": compareTC[int]{
			x:      Of(0),
			y:      Empty[int](),
			expect: 1,
		},
		"with non-empty int Optional with positive non-zero value and non-empty int Optional with zero value": compareTC[int]{
			x:      Of(123),
			y:      Of(0),
			expect: 1,
		},
		// Other test cases...
	})
}

func BenchmarkEmpty(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Empty[int]()
	}
}

func ExampleEmpty_int() {
	example.Print(Empty[int]())

	// Output: <empty>
}

func ExampleEmpty_string() {
	example.Print(Empty[string]())

	// Output: <empty>
}

type emptyTC[T any] struct {
	test.Control
}

func (tc emptyTC[T]) Test(t *testing.T) {
	opt := Empty[T]()
	value, present := opt.Get()
	assert.Zero(t, value, "expected zero value")
	assert.False(t, present, "expected emptiness")
}

func TestEmpty(t *testing.T) {
	test.RunCases(t, test.Cases{
		// Test cases for documented examples
		"with int":    emptyTC[int]{},
		"with string": emptyTC[string]{},
		// Other test cases...
	})
}

func BenchmarkFind(b *testing.B) {
	opts := []Optional[int]{Empty[int](), Empty[int](), Of(123)}
	for i := 0; i < b.N; i++ {
		_ = Find(opts...)
	}
}

func ExampleFind_int() {
	example.Print(Find[int]())
	example.Print(Find(Empty[int]()))
	example.Print(Find(Empty[int](), Of(0), Of(123)))
	example.Print(Find(Empty[int](), Of(123), Of(0)))

	// Output:
	// <empty>
	// <empty>
	// 0
	// 123
}

func ExampleFind_string() {
	example.Print(Find[string]())
	example.Print(Find(Empty[string]()))
	example.Print(Find(Empty[string](), Of(""), Of("abc")))
	example.Print(Find(Empty[string](), Of("abc"), Of("")))

	// Output:
	// <empty>
	// <empty>
	// ""
	// "abc"
}

type findTC[T any] struct {
	opts          []Optional[T]
	expectPresent bool
	expectValue   T
	test.Control
}

func (tc findTC[T]) Test(t *testing.T) {
	opt := Find(tc.opts...)
	value, present := opt.Get()
	assert.Equal(t, tc.expectValue, value, "unexpected value")
	assert.Equal(t, tc.expectPresent, present, "unexpected value presence")
}

func TestFind(t *testing.T) {
	test.RunCases(t, test.Cases{
		// Test cases for documented examples
		"with no int Optionals": findTC[int]{
			expectPresent: false,
			expectValue:   0,
		},
		"with empty int Optional": findTC[int]{
			opts:          []Optional[int]{Empty[int]()},
			expectPresent: false,
			expectValue:   0,
		},
		"with an empty int Optional and two non-empty int Optionals": findTC[int]{
			opts: []Optional[int]{
				Empty[int](),
				Of(0),
				Of(123),
			},
			expectPresent: true,
			expectValue:   0,
		},
		"with no string Optionals": findTC[string]{
			expectPresent: false,
			expectValue:   "",
		},
		"with empty string Optional": findTC[string]{
			opts:          []Optional[string]{Empty[string]()},
			expectPresent: false,
			expectValue:   "",
		},
		"with an empty string Optional and two non-empty string Optionals": findTC[string]{
			opts: []Optional[string]{
				Empty[string](),
				Of("abc"),
				Of(""),
			},
			expectPresent: true,
			expectValue:   "abc",
		},
		// Other test cases...
	})
}

func BenchmarkFlatMap(b *testing.B) {
	toString := func(value int) Optional[string] {
		if value == 0 {
			return Empty[string]()
		}
		return Of(strconv.FormatInt(int64(value), 10))
	}
	opt := Of(123)
	for i := 0; i < b.N; i++ {
		_ = FlatMap(opt, toString)
	}
}

func ExampleFlatMap_int() {
	mapper := func(value int) Optional[string] {
		if value == 0 {
			return Empty[string]()
		}
		return Of(strconv.FormatInt(int64(value), 10))
	}

	example.Print(FlatMap(Empty[int](), mapper))
	example.Print(FlatMap(Of(0), mapper))
	example.Print(FlatMap(Of(123), mapper))

	// Output:
	// <empty>
	// <empty>
	// "123"
}

func ExampleFlatMap_string() {
	mapper := func(value string) Optional[int] {
		if value == "" {
			return Empty[int]()
		}
		i, err := strconv.ParseInt(value, 10, 0)
		if err != nil {
			panic(err)
		}
		return OfZeroable(int(i))
	}

	example.Print(FlatMap(Empty[string](), mapper))
	example.Print(FlatMap(Of(""), mapper))
	example.Print(FlatMap(Of("0"), mapper))
	example.Print(FlatMap(Of("123"), mapper))

	// Output:
	// <empty>
	// <empty>
	// <empty>
	// 123
}

type flatMapTC[T, M any] struct {
	opt           Optional[T]
	fn            func(value T) Optional[M]
	expectPresent bool
	expectValue   M
	test.Control
}

func (tc flatMapTC[T, M]) Test(t *testing.T) {
	opt := FlatMap(tc.opt, tc.fn)
	value, present := opt.Get()
	assert.Equal(t, tc.expectValue, value, "unexpected value")
	assert.Equal(t, tc.expectPresent, present, "unexpected value presence")
}

func TestFlatMap(t *testing.T) {
	toInt := func(value string) Optional[int] {
		if value == "" {
			return Empty[int]()
		}
		i, err := strconv.ParseInt(value, 10, 0)
		if err != nil {
			panic(err)
		}
		return OfZeroable(int(i))
	}
	toString := func(value int) Optional[string] {
		if value == 0 {
			return Empty[string]()
		}
		return Of(strconv.FormatInt(int64(value), 10))
	}

	test.RunCases(t, test.Cases{
		// Test cases for documented examples
		"with empty int Optional": flatMapTC[int, string]{
			opt:           Empty[int](),
			fn:            toString,
			expectPresent: false,
		},
		"with non-empty int Optional with zero value": flatMapTC[int, string]{
			opt:           Of(0),
			fn:            toString,
			expectPresent: false,
		},
		"with non-empty int Optional with non-zero value": flatMapTC[int, string]{
			opt:           Of(123),
			fn:            toString,
			expectPresent: true,
			expectValue:   "123",
		},
		"with empty string Optional": flatMapTC[string, int]{
			opt:           Empty[string](),
			fn:            toInt,
			expectPresent: false,
		},
		"with non-empty string Optional with zero value": flatMapTC[string, int]{
			opt:           Of(""),
			fn:            toInt,
			expectPresent: false,
		},
		"with non-empty string Optional with zero-representing value": flatMapTC[string, int]{
			opt:           Of("0"),
			fn:            toInt,
			expectPresent: false,
		},
		"with non-empty string Optional with non-zero-representing value": flatMapTC[string, int]{
			opt:           Of("123"),
			fn:            toInt,
			expectPresent: true,
			expectValue:   123,
		},
		// Other test cases...
	})
}

func BenchmarkGetAny(b *testing.B) {
	opts := []Optional[int]{Empty[int](), Of(0), Of(123)}
	for i := 0; i < b.N; i++ {
		_ = GetAny(opts...)
	}
}

func ExampleGetAny_int() {
	example.PrintValues(GetAny[int]())
	example.PrintValues(GetAny(Empty[int]()))
	example.PrintValues(GetAny(Empty[int](), Of(0), Of(123)))

	// Output:
	// []
	// []
	// [0 123]
}

func ExampleGetAny_string() {
	example.PrintValues(GetAny[string]())
	example.PrintValues(GetAny(Empty[string]()))
	example.PrintValues(GetAny(Empty[string](), Of("abc"), Of("")))

	// Output:
	// []
	// []
	// ["abc" ""]
}

type getAnyTC[T any] struct {
	opts   []Optional[T]
	expect []T
	test.Control
}

func (tc getAnyTC[T]) Test(t *testing.T) {
	actual := GetAny(tc.opts...)
	assert.Equal(t, tc.expect, actual, "unexpected values")
}

func TestGetAny(t *testing.T) {
	test.RunCases(t, test.Cases{
		// Test cases for documented examples
		"with no int Optionals": getAnyTC[int]{
			expect: nil,
		},
		"with empty int Optional": getAnyTC[int]{
			opts:   []Optional[int]{Empty[int]()},
			expect: nil,
		},
		"with an empty int Optional and two non-empty int Optionals": getAnyTC[int]{
			opts: []Optional[int]{
				Empty[int](),
				Of(0),
				Of(123),
			},
			expect: []int{0, 123},
		},
		"with no string Optionals": getAnyTC[string]{
			expect: nil,
		},
		"with empty string Optional": getAnyTC[string]{
			opts:   []Optional[string]{Empty[string]()},
			expect: nil,
		},
		"with an empty string Optional and two non-empty string Optionals": getAnyTC[string]{
			opts: []Optional[string]{
				Empty[string](),
				Of("abc"),
				Of(""),
			},
			expect: []string{"abc", ""},
		},
		// Other test cases...
	})
}

func BenchmarkMap(b *testing.B) {
	toString := func(value int) string {
		return strconv.FormatInt(int64(value), 10)
	}
	opt := Of(123)
	for i := 0; i < b.N; i++ {
		_ = Map(opt, toString)
	}
}

func ExampleMap_int() {
	mapper := func(value int) string {
		return strconv.FormatInt(int64(value), 10)
	}

	example.Print(Map(Empty[int](), mapper))
	example.Print(Map(Of(0), mapper))
	example.Print(Map(Of(123), mapper))

	// Output:
	// <empty>
	// "0"
	// "123"
}

func ExampleMap_string() {
	mapper := func(value string) int {
		i, err := strconv.ParseInt(value, 10, 0)
		if err != nil {
			panic(err)
		}
		return int(i)
	}

	example.Print(Map(Empty[string](), mapper))
	example.Print(Map(Of("0"), mapper))
	example.Print(Map(Of("123"), mapper))

	// Output:
	// <empty>
	// 0
	// 123
}

type mapTC[T, M any] struct {
	opt           Optional[T]
	fn            func(value T) M
	expectPresent bool
	expectValue   M
	test.Control
}

func (tc mapTC[T, M]) Test(t *testing.T) {
	opt := Map(tc.opt, tc.fn)
	value, present := opt.Get()
	assert.Equal(t, tc.expectValue, value, "unexpected value")
	assert.Equal(t, tc.expectPresent, present, "unexpected value presence")
}

func TestMap(t *testing.T) {
	toInt := func(value string) int {
		i, err := strconv.ParseInt(value, 10, 0)
		if err != nil {
			panic(err)
		}
		return int(i)
	}
	toString := func(value int) string {
		return strconv.FormatInt(int64(value), 10)
	}

	test.RunCases(t, test.Cases{
		// Test cases for documented examples
		"with empty int Optional": mapTC[int, string]{
			opt:           Empty[int](),
			fn:            toString,
			expectPresent: false,
		},
		"with non-empty int Optional with zero value": mapTC[int, string]{
			opt:           Of(0),
			fn:            toString,
			expectPresent: true,
			expectValue:   "0",
		},
		"with non-empty int Optional with non-zero value": mapTC[int, string]{
			opt:           Of(123),
			fn:            toString,
			expectPresent: true,
			expectValue:   "123",
		},
		"with empty string Optional": mapTC[string, int]{
			opt:           Empty[string](),
			fn:            toInt,
			expectPresent: false,
		},
		"with non-empty string Optional with zero-representing value": mapTC[string, int]{
			opt:           Of("0"),
			fn:            toInt,
			expectPresent: true,
			expectValue:   0,
		},
		"with non-empty string Optional with non-zero-representing value": mapTC[string, int]{
			opt:           Of("123"),
			fn:            toInt,
			expectPresent: true,
			expectValue:   123,
		},
		// Other test cases...
	})
}

func BenchmarkMustFind(b *testing.B) {
	opts := []Optional[int]{Empty[int](), Of(0), Of(123)}
	for i := 0; i < b.N; i++ {
		_ = MustFind(opts...)
	}
}

func ExampleMustFind_int() {
	example.PrintValue(MustFind(Empty[int](), Of(0), Of(123)))

	// Output: 0
}

func ExampleMustFind_panic() {
	defer func() {
		fmt.Println(recover())
	}()

	MustFind(Empty[int]())

	// Output: go-optional: value not present
}

func ExampleMustFind_string() {
	example.PrintValue(MustFind(Empty[string](), Of("abc"), Of("")))

	// Output: "abc"
}

type mustFindTC[T any] struct {
	opts        []Optional[T]
	expectPanic bool
	expectValue T
	test.Control
}

func (tc mustFindTC[T]) Test(t *testing.T) {
	if tc.expectPanic {
		assert.Panics(t, func() {
			MustFind(tc.opts...)
		}, "expected panic")
	} else {
		var value T
		assert.NotPanics(t, func() {
			value = MustFind(tc.opts...)
		}, "unexpected panic")
		assert.Equal(t, tc.expectValue, value, "unexpected value")
	}
}

func TestMustFind(t *testing.T) {
	test.RunCases(t, test.Cases{
		// Test cases for documented examples
		"with no int Optionals": mustFindTC[int]{
			expectPanic: true,
		},
		"with empty int Optional": mustFindTC[int]{
			opts:        []Optional[int]{Empty[int]()},
			expectPanic: true,
		},
		"with an empty int Optional and two non-empty int Optionals": mustFindTC[int]{
			opts: []Optional[int]{
				Empty[int](),
				Of(0),
				Of(123),
			},
			expectValue: 0,
		},
		"with no string Optionals": mustFindTC[string]{
			expectPanic: true,
		},
		"with empty string Optional": mustFindTC[string]{
			opts:        []Optional[string]{Empty[string]()},
			expectPanic: true,
		},
		"with an empty string Optional and two non-empty string Optionals": mustFindTC[string]{
			opts: []Optional[string]{
				Empty[string](),
				Of("abc"),
				Of(""),
			},
			expectValue: "abc",
		},
		// Other test cases...
	})
}

func BenchmarkOf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Of(123)
	}
}

func ExampleOf_int() {
	example.Print(Of(0))
	example.Print(Of(123))

	// Output:
	// 0
	// 123
}

func ExampleOf_int_pointer() {
	example.Print(Of((*int)(nil)))
	example.Print(Of(ptrs.ZeroInt()))
	example.Print(Of(ptrs.Int(123)))

	// Output:
	// <nil>
	// &0
	// &123
}

func ExampleOf_string() {
	example.Print(Of(""))
	example.Print(Of("abc"))

	// Output:
	// ""
	// "abc"
}

func ExampleOf_string_pointer() {
	example.Print(Of((*string)(nil)))
	example.Print(Of(ptrs.ZeroString()))
	example.Print(Of(ptrs.String("abc")))

	// Output:
	// <nil>
	// &""
	// &"abc"
}

type ofTC[T any] struct {
	value T
	test.Control
}

func (tc ofTC[T]) Test(t *testing.T) {
	opt := Of(tc.value)
	value, present := opt.Get()
	assert.Equal(t, tc.value, value, "unexpected value")
	assert.True(t, present, "expected presence")
}

func TestOf(t *testing.T) {
	test.RunCases(t, test.Cases{
		// Test cases for documented examples
		"with zero int": ofTC[int]{
			value: 0,
		},
		"with non-zero int": ofTC[int]{
			value: 123,
		},
		"with nil int pointer": ofTC[*int]{
			value: nil,
		},
		"with zero int pointer": ofTC[*int]{
			value: ptrs.ZeroInt(),
		},
		"with non-zero int pointer": ofTC[*int]{
			value: ptrs.Int(123),
		},
		"with zero string": ofTC[string]{
			value: "",
		},
		"with non-zero string": ofTC[string]{
			value: "abc",
		},
		"with nil string pointer": ofTC[*string]{
			value: nil,
		},
		"with zero string pointer": ofTC[*string]{
			value: ptrs.ZeroString(),
		},
		"with non-zero string pointer": ofTC[*string]{
			value: ptrs.String("abc"),
		},
		// Other test cases...
	})
}

func BenchmarkOfNillable(b *testing.B) {
	value := 123
	for i := 0; i < b.N; i++ {
		OfNillable(&value)
	}
}

func ExampleOfNillable_int() {
	example.Print(OfNillable(0))
	example.Print(OfNillable(123))

	// Output:
	// 0
	// 123
}

func ExampleOfNillable_int_pointer() {
	example.Print(OfNillable((*int)(nil)))
	example.Print(OfNillable(ptrs.ZeroInt()))
	example.Print(OfNillable(ptrs.Int(123)))

	// Output:
	// <empty>
	// &0
	// &123
}

func ExampleOfNillable_string() {
	example.Print(OfNillable(""))
	example.Print(OfNillable("abc"))

	// Output:
	// ""
	// "abc"
}

func ExampleOfNillable_string_pointer() {
	example.Print(OfNillable((*string)(nil)))
	example.Print(OfNillable(ptrs.ZeroString()))
	example.Print(OfNillable(ptrs.String("abc")))

	// Output:
	// <empty>
	// &""
	// &"abc"
}

type ofNillableTC[T any] struct {
	value         T
	expectPresent bool
	test.Control
}

func (tc ofNillableTC[T]) Test(t *testing.T) {
	opt := OfNillable(tc.value)
	value, present := opt.Get()
	assert.Equal(t, tc.value, value, "unexpected value")
	assert.Equal(t, tc.expectPresent, present, "unexpected value presence")
}

func TestOfNillable(t *testing.T) {
	test.RunCases(t, test.Cases{
		// Test cases for documented examples
		"with zero int": ofNillableTC[int]{
			value:         0,
			expectPresent: true,
		},
		"with non-zero int": ofNillableTC[int]{
			value:         123,
			expectPresent: true,
		},
		"with nil int pointer": ofNillableTC[*int]{
			value:         nil,
			expectPresent: false,
		},
		"with zero int pointer": ofNillableTC[*int]{
			value:         ptrs.ZeroInt(),
			expectPresent: true,
		},
		"with non-zero int pointer": ofNillableTC[*int]{
			value:         ptrs.Int(123),
			expectPresent: true,
		},
		"with zero string": ofNillableTC[string]{
			value:         "",
			expectPresent: true,
		},
		"with non-zero string": ofNillableTC[string]{
			value:         "abc",
			expectPresent: true,
		},
		"with nil string pointer": ofNillableTC[*string]{
			value:         nil,
			expectPresent: false,
		},
		"with zero string pointer": ofNillableTC[*string]{
			value:         ptrs.ZeroString(),
			expectPresent: true,
		},
		"with non-zero string pointer": ofNillableTC[*string]{
			value:         ptrs.String("abc"),
			expectPresent: true,
		},
		// Other test cases...
	})
}

func BenchmarkOfPointer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		OfPointer(123)
	}
}

func ExampleOfPointer_int() {
	example.Print(OfPointer(0))
	example.Print(OfPointer(123))

	// Output:
	// &0
	// &123
}

func ExampleOfPointer_string() {
	example.Print(OfPointer(""))
	example.Print(OfPointer("abc"))

	// Output:
	// &""
	// &"abc"
}

type ofPointerTC[T any] struct {
	value T
	test.Control
}

func (tc ofPointerTC[T]) Test(t *testing.T) {
	opt := OfPointer(tc.value)
	value, present := opt.Get()
	assert.NotNil(t, value, "unexpected nil value")
	assert.Equal(t, tc.value, *value, "unexpected value")
	assert.True(t, present, "expected presence")
}

func TestOfPointer(t *testing.T) {
	test.RunCases(t, test.Cases{
		// Test cases for documented examples
		"with zero int": ofPointerTC[int]{
			value: 0,
		},
		"with non-zero int": ofPointerTC[int]{
			value: 123,
		},
		"with zero string": ofPointerTC[string]{
			value: "",
		},
		"with non-zero string": ofPointerTC[string]{
			value: "abc",
		},
		// Other test cases...
	})
}

func BenchmarkOfZeroable(b *testing.B) {
	for i := 0; i < b.N; i++ {
		OfZeroable(123)
	}
}

func ExampleOfZeroable_int() {
	example.Print(OfZeroable(0))
	example.Print(OfZeroable(123))

	// Output:
	// <empty>
	// 123
}

func ExampleOfZeroable_int_pointer() {
	example.Print(OfZeroable((*int)(nil)))
	example.Print(OfZeroable(ptrs.ZeroInt()))
	example.Print(OfZeroable(ptrs.Int(123)))

	// Output:
	// <empty>
	// &0
	// &123
}

func ExampleOfZeroable_string() {
	example.Print(OfZeroable(""))
	example.Print(OfZeroable("abc"))

	// Output:
	// <empty>
	// "abc"
}

func ExampleOfZeroable_string_pointer() {
	example.Print(OfZeroable((*string)(nil)))
	example.Print(OfZeroable(ptrs.ZeroString()))
	example.Print(OfZeroable(ptrs.String("abc")))

	// Output:
	// <empty>
	// &""
	// &"abc"
}

type ofZeroableTC[T any] struct {
	value         T
	expectPresent bool
	test.Control
}

func (tc ofZeroableTC[T]) Test(t *testing.T) {
	opt := OfZeroable(tc.value)
	value, present := opt.Get()
	assert.Equal(t, tc.value, value, "unexpected value")
	assert.Equal(t, tc.expectPresent, present, "unexpected value presence")
}

func TestOfZeroable(t *testing.T) {
	test.RunCases(t, test.Cases{
		// Test cases for documented examples
		"with zero int": ofZeroableTC[int]{
			value:         0,
			expectPresent: false,
		},
		"with non-zero int": ofZeroableTC[int]{
			value:         123,
			expectPresent: true,
		},
		"with nil int pointer": ofZeroableTC[*int]{
			value:         nil,
			expectPresent: false,
		},
		"with zero int pointer": ofZeroableTC[*int]{
			value:         ptrs.ZeroInt(),
			expectPresent: true,
		},
		"with non-zero int pointer": ofZeroableTC[*int]{
			value:         ptrs.Int(123),
			expectPresent: true,
		},
		"with zero string": ofZeroableTC[string]{
			value:         "",
			expectPresent: false,
		},
		"with non-zero string": ofZeroableTC[string]{
			value:         "abc",
			expectPresent: true,
		},
		"with nil string pointer": ofZeroableTC[*string]{
			value:         nil,
			expectPresent: false,
		},
		"with zero string pointer": ofZeroableTC[*string]{
			value:         ptrs.ZeroString(),
			expectPresent: true,
		},
		"with non-zero string pointer": ofZeroableTC[*string]{
			value:         ptrs.String("abc"),
			expectPresent: true,
		},
		// Other test cases...
	})
}

func BenchmarkRequireAny(b *testing.B) {
	opts := []Optional[int]{Empty[int](), Of(0), Of(123)}
	for i := 0; i < b.N; i++ {
		_ = RequireAny(opts...)
	}
}

func ExampleRequireAny_int() {
	example.PrintValues(RequireAny(Empty[int](), Of(0), Of(123)))

	// Output: [0 123]
}

func ExampleRequireAny_panic() {
	defer func() {
		fmt.Println(recover())
	}()

	RequireAny(Empty[int]())

	// Output: go-optional: value not present
}

func ExampleRequireAny_string() {
	example.PrintValues(RequireAny(Empty[string](), Of(""), Of("abc")))

	// Output: ["" "abc"]
}

type requireAnyTC[T any] struct {
	opts         []Optional[T]
	expectPanic  bool
	expectValues []T
	test.Control
}

func (tc requireAnyTC[T]) Test(t *testing.T) {
	if tc.expectPanic {
		assert.Panics(t, func() {
			RequireAny(tc.opts...)
		}, "expected panic")
	} else {
		var values []T
		assert.NotPanics(t, func() {
			values = RequireAny(tc.opts...)
		}, "unexpected panic")
		assert.Equal(t, tc.expectValues, values, "unexpected values")
	}
}

func TestRequireAny(t *testing.T) {
	test.RunCases(t, test.Cases{
		// Test cases for documented examples
		"with no int Optionals": requireAnyTC[int]{
			expectPanic: true,
		},
		"with empty int Optional": requireAnyTC[int]{
			opts:        []Optional[int]{Empty[int]()},
			expectPanic: true,
		},
		"with an empty int Optional and two non-empty int Optionals": requireAnyTC[int]{
			opts: []Optional[int]{
				Empty[int](),
				Of(0),
				Of(123),
			},
			expectValues: []int{0, 123},
		},
		"with no string Optionals": requireAnyTC[string]{
			expectPanic: true,
		},
		"with empty string Optional": requireAnyTC[string]{
			opts:        []Optional[string]{Empty[string]()},
			expectPanic: true,
		},
		"with an empty string Optional and two non-empty string Optionals": requireAnyTC[string]{
			opts: []Optional[string]{
				Empty[string](),
				Of("abc"),
				Of(""),
			},
			expectValues: []string{"abc", ""},
		},
		// Other test cases...
	})
}

func BenchmarkTryFlatMap(b *testing.B) {
	toString := func(value int) (Optional[string], error) {
		if value == 0 {
			return Empty[string](), nil
		}
		return Of(strconv.FormatInt(int64(value), 10)), nil
	}
	opt := Of(123)
	for i := 0; i < b.N; i++ {
		_, _ = TryFlatMap(opt, toString)
	}
}

func ExampleTryFlatMap_int() {
	mapper := func(value int) (Optional[string], error) {
		if value == 0 {
			return Empty[string](), nil
		}
		return Of(strconv.FormatInt(int64(value), 10)), nil
	}

	example.PrintTry(TryFlatMap(Empty[int](), mapper))
	example.PrintTry(TryFlatMap(Of(0), mapper))
	example.PrintTry(TryFlatMap(Of(123), mapper))

	// Output:
	// <empty> <nil>
	// <empty> <nil>
	// "123" <nil>
}

func ExampleTryFlatMap_string() {
	mapper := func(value string) (Optional[int], error) {
		if value == "" {
			return Empty[int](), nil
		}
		i, err := strconv.ParseInt(value, 10, 0)
		if err != nil {
			return Empty[int](), err
		}
		return OfZeroable(int(i)), nil
	}

	example.PrintTry(TryFlatMap(Empty[string](), mapper))
	example.PrintTry(TryFlatMap(Of(""), mapper))
	example.PrintTry(TryFlatMap(Of("0"), mapper))
	example.PrintTry(TryFlatMap(Of("123"), mapper))
	example.PrintTry(TryFlatMap(Of("abc"), mapper))

	// Output:
	// <empty> <nil>
	// <empty> <nil>
	// <empty> <nil>
	// 123 <nil>
	// <empty> "strconv.ParseInt: parsing \"abc\": invalid syntax"
}

type tryFatMapTC[T, M any] struct {
	opt           Optional[T]
	fn            func(value T) (Optional[M], error)
	expectError   bool
	expectPresent bool
	expectValue   M
	test.Control
}

func (tc tryFatMapTC[T, M]) Test(t *testing.T) {
	opt, err := TryFlatMap(tc.opt, tc.fn)
	value, present := opt.Get()
	assert.Equalf(t, tc.expectError, err != nil, "unexpected error: %v", err)
	assert.Equal(t, tc.expectValue, value, "unexpected value")
	assert.Equal(t, tc.expectPresent, present, "unexpected value presence")
}

func TestTryFlatMap(t *testing.T) {
	toInt := func(value string) (Optional[int], error) {
		if value == "" {
			return Empty[int](), nil
		}
		i, err := strconv.ParseInt(value, 10, 0)
		if err != nil {
			return Empty[int](), err
		}
		return OfZeroable(int(i)), nil
	}
	toString := func(value int) (Optional[string], error) {
		if value == 0 {
			return Empty[string](), nil
		}
		return Of(strconv.FormatInt(int64(value), 10)), nil
	}

	test.RunCases(t, test.Cases{
		// Test cases for documented examples
		"with empty int Optional": tryFatMapTC[int, string]{
			opt:           Empty[int](),
			fn:            toString,
			expectPresent: false,
		},
		"with non-empty int Optional with zero value": tryFatMapTC[int, string]{
			opt:           Of(0),
			fn:            toString,
			expectPresent: false,
		},
		"with non-empty int Optional with non-zero value": tryFatMapTC[int, string]{
			opt:           Of(123),
			fn:            toString,
			expectPresent: true,
			expectValue:   "123",
		},
		"with empty string Optional": tryFatMapTC[string, int]{
			opt:           Empty[string](),
			fn:            toInt,
			expectPresent: false,
		},
		"with non-empty string Optional with zero value": tryFatMapTC[string, int]{
			opt:           Of(""),
			fn:            toInt,
			expectPresent: false,
		},
		"with non-empty string Optional with zero-representing value": tryFatMapTC[string, int]{
			opt:           Of("0"),
			fn:            toInt,
			expectPresent: false,
		},
		"with non-empty string Optional with non-zero-representing value": tryFatMapTC[string, int]{
			opt:           Of("123"),
			fn:            toInt,
			expectPresent: true,
			expectValue:   123,
		},
		"with non-empty string Optional with erroneous value": tryFatMapTC[string, int]{
			opt:         Of("abc"),
			fn:          toInt,
			expectError: true,
		},
		// Other test cases...
	})
}

func BenchmarkTryMap(b *testing.B) {
	toString := func(value int) (string, error) {
		return strconv.FormatInt(int64(value), 10), nil
	}
	opt := Of(123)
	for i := 0; i < b.N; i++ {
		_, _ = TryMap(opt, toString)
	}
}

func ExampleTryMap_int() {
	mapper := func(value int) (string, error) {
		return strconv.FormatInt(int64(value), 10), nil
	}

	example.PrintTry(TryMap(Empty[int](), mapper))
	example.PrintTry(TryMap(Of(0), mapper))
	example.PrintTry(TryMap(Of(123), mapper))

	// Output:
	// <empty> <nil>
	// "0" <nil>
	// "123" <nil>
}

func ExampleTryMap_string() {
	mapper := func(value string) (int, error) {
		i, err := strconv.ParseInt(value, 10, 0)
		return int(i), err
	}

	example.PrintTry(TryMap(Empty[string](), mapper))
	example.PrintTry(TryMap(Of("0"), mapper))
	example.PrintTry(TryMap(Of("123"), mapper))
	example.PrintTry(TryMap(Of("abc"), mapper))

	// Output:
	// <empty> <nil>
	// 0 <nil>
	// 123 <nil>
	// <empty> "strconv.ParseInt: parsing \"abc\": invalid syntax"
}

type tryMapTC[T, M any] struct {
	opt           Optional[T]
	fn            func(value T) (M, error)
	expectError   bool
	expectPresent bool
	expectValue   M
	test.Control
}

func (tc tryMapTC[T, M]) Test(t *testing.T) {
	opt, err := TryMap(tc.opt, tc.fn)
	value, present := opt.Get()
	assert.Equalf(t, tc.expectError, err != nil, "unexpected error: %v", err)
	assert.Equal(t, tc.expectValue, value, "unexpected value")
	assert.Equal(t, tc.expectPresent, present, "unexpected value presence")
}

func TestTryMap(t *testing.T) {
	toInt := func(value string) (int, error) {
		i, err := strconv.ParseInt(value, 10, 0)
		return int(i), err
	}
	toString := func(value int) (string, error) {
		return strconv.FormatInt(int64(value), 10), nil
	}

	test.RunCases(t, test.Cases{
		// Test cases for documented examples
		"with empty int Optional": tryMapTC[int, string]{
			opt:           Empty[int](),
			fn:            toString,
			expectPresent: false,
		},
		"with non-empty int Optional with zero value": tryMapTC[int, string]{
			opt:           Of(0),
			fn:            toString,
			expectPresent: true,
			expectValue:   "0",
		},
		"with non-empty int Optional with non-zero value": tryMapTC[int, string]{
			opt:           Of(123),
			fn:            toString,
			expectPresent: true,
			expectValue:   "123",
		},
		"with empty string Optional": tryMapTC[string, int]{
			opt:           Empty[string](),
			fn:            toInt,
			expectPresent: false,
		},
		"with non-empty string Optional with zero-representing value": tryMapTC[string, int]{
			opt:           Of("0"),
			fn:            toInt,
			expectPresent: true,
			expectValue:   0,
		},
		"with non-empty string Optional with non-zero-representing value": tryMapTC[string, int]{
			opt:           Of("123"),
			fn:            toInt,
			expectPresent: true,
			expectValue:   123,
		},
		"with non-empty string Optional with erroneous value": tryMapTC[string, int]{
			opt:         Of("abc"),
			fn:          toInt,
			expectError: true,
		},
		// Other test cases...
	})
}
