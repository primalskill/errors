package errors

import (
	"testing"
)

func TestStdErrors(t *testing.T) {
	t.Run("it should have the same error signature", isError)
}

func isError(t *testing.T) {
	var validationErr = E("validation error test string", WithMeta("field1", "required"))

	err1 := E("err1 error", validationErr)
	err2 := E("err2 error", err1)

	if !Is(err2, validationErr) {
		t.Fatalf("Errors should contain validationErr, got: %+v", err2)
	}
}
