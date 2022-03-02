Error module used internally at Primal Skill, inspired by [upspin/errors](https://pkg.go.dev/github.com/palager/upspin/errors) and [pkg/errors](https://pkg.go.dev/github.com/pkg/errors).


## Install

```bash
go get -u github.com/primalskill/errors
```

## Basic Usage

```go
package main

import (
  "fmt"
  "github.com/primalskill/errors"
)

func main() {
  err := errors.E("this is an error")
  fmt.Printf("%s", err.Error())
}
```
