package errors

import (
	"fmt"
	"strconv"
	"strings"
)

// Error returns the error message and satisfies the stdlib Error interface.
func (e *Error) Error() string {
	return e.Msg
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

// String returns Stack in FilePath:FuncName:Line format.
func (s *Stack) String() string {
	
	var ret strings.Builder
	ret.WriteString(s.FilePath)
	ret.WriteString(":")
	ret.WriteString(s.FuncName)
	ret.WriteString(":")
	ret.WriteString(strconv.Itoa(s.Line))

	return ret.String()
}

// PrettyPrint returns Stack pretty formatted.
func (s *Stack) PrettyPrint() string {
	return s.stackPrettyString(false)
}

// String returns Meta in [key1:val1 key2:val2 ...] format.
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

// PrettyPrint returns Meta pretty formatted.
func (p *Meta) PrettyPrint() string {
	return p.metaPrettyString(false)
}



func (p *Stack) stackPrettyString(isSub bool) string {
	if len(p.FilePath) == 0 && len(p.FuncName) == 0 && p.Line == 0 {
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
	s.WriteString(fmt.Sprintf("\n%s%sStack:\n", padding, pipe))

	if len(p.FilePath) > 0 {
		s.WriteString(fmt.Sprintf("%s%s%sFile Path: ", padding, subPadding, subPipe))
		s.WriteString(p.FilePath)
		s.WriteString("\n")
	}

	if len(p.FuncName) > 0 {
		s.WriteString(fmt.Sprintf("%s%s%sFunction Name: ", padding, subPadding, subPipe))
		s.WriteString(p.FuncName)
		s.WriteString("\n")
	}

	if p.Line > 0 {
		s.WriteString(fmt.Sprintf("%s%s%sLine Number: ", padding, subPadding, subPipe))
		s.WriteString(strconv.Itoa(p.Line))
	}

	return s.String()
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
