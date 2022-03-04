package errors

import (
	"encoding/json"
	"testing"
)

func TestStack(t *testing.T) {
	t.Run("it should marshal JSON", marshalJSON)
}

func marshalJSON(t *testing.T) {
	err := E("test error")
	var e *Error
	As(err, &e)

	retJSON, _ := e.Stack.MarshalJSON()

	type cmp struct {
		Path string `json:"path"`
		Func string `json:"func"`
		Line int    `json:"line"`
	}

	var ret cmp

	uErr := json.Unmarshal(retJSON, &ret)
	if uErr != nil {
		t.Fatalf("Cannot unmarshal Stack JSON, error: %s", uErr.Error())
	}

	if len(ret.Path) == 0 {
		t.Fatalf("Unmarshaled Stack JSON path is empty")
	}

	if len(ret.Func) == 0 {
		t.Fatalf("Unmarshaled Stack JSON function name is empty")
	}
}
