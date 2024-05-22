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

// Package example provides helpers for generating concise examples for the module.
package example

import (
	"fmt"
	"reflect"
	"strings"
)

// optional defines the minimal interface for optional.Optional that will satisfy the example printing whilst avoiding a
// circular dependency.
type optional[T any] interface {
	Get() (T, bool)
}

// Print prints the formatted value of the given optional.Optional, if present, otherwise its string representation.
func Print[T any, O optional[T]](opt O) {
	if value, present := opt.Get(); present {
		printValue(value)
		fmt.Println()
	} else {
		fmt.Println(opt)
	}
}

// PrintGet prints the formatted value provided and whether its present.
//
// Only intended for printing the result of calling optional.Optional.Get.
func PrintGet(value any, present bool) {
	printValue(value)
	fmt.Printf(" %v\n", present)
}

// PrintMarshalled prints the output marshalling an optional.Optional via an encoder as well as err.
func PrintMarshalled(data []byte, err error) {
	fmt.Print(strings.TrimSpace(string(data)))
	fmt.Print(" ")
	printError(err)
	fmt.Println()
}

// PrintSlice prints the formatted value of each given optional.Optional, if present, otherwise their string
// representation, as a slice.
func PrintSlice[T any, O optional[T]](opts []O) {
	fmt.Print("[")
	for i, opt := range opts {
		if i > 0 {
			fmt.Print(" ")
		}
		if value, present := opt.Get(); present {
			printValue(value)
		} else {
			fmt.Print(opt)
		}
	}
	fmt.Print("]")
}

// PrintTry prints the formatted value of the given optional.Optional, if present, otherwise its string representation.
// The err provided is also printed.
func PrintTry[T any, O optional[T]](opt O, err error) {
	if value, present := opt.Get(); present {
		printValue(value)
	} else {
		fmt.Print(opt)
	}
	fmt.Print(" ")
	printError(err)
	fmt.Println()
}

// PrintTryValue prints the formatted value provided as well as err.
func PrintTryValue[T any](value T, err error) {
	printValue(value)
	fmt.Print(" ")
	printError(err)
	fmt.Println()
}

// PrintValue prints the formatted value provided.
func PrintValue[T any](value T) {
	printValue(value)
	fmt.Println()
}

// PrintValues prints the formatted values provided.
func PrintValues(values any) {
	rv := reflect.ValueOf(values)
	if rv.Kind() != reflect.Slice {
		panic(fmt.Errorf("printValues argument must be a slice: %s", rv.Kind()))
	}
	if rv.Len() == 0 {
		fmt.Println(values)
		return
	}
	if rv.Type().Elem().Kind() == reflect.String {
		fmt.Printf("%q\n", values)
	} else {
		fmt.Println(values)
	}
}

// printError formats and prints the error provided, quoting err if not nil.
func printError(err error) {
	if err == nil {
		fmt.Print(err)
	} else {
		fmt.Printf("%q", err)
	}
}

// printValue formats and prints the value provided, quoting value if a string and prefixing with an ampersand if a
// pointer.
func printValue(value any) {
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.String:
		fmt.Printf("%q", value)
	case reflect.Pointer:
		if rv.IsNil() {
			fmt.Print(value)
		} else if ert := rv.Type().Elem(); ert.Kind() == reflect.String {
			fmt.Printf("&%q", reflect.Indirect(rv).Interface())
		} else {
			fmt.Printf("&%v", reflect.Indirect(rv).Interface())
		}
	default:
		fmt.Print(value)
	}
}
