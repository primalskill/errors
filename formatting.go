package errors

import (
	"fmt"
	"strings"
)

// Error returns the error message and satisfies the stdlib Error interface.
func (e *Error) Error() string {
	return e.Msg
}

// PrettyPrint is a helper method to *Error.PrettyPrint. This should be used in development.
func PrettyPrint(err error) string {
	var e *Error
	is := As(err, &e)

	if is == false {
		return ""
	}

	return e.PrettyPrint()
}

// PrettyPrint will recursively print all embedded errors including all information on the error it can found.
// This should be used in development.
func (e *Error) PrettyPrint() string {
	var s strings.Builder
	err := Flatten(e)

	s.WriteString("\n")

	for i, elem := range err {
		if len(e.Msg) == 0 {
			s.WriteString("[empty]")
		} else {
			s.WriteString(elem.Msg)
		}

		s.WriteString(elem.Stack.stackPrettyString(true))
		s.WriteString(elem.Meta.metaPrettyString(true))

		if i < len(err) {
			s.WriteString("\n")
		}

	}

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

func (p *Stack) stackPrettyString(isSub bool) string {
	if len(string(*p)) == 0 {
		return ""
	}

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
