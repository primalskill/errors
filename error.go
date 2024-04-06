package errors

import (
	stderrors "errors"
)

type Error struct {
	withFlag error
	err      error
	Msg      string `json:"msg"`
	Source   Source `json:"source,omitempty"`
	Meta     Meta   `json:"meta,omitempty"`
}

// parseArgTypes parses the arguments passed to the function
func (e *Error) parseArgTypes(args ...any) {
	for _, arg := range args {
		switch arg := arg.(type) {

		case Meta:
			// Merge the meta into the existing map on the err
			if len(e.Meta) == 0 {
				e.Meta = arg
			} else {
				for k, v := range arg {
					e.Meta.Set(k, v)
				}
			}

		case error:
			e.err = arg
		}
	}
}

// E return a new error and sets the required msg argument as the error message. Additional arguments like a Meta map or another error can be passed
// into the function that will be set on the error.
func E(msg string, args ...any) error {
	e := &Error{}
	e.Msg = msg
	e.Source = getSource()

	e.parseArgTypes(args...)

	return e
}

// M preloads err with all its Meta and wrapped errors if err is of type Error, otherwise it creates a new error of type Error and
// adds args on it. Passing in a regular error as err in the argument converts err to Error.
func M(err error, args ...any) error {
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

	// Overwrite the source to where M() was called, otherwise source will point to where err was instantiated.
	e.Source = getSource()

	// Parse the args too
	e.parseArgTypes(args...)

	return e
}

// With is deprecated, see M. It is kept for backwards compatibility.
func With(err error, args ...any) error {
	return M(err, args...)
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

// MergeMeta will merge m to err.Meta if err is of type errors.Error and returns TRUE if the operation was successful,
// FALSE otherwise.
func MergeMeta(err error, m Meta) (bool, error) {
	var e *Error

	isError := As(err, &e)
	if isError == false {
		return false, err
	}

	if e.Meta == nil {
		e.Meta = make(Meta, 1)
	}

	for k, v := range m {
		e.Meta.Set(k, v)
	}

	return true, e
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

func (e Error) As(target any) bool {
	if stderrors.As(e.withFlag, target) {
		return true
	}

	return stderrors.As(e.err, target)
}

// Unwrap returns the error one level deep otherwise nil. This is a proxy method for Unwrap().
func (e *Error) Unwrap() error {
	return e.err
}
