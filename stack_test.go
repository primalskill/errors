package errors

import (
	"encoding/json"
	"testing"
)

type jsonCmp struct {
	Path string `json:"path"`
	Func string `json:"func"`
	Line int    `json:"line"`
}

func TestStack(t *testing.T) {
	t.Run("it should marshal JSON", marshalJSON)
	t.Run("it should unmarshal JSON", unmarshalJSON)
	t.Run("it should unmarhsal null JSON", unmarshalNullJSON)
}

func marshalJSON(t *testing.T) {
	err := E("test error")
	var e *Error
	As(err, &e)

	retJSON, _ := e.Stack.MarshalJSON()

	var ret jsonCmp

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

func unmarshalJSON(t *testing.T) {
	var cmp jsonCmp

	cmp.Path = "test path"
	cmp.Func = "test func"
	cmp.Line = 10

	j, jErr := json.Marshal(cmp)
	if jErr != nil {
		t.Fatalf("Cannot marshal JSON for comparison, got error: %s", jErr.Error())
	}

	err := E("test error")
	var e *Error
	As(err, &e)

	uErr := e.Stack.UnmarshalJSON(j)
	if uErr != nil {
		t.Fatalf("Cannot unmarshal JSON into Stack, got error: %s", uErr.Error())
	}

	if e.Stack.FilePath != "test path" {
		t.Fatalf("UnmarshalJSON file path mismatch, expected: test path, got: %s", e.Stack.FilePath)
	}

	if e.Stack.FuncName != "test func" {
		t.Fatalf("UnmarshalJSON func name mismatch, expected: test func, got: %s", e.Stack.FuncName)
	}

	if e.Stack.Line != 10 {
		t.Fatalf("UnmarshalJSON line mismatch, expected: 10, got: %d", e.Stack.Line)
	}
}


func unmarshalNullJSON(t *testing.T) {
	
	err := E("test error")
	var e *Error
	As(err, &e)

	uErr := e.Stack.UnmarshalJSON([]byte("null"))
	if uErr != nil {
		t.Fatalf("Cannot unmarshal JSON into Stack, got error: %s", uErr.Error())
	}

}
