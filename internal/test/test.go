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

// Package test provides helpers for testing the module.
package test

import "testing"

type (
	// Case is a test case that contains its own test logic.
	Case interface {
		// IsSkipped returns whether the test runner should skip the Case.
		IsSkipped() bool
		// Test runs the test logic for the Case.
		Test(t *testing.T)
	}

	// Cases is a mapping of test names to their Cases.
	Cases map[string]Case

	// Control contains fields used for controlling the execution of a Case.
	//
	// It's expected that struct implementations of Case will embed Control, granting greater focus on test logic.
	Control struct {
		// Skip is whether the Case should be skipped.
		Skip bool
	}
)

// IsSkipped returns whether the test runner should skip the Case.
func (c Control) IsSkipped() bool {
	return c.Skip
}

// RunCases runs all the provided test cases, applying controls as needed.
func RunCases(t *testing.T, cases Cases) {
	for name, c := range cases {
		t.Run(name, func(tt *testing.T) {
			if c.IsSkipped() {
				tt.Skip()
			}
			c.Test(tt)
		})
	}
}
