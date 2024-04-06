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

func ExamplePrettyPrint() {
	err := errors.E("my error", errors.WithMeta("key", "value"))
	fmt.Printf("%s", errors.PrettyPrint(err))
}

func ExampleError_PrettyPrint() {
	err := errors.E("my error", errors.WithMeta("key", "value"))

	var e *errors.Error
	errors.As(err, &e)

	fmt.Printf("%s", e.PrettyPrint())
}

func ExampleWithMeta() {
	// Valid
	mValid := errors.WithMeta("key1", 158, "key2", "some value", "anotherKey", true)
	fmt.Printf("%#v", mValid)

	// no value defined it will add !BADVALUE to noValueKey
	mBadValue := errors.WithMeta("noValueKey")
	fmt.Printf("\n%#v", mBadValue)

	// 16 is not a string for a key it will replace it with !BADKEY2
	mBadKey := errors.WithMeta("key1", "val1", 16, "val2", "key3", "val3")
	fmt.Printf("\n%#v", mBadKey)

	// 10 is not a string for a key it will replace it with !BADKEY2
	// key3 doesn't have a value it will add !BADVALUE to key3
	mBadKeyValue := errors.WithMeta("key1", "val1", 10, "val2", "key3")
	fmt.Printf("\n%#v", mBadKeyValue)
}

func ExampleMeta_Merge() {
	m := errors.WithMeta("key1", "val1")
	m = m.Merge("key2", "val2")

	fmt.Printf("%#v", m)

	// Output: errors.Meta{"key1":"val1", "key2":"val2"}
}
