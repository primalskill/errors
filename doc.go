// Package errors implements custom error handling by decorating the standard Error interface with additional meta data.
//
// Installation
//
// To install the module, you will need Go 1.13 or a newer version.
//
// Getting the module:
//
//  $ go get -u github.com/primalskill/errors
//
// Import it in your code:
//
//  import "github.com/primalskill/errors"
//
// Overview
//
// Always use the E() function to create new errors, to which you can add additional data.
//
//  package main
//
//  import (
//    "fmt"
//    "github.com/primalskill/errors"
//  )
//
//  func main() {
//    err1 := errors.E("this is an error", errors.WithMeta("metaKey1", "meta value 1", "isAuth", true))
//    err2 := errors.E("embedded error", err1, errors.WithMeta("additionalMeta", 246))
//
//    fmt.Println(err2.Error()) // outputs the error message: embedded error
//
//    // Convert stdlib error to errors.Error
//    var ee *errors.Error
//    errors.As(err2, &ee)
//
//    fmt.Printf("%+v", ee.Meta) // outputs err2 Meta
//
//    m := errors.GetMeta(err2) // get the Meta with a helper func
//    fmt.Printf("%+v", m) // outputs the Meta attached to err2
//
//    uErr := errors.Unwrap(err2) // unwraps err2 to get err1
//    fmt.Println(uErr.Error()) // outputs: this is an error
//  }
//
package errors
