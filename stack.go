package errors

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

// Stack represents an error stack captuing the file path, function name and file line number where the error happened.
// A stack is always attached to an error automatically.
type Stack struct {
	FilePath string
	FuncName string
	Line     int
}

// MarshalJSON satisfies the encoding/json Marshaler interface and generates valid JSON from  Stack struct s.
func (s *Stack) MarshalJSON() ([]byte, error) {
	var ret strings.Builder

	ret.WriteString("{")
	ret.WriteString("\"path\":\"")
	ret.WriteString(s.FilePath)
	ret.WriteString("\",\"func\":\"")
	ret.WriteString(s.FuncName)
	ret.WriteString("\",\"line\":")
	ret.WriteString(strconv.Itoa(s.Line))
	ret.WriteString("}")

	return []byte(ret.String()), nil
}

// UnmarshalJSON satisfies the encoding/json Unmarshaler interface and parses JSON into Stack struct s.
func (s *Stack) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	if s == nil {
		return fmt.Errorf("nil receiver passed to UnmarshalJSON")
	}

	type alias struct{
		Path string `json:"path"`
		Func string `json:"func"`
		Line int    `json:"line"`
	}

	var ret alias

	err := json.Unmarshal(data, &ret)
	if err != nil {
		return err
	}

	s.FilePath = ret.Path
	s.FuncName = ret.Func
	s.Line = ret.Line

	return nil
}

// getStack will get the file path, function name and line number where the error happened.
func getStack() (s Stack) {
	// Index 3 will show the calling function data
	targetFrameIndex := 3

	// Set size to targetFrameIndex + 2 to ensure we have room for one more caller than we need.
	programCounters := make([]uintptr, targetFrameIndex+2)
	n := runtime.Callers(0, programCounters)

	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()
			if frameIndex == targetFrameIndex {
				s.FilePath = frameCandidate.File
				s.FuncName = frameCandidate.Function
				s.Line = frameCandidate.Line
			}
		}
	}

	// Can't extract the file path and line number
	if len(s.FilePath) == 0 {
		s.FilePath = "unknown"
		s.FuncName = "unknown"
		return
	}

	// Parse the function name
	i := strings.LastIndex(s.FuncName, "/")
	s.FuncName = s.FuncName[i+1:]
	i = strings.Index(s.FuncName, ".")
	s.FuncName = s.FuncName[i+1:]

	return
}
