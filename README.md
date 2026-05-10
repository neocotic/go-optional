# go-optional

[![Go Reference](https://img.shields.io/badge/go.dev-reference-007d9c?style=for-the-badge&logo=go&logoColor=white)](https://pkg.go.dev/github.com/airmrcr/go-optional)
[![Build Status](https://img.shields.io/github/actions/workflow/status/airmrcr/go-optional/ci.yml?style=for-the-badge)](https://github.com/airmrcr/go-optional/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/airmrcr/go-optional?style=for-the-badge)](https://github.com/airmrcr/go-optional)
[![License](https://img.shields.io/github/license/airmrcr/go-optional?style=for-the-badge)](https://github.com/airmrcr/go-optional/blob/main/LICENSE.md)

Easy-to-use generic optional values for Go (golang).

Inspired by Java's `Optional` class, it enables the ability to differentiate a value that has its zero value due to not
being set from having a zero value that was explicitly set, drastically reducing the need for pointers in a lot of use
cases.

## Installation

Install using [go install](https://go.dev/ref/mod#go-install):

``` sh
go get github.com/airmrcr/go-optional
```

Then import the package into your own code:

``` go
import "github.com/airmrcr/go-optional"
```

## Documentation

Documentation is available on [pkg.go.dev](https://pkg.go.dev/github.com/airmrcr/go-optional#section-documentation). It
contains an overview and reference.

### Example

Basic usage:

``` go
opt := optional.Of(0)
opt.IsPresent()  // true
opt.Get()        // 0, true

opt := optional.Of(123)
opt.IsPresent()  // true
opt.Get()        // 123, true

var external []int
opt := optional.OfNillable(external)
opt.IsPresent()      // false
opt.Get()            // nil, false
opt.OrElse([]int{})  // []int{}

opt := optional.Empty[string]()
opt.IsPresent()  // false
opt.Get()        // "", false
```

Optionals are also intended to be used as struct fields and support JSON, XML, and YAML marshaling and unmarshaling
out-of-the-box. That said; the very nature of optionals leans more towards unmarshaling.

``` go
type Example struct {
    Number Optional[int]    `json:"number"`
    Text   Optional[string] `json:"text"`
}

var example Example
err := json.Unmarshal([]byte(`{"text": "abc"}`, &example)
if err != nil {
    panic(err)
}
example.Number.Get()  // 0, false
example.Text.Get()    // "abc", true
```

There's a load of other functions and methods on `Optional` (along with its own `sort` sub-package) to explore, all with
documented examples.

## Issues

If you have any problems or would like to see changes currently in development you can do so
[here](https://github.com/airmrcr/go-optional/issues).

## Contributors

If you want to contribute, you're a legend! Information on how you can do so can be found in
[CONTRIBUTING.md](https://github.com/airmrcr/go-optional/blob/main/CONTRIBUTING.md). We want your suggestions and pull
requests!

A list of contributors can be found in [AUTHORS.md](https://github.com/airmrcr/go-optional/blob/main/AUTHORS.md).

## License

Copyright © 2026 airmrcr

See [LICENSE.md](https://github.com/airmrcr/go-optional/raw/main/LICENSE.md) for more information on our MIT license.
