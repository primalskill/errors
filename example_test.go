package errors_test

import (
	"fmt"
	"github.com/primalskill/errors"
)

func ExampleE() {
	err := errors.E("this is an error")
	fmt.Println(err.Error())

	// Output: this is an error
}

func ExampleE_meta() {

	// err will carry a Meta map.
	err := errors.E("this is an error with meta", errors.WithMeta("myKey", "value_testing"))

	var e *errors.Error
	errors.As(err, &e)

	fmt.Printf("%+v", e.Meta)

	// Output: [myKey:value_testing]
}

func ExampleUnwrap() {
	err1 := errors.E("error 1")
	err2 := errors.E("error 2", err1)

	e := errors.Unwrap(err2)
	fmt.Printf("%+v", e)

	// Output: error 1
}

func ExampleUnwrapAll() {
	err1 := errors.E("error1", errors.WithMeta("err1Key", "err1 value"))
	err2 := errors.E("error2", err1, errors.WithMeta("err2Key", "err2 value"))

	errs := errors.UnwrapAll(err2)
	fmt.Printf("%+v", errs)

	// Output: [error2 error1]
}

func ExampleErrorFull() {
	err := errors.E("my error", errors.WithMeta("key", "value"))
	fmt.Printf("%s", ErrorFull(err))
}

func ExampleError_ErrorFull() {
	err := errors.E("my error", errors.WithMeta("key", "value"))

	var e *errors.Error
	errors.As(err, &e)

	fmt.Printf("%s", e.ErrorFull())
}

func ExampleWithMeta() {
	// Valid
	mValid := errors.WithMeta("key1", 158, "key2", "some value", "anotherKey", true)

	fmt.Printf("%#v", mValid)

	// Invalid Outputs
	// errors.WithMeta(15) <-- results in compile error
	// errors.WithMeta("noValueKey") <-- returns empty Meta map because there is no value added to noValueKey
	// errors.WithMeta("key1", "value 1", 16, "key3", "key4", "some value") <-- skips 16 and key3 pairs because 16 is int and not string.
	// errors.WithMeta("key1", "val1", 10, "val2", "key3") <-- skips 10 and val2 pairs, output: [key1:val1 key3:]

	// Output: errors.Meta{"anotherKey":true, "key1":158, "key2":"some value"}
}

func ExampleWithMeta_emptyValue() {
	mValid := errors.WithMeta("key1", "value1", "key2")
	fmt.Printf("%#v", mValid)
}

func ExampleMeta_Merge() {
	m := errors.WithMeta("key1", "val1")
	m = m.Merge("key2", "val2")

	fmt.Printf("%#v", m)

	// Output: errors.Meta{"key1":"val1", "key2":"val2"}
}

func ExampleStack_String() {
	err := errors.E("my error")

	var e *errors.Error
	errors.As(err, &e)

	fmt.Printf("%s", e.Stack.String())
}

func ExampleStack_PrettyPrint() {
	err := errors.E("my error")

	var e *errors.Error
	errors.As(err, &e)

	fmt.Printf("%s", e.Stack.PrettyPrint())
}
