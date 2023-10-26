package errors

import (
	stderrors "errors"
)

type Error struct {
	withFlag error
	err      error
	Msg      string `json:"msg"`
	Stack    Stack  `json:"stack"`
	Meta     Meta   `json:"meta"`
}

// parseArgTypes parses the arguments passed to the function
func (err *Error) parseArgTypes(args ...interface{}) {
	for _, arg := range args {
		switch arg := arg.(type) {

		case Meta:
			// Merge the meta into the existing map on the err
			if len(err.Meta) == 0 {
				err.Meta = arg
			} else {
				for k, v := range arg {
					err.Meta.Set(k, v)
				}
			}

		case error:
			err.err = arg
		}
	}
}

// E return a new error and sets the required msg argument as the error message. Additional arguments like a Meta map or another error can be passed
// into the function that will be set on the error.
func E(msg string, args ...interface{}) error {
	e := &Error{}
	e.Msg = msg
	e.Stack = getStack()

	e.parseArgTypes(args...)

	return e
}

// With adds args to err if err is of type Error, othwerwise it creates a new error of type Error and adds args
// on that error. Passing in a regular error (not errors.Error) to the err argument will convert the error to
// errors.Error therefore trying to compare errors with Is() will return FALSE.
func With(err error, args ...interface{}) error {
	e := &Error{}
	ec, is := err.(*Error)

	if is == false {
		// We're dealing with an error other than errors.Error
		e.Msg = err.Error()
	} else {
		// We're dealing with error.Error, copy over the data

		e.err = ec.err
		e.Msg = ec.Msg

		// If the original error have Meta, copy over onto the new error
		if len(ec.Meta) > 0 {
			e.Meta = make(Meta, 1)

			for k, v := range ec.Meta {
				e.Meta[k] = v
			}
		}
	}

	// Set the withFlag to the original so the Is() and As() functions still have the correct behavior.
	e.withFlag = err

	// Overwrite the stack to where With() was called, otherwise stack will point to where err was instantiated.
	e.Stack = getStack()

	// Parse the args too
	e.parseArgTypes(args...)

	return e
}

// GetMeta returns a Meta map or an empty Meta if the error doesn't contain a Meta or the error is not of type
// errors.Error. The second returned argument is TRUE if the err has a Meta, FALSE otherwise.
func GetMeta(err error) (Meta, bool) {
	eerr, ok := err.(*Error)

	if !ok {
		return make(Meta, 1), false
	}

	return eerr.Meta, true
}

// MergeMeta will merge m to err.Meta if err is of type errors.Error and returns TRUE if the operation was succesful,
// FALSE otherwise.
func MergeMeta(err error, m Meta) (error, bool) {
	var e *Error

	isError := As(err, &e)
	if isError == false {
		return err, false
	}

	if e.Meta == nil {
		e.Meta = make(Meta, 1)
	}

	for k, v := range m {
		e.Meta.Set(k, v)
	}

	return e, true
}

// Flatten returns a slice of Error from embedded err.
func Flatten(err error) (ret []Error) {
	if err == nil {
		return
	}

	uErr := err

	for ok := true; ok; ok = (uErr != nil) {
		var e *Error
		cOk := As(uErr, &e)

		if cOk == true {
			ret = append(ret, *e)
			uErr = Unwrap(uErr)

			continue
		}

		// This is a regular error, convert it to errors.Error
		var ee Error
		ee.Msg = uErr.Error()
		ee.err = uErr

		ret = append(ret, ee)
		uErr = Unwrap(uErr)
	}

	return
}

func (e Error) Is(target error) bool {
	if stderrors.Is(e.withFlag, target) {
		return true
	}

	return stderrors.Is(e.err, target)
}

func (e Error) As(target interface{}) bool {
	if stderrors.As(e.withFlag, target) {
		return true
	}

	return stderrors.As(e.err, target)
}

// Unwrap returns the error one level deep otherwise nil. This is a proxy method for Unwrap().
func (e *Error) Unwrap() error {
	return e.err
}
