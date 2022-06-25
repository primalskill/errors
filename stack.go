package errors

import (
	"runtime"
	"strconv"
	"strings"
)

// Stack represents an error stack captuing the file path and line number where the error happened in the format
// <file path>:<line nunmber>. A Stack is always attached to an error automatically.
type Stack string

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
				var sb strings.Builder

				filePath := frameCandidate.File
				if len(filePath) == 0 {
					filePath = "unknown"
				}

				sb.WriteString(filePath)
				sb.WriteString(":")
				sb.WriteString(strconv.Itoa(frameCandidate.Line))

				s = Stack(sb.String())
			}
		}
	}

	return
}
