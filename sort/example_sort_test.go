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

package sort

import (
	"fmt"
	"github.com/neocotic/go-optional"
	"github.com/neocotic/go-optional/internal/example"
)

func ExampleAsc_int() {
	opts := []optional.Optional[int]{optional.Of(0), optional.Of(123), optional.Empty[int]()}
	Asc(opts)
	example.PrintSlice(opts)

	// Output: [<empty> 0 123]
}

func ExampleAsc_string() {
	opts := []optional.Optional[string]{optional.Of(""), optional.Of("abc"), optional.Empty[string]()}
	Asc(opts)
	example.PrintSlice(opts)

	// Output: [<empty> "" "abc"]
}

func ExampleDesc_int() {
	opts := []optional.Optional[int]{optional.Of(0), optional.Of(123), optional.Empty[int]()}
	Desc(opts)
	example.PrintSlice(opts)

	// Output: [123 0 <empty>]
}

func ExampleDesc_string() {
	opts := []optional.Optional[string]{optional.Of(""), optional.Of("abc"), optional.Empty[string]()}
	Desc(opts)
	example.PrintSlice(opts)

	// Output: ["abc" "" <empty>]
}

func ExampleIsAsc_int() {
	fmt.Println(IsAsc(([]optional.Optional[int])(nil)))
	fmt.Println(IsAsc([]optional.Optional[int]{optional.Of(0), optional.Of(123), optional.Empty[int]()}))
	fmt.Println(IsAsc([]optional.Optional[int]{optional.Empty[int](), optional.Of(0), optional.Of(123)}))

	// Output:
	// true
	// false
	// true
}

func ExampleIsAsc_string() {
	fmt.Println(IsAsc(([]optional.Optional[string])(nil)))
	fmt.Println(IsAsc([]optional.Optional[string]{optional.Of(""), optional.Of("abc"), optional.Empty[string]()}))
	fmt.Println(IsAsc([]optional.Optional[string]{optional.Empty[string](), optional.Of(""), optional.Of("abc")}))

	// Output:
	// true
	// false
	// true
}

func ExampleIsDesc_int() {
	fmt.Println(IsDesc(([]optional.Optional[int])(nil)))
	fmt.Println(IsDesc([]optional.Optional[int]{optional.Of(0), optional.Of(123), optional.Empty[int]()}))
	fmt.Println(IsDesc([]optional.Optional[int]{optional.Of(123), optional.Of(0), optional.Empty[int]()}))

	// Output:
	// true
	// false
	// true
}

func ExampleIsDesc_string() {
	fmt.Println(IsDesc(([]optional.Optional[string])(nil)))
	fmt.Println(IsDesc([]optional.Optional[string]{optional.Of(""), optional.Of("abc"), optional.Empty[string]()}))
	fmt.Println(IsDesc([]optional.Optional[string]{optional.Of("abc"), optional.Of(""), optional.Empty[string]()}))

	// Output:
	// true
	// false
	// true
}
