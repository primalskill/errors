package errors

type Error struct {
	Err   error
	Msg   string
	Stack Stack
	Meta  Meta
}

// E a new error and sets the required msg argument as the error message. Additional arguments like a Meta map or another error can be passed
// into the function that will be set on the error.
func E(msg string, args ...interface{}) error {
	e := &Error{}
	e.Msg = msg
	e.Stack = getStack()

	// Parse the arguments passed to the function
	for _, arg := range args {
		switch arg := arg.(type) {
		case Meta:
			e.Meta = arg
		case error:
			e.Err = arg
		}
	}

	return e
}

// GetMeta returns a Meta map or an empty Meta if the error doesn't contain a Meta or the error is not of type errors.Error.
func GetMeta(err error) Meta {
	eerr, ok := err.(*Error)

	if !ok {
		return make(Meta, 1)
	}

	return eerr.Meta
}

// HasMeta returns true if the error has Meta attached to it, otherwise false.
func HasMeta(err error) bool {
	ee, ok := err.(*Error)

	if !ok {
		return false
	}

	if len(ee.Meta) == 0 {
		return false
	}

	return true
}

// UnwrapAll returns err unwrapped into a slice of errors.
func UnwrapAll(err error) (errs []error) {
	if err == nil {
		return
	}

	var rec func(err error) error

	rec = func(recErr error) error {
		if recErr == nil {
			return nil
		}

		errs = append(errs, recErr)

		return rec(Unwrap(recErr))
	}

	rec(err)

	return
}

/*
// Flatten returns a slice of Error from embedded err.
func Flatten(err error) (ret []Error) {
	return
}
*/

// Unwrap returns the error one level deep otherwise nil. This is a proxy method for Unwrap().
func (e *Error) Unwrap() error {
	return e.Err
}
