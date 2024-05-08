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

// Package sort provides basic support for sorting slices of optional.Optional.
package sort

import (
	"cmp"
	"github.com/neocotic/go-optional"
	"sort"
)

// Asc sorts the given slice using optional.Compare in ascending order.
//
// For example;
//
//	opts := []Optional[int]{Of(0), Of(123), Empty[int]()}
//	Asc(opts)
//	opts       // []Optional[int]{Empty[int](), Of(0), Of(123)}
//
//	opts := []Optional[string]{Of(""), Of("abc"), Empty[string]()}
//	Asc(opts)
//	opts       // []Optional[string]{Empty[string](), Of(""), Of("abc")}
func Asc[T cmp.Ordered](opts []optional.Optional[T]) {
	if len(opts) == 0 {
		return
	}
	sort.Slice(opts, func(i, j int) bool {
		return optional.Compare(opts[i], opts[j]) < 0
	})
}

// Desc sorts the given slice using optional.Compare in descending order.
//
// For example;
//
//	opts := []Optional[int]{Of(0), Of(123), Empty[int]()}
//	Desc(opts)
//	opts        // []Optional[int]{Of(123), Of(0), Empty[int]()}
//
//	opts := []Optional[string]{Of(""), Of("abc"), Empty[string]()}
//	Desc(opts)
//	opts        // []Optional[string]{Of("abc"), Of(""), Empty[string]()}
func Desc[T cmp.Ordered](opts []optional.Optional[T]) {
	if len(opts) == 0 {
		return
	}
	sort.Slice(opts, func(i, j int) bool {
		return optional.Compare(opts[i], opts[j]) > 0
	})
}

// IsAsc returns whether the given slice is sorted using optional.Compare in ascending order.
//
// For example;
//
//	IsAsc(([]Optional[int])(nil))                         // true
//	IsAsc([]Optional[int]{Of(0), Of(123), Empty[int]()})  // false
//	IsAsc([]Optional[int]{Empty[int](), Of(0), Of(123)})  // true
//
//	IsAsc(([]Optional[string])(nil))                               // true
//	IsAsc([]Optional[string]{Of(""), Of("abc"), Empty[string]()})  // false
//	IsAsc([]Optional[string]{Empty[string](), Of(""), Of("abc")})  // true
func IsAsc[T cmp.Ordered](opts []optional.Optional[T]) bool {
	if len(opts) == 0 {
		return true
	}
	return sort.SliceIsSorted(opts, func(i, j int) bool {
		return optional.Compare(opts[i], opts[j]) < 0
	})
}

// IsDesc returns whether the given slice is sorted using optional.Compare in descending order.
//
// For example;
//
//	IsDesc(([]Optional[int])(nil))                         // true
//	IsDesc([]Optional[int]{Of(0), Of(123), Empty[int]()})  // false
//	IsDesc([]Optional[int]{Of(123), Of(0), Empty[int]()})  // true
//
//	IsDesc(([]Optional[string])(nil))                               // true
//	IsDesc([]Optional[string]{Of(""), Of("abc"), Empty[string]()})  // false
//	IsDesc([]Optional[string]{Of("abc"), Of(""), Empty[string]()})  // true
func IsDesc[T cmp.Ordered](opts []optional.Optional[T]) bool {
	if len(opts) == 0 {
		return true
	}
	return sort.SliceIsSorted(opts, func(i, j int) bool {
		return optional.Compare(opts[i], opts[j]) > 0
	})
}
