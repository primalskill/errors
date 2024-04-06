package errors

import (
	"encoding/json"
)

// MarshalJSON implements json.Marshaler for Error. Need to use a workaround for encoding because it will create
// an infinite loop if json.Marshal() is used in MarshalJSON().
func (e *Error) MarshalJSON() ([]byte, error) {
	errs := Flatten(e)

	var b []byte

	if len(errs) == 1 {
		b, err := marshalJSONErr(errs[0])
		return b, err
	}

	// open json array
	b = append(b, '[')

	for _, err := range errs {
		ib, ierr := marshalJSONErr(err)
		if ierr != nil {
			return ib, ierr
		}

		b = append(b, ib...)
		b = append(b, ',')
	}

	// remove last ,
	b = b[:len(b)-1]

	// close json array
	b = append(b, ']')

	return b, nil
}

func marshalJSONErr(err Error) ([]byte, error) {
	b, merr := json.Marshal(err)

	return b, merr
}
