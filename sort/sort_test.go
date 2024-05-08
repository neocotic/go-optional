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
	"cmp"
	"github.com/neocotic/go-optional"
	"github.com/neocotic/go-optional/internal/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Benchmark_Asc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Asc([]optional.Optional[int]{
			optional.Empty[int](),
			optional.Of(123),
			optional.Of(12),
			optional.Of(1),
			optional.Of(0),
			optional.Of(-1),
			optional.Of(-12),
			optional.Of(-123),
			optional.Empty[int](),
		})
	}
}

type ascTC[T cmp.Ordered] struct {
	opts   []optional.Optional[T]
	expect []optional.Optional[T]
	test.Control
}

func (tc ascTC[T]) Test(t *testing.T) {
	Asc(tc.opts)
	assert.Equal(t, tc.expect, tc.opts, "unexpected optionals")
}

func Test_Asc(t *testing.T) {
	test.RunCases(t, test.Cases{
		// Test cases for documented examples
		"with int Optionals": ascTC[int]{
			opts: []optional.Optional[int]{
				optional.Of(0),
				optional.Of(123),
				optional.Empty[int](),
			},
			expect: []optional.Optional[int]{
				optional.Empty[int](),
				optional.Of(0),
				optional.Of(123),
			},
		},
		"with string Optionals": ascTC[string]{
			opts: []optional.Optional[string]{
				optional.Of(""),
				optional.Of("abc"),
				optional.Empty[string](),
			},
			expect: []optional.Optional[string]{
				optional.Empty[string](),
				optional.Of(""),
				optional.Of("abc"),
			},
		},
		// Other test cases...
	})
}

func Benchmark_Desc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Desc([]optional.Optional[int]{
			optional.Empty[int](),
			optional.Of(-123),
			optional.Of(-12),
			optional.Of(-1),
			optional.Of(0),
			optional.Of(1),
			optional.Of(12),
			optional.Of(123),
			optional.Empty[int](),
		})
	}
}

type descTC[T cmp.Ordered] struct {
	opts   []optional.Optional[T]
	expect []optional.Optional[T]
	test.Control
}

func (tc descTC[T]) Test(t *testing.T) {
	Desc(tc.opts)
	assert.Equal(t, tc.expect, tc.opts, "unexpected optionals")
}

func Test_Desc(t *testing.T) {
	test.RunCases(t, test.Cases{
		// Test cases for documented examples
		"with int Optionals": descTC[int]{
			opts: []optional.Optional[int]{
				optional.Of(0),
				optional.Of(123),
				optional.Empty[int](),
			},
			expect: []optional.Optional[int]{
				optional.Of(123),
				optional.Of(0),
				optional.Empty[int](),
			},
		},
		"with string Optionals": descTC[string]{
			opts: []optional.Optional[string]{
				optional.Of(""),
				optional.Of("abc"),
				optional.Empty[string](),
			},
			expect: []optional.Optional[string]{
				optional.Of("abc"),
				optional.Of(""),
				optional.Empty[string](),
			},
		},
		// Other test cases...
	})
}

func Benchmark_IsAsc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = IsAsc([]optional.Optional[int]{
			optional.Empty[int](),
			optional.Empty[int](),
			optional.Of(-123),
			optional.Of(-12),
			optional.Of(-1),
			optional.Of(0),
			optional.Of(1),
			optional.Of(12),
			optional.Of(123),
		})
	}
}

type isAscTC[T cmp.Ordered] struct {
	opts   []optional.Optional[T]
	expect bool
	test.Control
}

func (tc isAscTC[T]) Test(t *testing.T) {
	actual := IsAsc(tc.opts)
	assert.Equalf(t, tc.expect, actual, "unexpected sorting: %v", tc.opts)
}

func Test_IsAsc(t *testing.T) {
	test.RunCases(t, test.Cases{
		// Test cases for documented examples
		"with no int Optionals": isAscTC[int]{
			expect: true,
		},
		"with int Optionals not sorted in ascending order": isAscTC[int]{
			opts: []optional.Optional[int]{
				optional.Of(0),
				optional.Of(123),
				optional.Empty[int](),
			},
			expect: false,
		},
		"with int Optionals sorted in ascending order": isAscTC[int]{
			opts: []optional.Optional[int]{
				optional.Empty[int](),
				optional.Of(0),
				optional.Of(123),
			},
			expect: true,
		},
		"with no string Optionals": isAscTC[string]{
			expect: true,
		},
		"with string Optionals not sorted in ascending order": isAscTC[string]{
			opts: []optional.Optional[string]{
				optional.Of(""),
				optional.Of("abc"),
				optional.Empty[string](),
			},
			expect: false,
		},
		"with string Optionals sorted in ascending order": isAscTC[string]{
			opts: []optional.Optional[string]{
				optional.Empty[string](),
				optional.Of(""),
				optional.Of("abc"),
			},
			expect: true,
		},
		// Other test cases...
	})
}

func Benchmark_IsDesc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = IsDesc([]optional.Optional[int]{
			optional.Of(123),
			optional.Of(12),
			optional.Of(1),
			optional.Of(0),
			optional.Of(-1),
			optional.Of(-12),
			optional.Of(-123),
			optional.Empty[int](),
			optional.Empty[int](),
		})
	}
}

type isDescTC[T cmp.Ordered] struct {
	opts   []optional.Optional[T]
	expect bool
	test.Control
}

func (tc isDescTC[T]) Test(t *testing.T) {
	actual := IsDesc(tc.opts)
	assert.Equalf(t, tc.expect, actual, "unexpected sorting: %v", tc.opts)
}

func Test_IsDesc(t *testing.T) {
	test.RunCases(t, test.Cases{
		// Test cases for documented examples
		"with no int Optionals": isDescTC[int]{
			expect: true,
		},
		"with int Optionals not sorted in descending order": isDescTC[int]{
			opts: []optional.Optional[int]{
				optional.Of(0),
				optional.Of(123),
				optional.Empty[int](),
			},
			expect: false,
		},
		"with int Optionals sorted in descending order": isDescTC[int]{
			opts: []optional.Optional[int]{
				optional.Of(123),
				optional.Of(0),
				optional.Empty[int](),
			},
			expect: true,
		},
		"with no string Optionals": isDescTC[string]{
			expect: true,
		},
		"with string Optionals not sorted in descending order": isDescTC[string]{
			opts: []optional.Optional[string]{
				optional.Of(""),
				optional.Of("abc"),
				optional.Empty[string](),
			},
			expect: false,
		},
		"with string Optionals sorted in descending order": isDescTC[string]{
			opts: []optional.Optional[string]{
				optional.Of("abc"),
				optional.Of(""),
				optional.Empty[string](),
			},
			expect: true,
		},
		// Other test cases...
	})
}
