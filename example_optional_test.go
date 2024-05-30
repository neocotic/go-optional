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
	"context"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/neocotic/go-optional/internal/example"
	ptrs "github.com/neocotic/go-pointers"
	"gopkg.in/yaml.v3"
	"log"
	"strconv"
	"strings"
	"unicode"
)

var (
	ctx context.Context
	db  *sql.DB
)

func ExampleOptional_Equal_int() {
	fmt.Println(Empty[int]().Equal(Empty[int]()))
	fmt.Println(Empty[int]().Equal(Of(0)))
	fmt.Println(Of(0).Equal(Empty[int]()))
	fmt.Println(Of(0).Equal(Of(0)))
	fmt.Println(Of(0).Equal(Of(123)))
	fmt.Println(Of(123).Equal(Of(0)))
	fmt.Println(Of(123).Equal(Of(123)))
	fmt.Println(Of(123).Equal(Of(-123)))
	fmt.Println(Of(123).Equal(Empty[int]()))

	// Output:
	// true
	// false
	// false
	// true
	// false
	// false
	// true
	// false
	// false
}

func ExampleOptional_Equal_string() {
	fmt.Println(Empty[string]().Equal(Empty[string]()))
	fmt.Println(Empty[string]().Equal(Of("")))
	fmt.Println(Of("").Equal(Empty[string]()))
	fmt.Println(Of("").Equal(Of("")))
	fmt.Println(Of("").Equal(Of("abc")))
	fmt.Println(Of("abc").Equal(Of("")))
	fmt.Println(Of("abc").Equal(Of("abc")))
	fmt.Println(Of("abc").Equal(Of("ABC")))
	fmt.Println(Of("abc").Equal(Empty[string]()))

	// Output:
	// true
	// false
	// false
	// true
	// false
	// false
	// true
	// false
	// false
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

func ExampleOptional_Scan() {
	rows, err := db.QueryContext(ctx, "SELECT name, age FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	users := make(map[string]Optional[int])
	for rows.Next() {
		var (
			age  Optional[int]
			name string
		)
		if err = rows.Scan(&name, &age); err != nil {
			log.Fatal(err)
		}
		users[name] = age
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	log.Printf("user demographics: %s", users)
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
			log.Fatal(err)
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
			log.Fatal(err)
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
			log.Fatal(err)
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

func ExampleOptional_Value() {
	username := "alex"
	age := Of(30)
	result, err := db.ExecContext(ctx, "UPDATE users SET age = ? WHERE username = ?", age, username)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	if rows != 1 {
		log.Fatalf("expected to affect 1 row, affected %d", rows)
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

func ExampleEmpty_int() {
	example.Print(Empty[int]())

	// Output: <empty>
}

func ExampleEmpty_string() {
	example.Print(Empty[string]())

	// Output: <empty>
}

func ExampleEqual_int() {
	fmt.Println(Equal(Empty[int](), Empty[int]()))
	fmt.Println(Equal(Empty[int](), Of(0)))
	fmt.Println(Equal(Of(0), Empty[int]()))
	fmt.Println(Equal(Of(0), Of(0)))
	fmt.Println(Equal(Of(0), Of(123)))
	fmt.Println(Equal(Of(123), Of(0)))
	fmt.Println(Equal(Of(123), Of(123)))
	fmt.Println(Equal(Of(123), Of(-123)))
	fmt.Println(Equal(Of(123), Empty[int]()))

	// Output:
	// true
	// false
	// false
	// true
	// false
	// false
	// true
	// false
	// false
}

func ExampleEqual_mixed() {
	fmt.Println(Equal(Empty[any](), Empty[int]()))
	fmt.Println(Equal(Empty[any](), Of(0)))
	fmt.Println(Equal(Of[any](0), Of(0)))
	fmt.Println(Equal(Of[any](123), Of(123)))
	fmt.Println(Equal(Of(0), Of("0")))

	// Output:
	// true
	// false
	// true
	// true
	// false
}

func ExampleEqual_string() {
	fmt.Println(Equal(Empty[string](), Empty[string]()))
	fmt.Println(Equal(Empty[string](), Of("")))
	fmt.Println(Equal(Of(""), Empty[string]()))
	fmt.Println(Equal(Of(""), Of("")))
	fmt.Println(Equal(Of(""), Of("abc")))
	fmt.Println(Equal(Of("abc"), Of("")))
	fmt.Println(Equal(Of("abc"), Of("abc")))
	fmt.Println(Equal(Of("abc"), Of("ABC")))
	fmt.Println(Equal(Of("abc"), Empty[string]()))

	// Output:
	// true
	// false
	// false
	// true
	// false
	// false
	// true
	// false
	// false
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
			log.Fatal(err)
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
			log.Fatal(err)
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

func ExampleOf_int() {
	example.Print(Of(0))
	example.Print(Of(123))

	// Output:
	// 0
	// 123
}

func ExampleOf_intPointer() {
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

func ExampleOf_stringPointer() {
	example.Print(Of((*string)(nil)))
	example.Print(Of(ptrs.ZeroString()))
	example.Print(Of(ptrs.String("abc")))

	// Output:
	// <nil>
	// &""
	// &"abc"
}

func ExampleOfNillable_int() {
	example.Print(OfNillable(0))
	example.Print(OfNillable(123))

	// Output:
	// 0
	// 123
}

func ExampleOfNillable_intPointer() {
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

func ExampleOfNillable_stringPointer() {
	example.Print(OfNillable((*string)(nil)))
	example.Print(OfNillable(ptrs.ZeroString()))
	example.Print(OfNillable(ptrs.String("abc")))

	// Output:
	// <empty>
	// &""
	// &"abc"
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

func ExampleOfZeroable_int() {
	example.Print(OfZeroable(0))
	example.Print(OfZeroable(123))

	// Output:
	// <empty>
	// 123
}

func ExampleOfZeroable_intPointer() {
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

func ExampleOfZeroable_stringPointer() {
	example.Print(OfZeroable((*string)(nil)))
	example.Print(OfZeroable(ptrs.ZeroString()))
	example.Print(OfZeroable(ptrs.String("abc")))

	// Output:
	// <empty>
	// &""
	// &"abc"
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
