package errors

import (
	"fmt"
)

// Error returns the error message and satisfies the stdlib Error interface.
func (e *Error) Error() string {
	if len(e.Msg) == 0 {
		return "<empty>"
	}

	return e.Msg
}

// PrettyPrint is a helper method to *Error.PrettyPrint. This should only be used in development.
func PrettyPrint(err error) string {
	var e *Error
	is := As(err, &e)

	if is == false {
		return err.Error()
	}

	return e.PrettyPrint()
}

// PrettyPrint will recursively print all embedded errors including all information on the error it can found.
// This should only be used in development.
func (e *Error) PrettyPrint() string {
	err := Flatten(e)

	var b []byte

	b = append(b, '\n')

	for i, elem := range err {
		if len(e.Msg) == 0 {
			b = append(b, "<empty>"...)
		} else {
			b = append(b, elem.Msg...)
		}

		b = append(b, elem.Source.sourcePrettyString()...)
		b = append(b, elem.Meta.metaPrettyString()...)

		if i < len(err) {
			b = append(b, '\n')
		}
	}

	return string(b)
}

// String returns Meta in [key1:val1 key2:val2 ...] format and satisfies the fmt.Stringer interface.
func (p Meta) String() string {
	var b []byte

	b = append(b, '[')

	i := 0
	for k, v := range p {
		str, ok := v.(string)
		if ok == true {
			b = fmt.Appendf(b, "%s:%s", k, str)
		}

		i++

		if i < len(p) {
			b = append(b, ' ')
		}
	}

	b = append(b, ']')

	return string(b)
}

func (p *Source) sourcePrettyString() string {
	if len(string(*p)) == 0 {
		return ""
	}

	var b []byte
	b = fmt.Appendf(b, "\n%*s|- Source : %s", 2, " ", string(*p))

	return string(b)
}

func (p Meta) metaPrettyString() string {
	if len(p) == 0 {
		return ""
	}

	var b []byte

	b = fmt.Appendf(b, "\n%*s|- Meta :", 2, " ")

	for k := range p {
		b = fmt.Appendf(b, "\n%*s|- %s : %+v", 4, " ", k, p[k])
	}

	return string(b)
}
