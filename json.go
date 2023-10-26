package errors

import (
	"encoding/json"
	"fmt"
)

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
	var b []byte
	b = append(b, '{')
	b = fmt.Appendf(b, "\"msg\":\"%s\"", err.Error())

	if len(err.Stack) > 0 {
		b = fmt.Appendf(b, ",\"stack\":\"%s\"", err.Stack)
	}

	if len(err.Meta) > 0 {
		m, merr := json.Marshal(err.Meta)
		if merr != nil {
			return b, merr
		}

		b = fmt.Appendf(b, ",\"meta\":%s", m)
	}

	b = append(b, '}')

	return b, nil
}
