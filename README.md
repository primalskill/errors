Error module used internally at Primal Skill, inspired by [upspin/errors](https://pkg.go.dev/github.com/palager/upspin/errors) and [pkg/errors](https://pkg.go.dev/github.com/pkg/errors).

## Docs

https://pkg.go.dev/github.com/primalskill/errors


## Install

Get the module:

```bash
go get -u github.com/primalskill/errors
```

Import it in your code:

```go
import "github.com/primalskill/errors"
```

## Running the Tests

Navigate to the module's directory and execute:

```bash
make test
```

## Example - Basic Usage

You can find more examples in the [docs](https://pkg.go.dev/github.com/primalskill/errors#pkg-examples).

```go
package main

import (
  "fmt"
  "github.com/primalskill/errors"
)

func main() {

  // Define an error
  err1 := errors.E(
    "some error", 
    errors.WithMeta(
      "metaKey1", "meta value", 
      "isAuth", true,
    ),
  )

  // Define another error and wrap err1
  err2 := errors.E(
    "wrapped error", 
    err1, 
    errors.WithMeta(
      "additionalMeta", 246,
    ),
  )

  fmt.Println(err2.Error()) // output: "wrapped error"

  // Convert stdlib error to errors.Error
  var ee *errors.Error
  errors.As(err2, &ee)
  fmt.Printf("%+v", ee.Meta) // output: [additionalMeta:246]

  // GetMeta helper func
  m, ok := errors.GetMeta(err1)
  fmt.Printf("\n%+v\n%+v\n", ok, m) // output: true \n [metaKey1:meta value isAuth:true]

  // Unwrap err2 to get err1
  uErr := errors.Unwrap(err2)
  fmt.Println(uErr.Error()) // output: "some error"
}
```

## Example - Mirror an Existing Error

```go
package main

import (
  "fmt"
  "github.com/primalskill/errors"
)

func main() {

  // Define an error
  err1 := errors.E("this is an error", errors.WithMeta("key1", "val1"))

  // Preload or "mirror" err1 in err2 with a few rules:
  // - err2 preloads err1 message
  // - err2 preloads err1 Meta key/value pairs
  // - if additional metas are defined on err2 it will merge it to the others
  // - err2 overwrites err1 source location to correctly show the location where err2 was executed
  err2 := errors.M(err1, errors.WithMeta("key2", "val2"))

  fmt.Printf("%+v", errors.PrettyPrint(err2)) // PrettyPrint should only be used in development to have a nicer output
  /* 
  output:

  this is an error
     |- Source : /goprograms/errors/example_test.go:20
    |- Meta :
      |- key1 : val1
      |- key2 : val2
  */
}
```

