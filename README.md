Error module used internally at Primal Skill, inspired by [upspin/errors](https://pkg.go.dev/github.com/palager/upspin/errors) and [pkg/errors](https://pkg.go.dev/github.com/pkg/errors).


## Install

Get the module:

```bash
go get -u github.com/primalskill/errors
```

Import it in your code:

```go
import "github.com/primalskill/errors"
```

## Basic Usage

```go
package main

import (
  "fmt"
  "github.com/primalskill/errors"
)

func main() {
  err1 := errors.E("this is an error", errors.WithMeta("metaKey1", "meta value 1", "isAuth", true))
  err2 := errors.E("embedded error", err1, errors.WithMeta("additionalMeta", 246))

  fmt.Println(err2.Error()) // outputs the error message: embedded error

  // Convert stdlib error to errors.Error
  var ee *errors.Error
  errors.As(err2, &ee)

  fmt.Printf("%+v", ee.Meta) // outputs err2 Meta

  m := errors.GetMeta(err2) // get the Meta with a helper func
  fmt.Printf("%+v", m) // outputs the Meta attached to err2

  uErr := errors.Unwrap(err2) // unwraps err2 to get err1
  fmt.Println(uErr.Error()) // outputs: this is an error
}
```
