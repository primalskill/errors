package errors

import (
	"fmt"
	"strings"
)

// Error returns the error message and satisfies the stdlib Error interface.
func (e *Error) Error() string {
	return e.Msg
}

// ErrorFull returns the error string and any additional meta data and stack or an empty string if err
// is not of type Error. It doesn't unwrap the error.
func ErrorFull(err error) string {
	var e *Error
	isError := As(err, &e)

	if isError == false {
		return ""
	}

	return e.ErrorFull()
}


// ErrorFull returns the error string and any additional meta data and stack. It doesn't unwrap the error.
func (e *Error) ErrorFull() string {
	var s strings.Builder

	// Msg
	if len(e.Msg) == 0 {
		s.WriteString("[empty]")
	} else {
		s.WriteString(e.Msg)
	}

	s.WriteString(e.Stack.stackPrettyString(true))
	s.WriteString(e.Meta.metaPrettyString(true))
	s.WriteString("\n")

	return s.String()
}

// String returns Meta in [key1:val1 key2:val2 ...] format and satisfies the fmt.Stringer interface.
func (p Meta) String() string {
	var ret strings.Builder

	ret.WriteString("[")

	i := 0
	for k, v := range p {
		ret.WriteString(k)
		ret.WriteString(":")
		ret.WriteString(v.(string))
		i++

		if i < len(p) {
			ret.WriteString(" ")
		}
	}

	ret.WriteString("]")

	return ret.String()
}

// PrettyPrint returns Meta pretty formatted. Should be used in development.
func (p *Meta) PrettyPrint() string {
	return p.metaPrettyString(false)
}

func (p *Stack) stackPrettyString(isSub bool) string {
	padding := ""
	pipe := ""

	if isSub {
		padding = "  "
		pipe = "|- "
	}

	
	return fmt.Sprintf("\n%s%sStack: %s", padding, pipe, string(*p))	
}

func (p Meta) metaPrettyString(isSub bool) string {
	if len(p) == 0 {
		return ""
	}

	padding := ""
	pipe := ""
	subPadding := "  "
	subPipe := "|- "

	if isSub {
		padding = "  "
		pipe = "|- "
	}

	var s strings.Builder
	s.WriteString(fmt.Sprintf("\n%s%sMeta:", padding, pipe))

	for k := range p {
		s.WriteString(fmt.Sprintf("\n%s%s%s", padding, subPadding, subPipe))
		s.WriteString(fmt.Sprintf("%s: %+v", k, p[k]))
	}

	return s.String()
}
