package errors

import (
	"encoding/json"
	"fmt"
	"testing"
)

type dummyType struct {
	Err error `json:"err"`
}

func TestJSONMarshaling(t *testing.T) {
	t.Run("it should generate correct json payload", correctJSONPayload)
	t.Run("it should work with single error", singleErrorJSONPayload)
}

func correctJSONPayload(t *testing.T) {
	regErr1 := fmt.Errorf("reg error 1")
	regErr2 := fmt.Errorf("reg error 2: %w", regErr1)

	err1 := E("err1 error", WithMeta("test", 3445), regErr2)
	err2 := E("err2", err1)

	var payload dummyType
	payload.Err = err2

	b, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("expected json marshal nil error, got: %s", err.Error())
	}

	jsonCmp := `{"err":[{"msg":"err2","source":"json_test.go:23"},{"msg":"err1 error","source":"json_test.go:22","meta":{"test":3445}},{"msg":"reg error 2: reg error 1"},{"msg":"reg error 1"}]}`
	if string(b) != jsonCmp {
		t.Fatalf("json payload mismatch:\nexpected: %s\ngot: %s", jsonCmp, string(b))
	}
}

func singleErrorJSONPayload(t *testing.T) {
	var payload dummyType
	payload.Err = E("single")

	b, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("expected json marshal nil error, got: %s", err.Error())
	}

	jsonCmp := `{"err":{"msg":"single","source":"json_test.go:41"}}`
	if string(b) != jsonCmp {
		t.Fatalf("json payload mismatch:\nexpected: %s\ngot: %s", jsonCmp, string(b))
	}
}
