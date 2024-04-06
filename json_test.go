package errors

import (
	"encoding/json"
	"fmt"
	"testing"
)

type dummyType struct {
	Err error `json:"err"`
}

type cmpMultiDummyType struct {
	Err []Error `json:"err"`
}

type cmpSingleDummyType struct {
	Err Error `json:"err"`
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

	var cmp cmpMultiDummyType
	err = json.Unmarshal(b, &cmp)
	if err != nil {
		t.Fatalf("expected json unmarshal nil error, got: %s", err.Error())
	}

	if len(cmp.Err) != 4 {
		t.Fatalf("expected cmp.Err length 4, got: %d", len(cmp.Err))
	}

	for i, elem := range cmp.Err {
		if i == 0 && elem.Msg != "err2" {
			t.Fatalf("expected cmp.Err[0].Msg = err2, got: %s", elem.Msg)
		}
		if i == 0 && len(elem.Meta) != 0 {
			t.Fatalf("expected cmp.Err[0].Meta length 0, got: %d", len(elem.Meta))
		}

		if i == 1 && elem.Msg != "err1 error" {
			t.Fatalf("expected cmp.Err[1].Msg = err1 error, got: %s", elem.Msg)
		}
		if i == 1 && len(elem.Meta) != 1 {
			t.Fatalf("expected cmp.Err[1].Meta length 1, got: %d", len(elem.Meta))
		}

		if i == 2 && elem.Msg != "reg error 2: reg error 1" {
			t.Fatalf("expected cmp.Err[2].Msg = reg error 2: reg error 1, got: %s", elem.Msg)
		}
		if i == 2 && len(elem.Meta) != 0 {
			t.Fatalf("expected cmp.Err[2].Meta length 0, got: %d", len(elem.Meta))
		}

		if i == 3 && elem.Msg != "reg error 1" {
			t.Fatalf("expected cmp.Err[3].Msg = reg error 1, got: %s", elem.Msg)
		}
		if i == 3 && len(elem.Meta) != 0 {
			t.Fatalf("expected cmp.Err[3].Meta length 0, got: %d", len(elem.Meta))
		}
	}
}

func singleErrorJSONPayload(t *testing.T) {
	var payload dummyType
	payload.Err = E("single")

	b, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("expected json marshal nil error, got: %s", err.Error())
	}

	var cmp cmpSingleDummyType
	err = json.Unmarshal(b, &cmp)
	if err != nil {
		t.Fatalf("expected json unmarshal nil error, got: %s", err.Error())
	}

	if cmp.Err.Msg != "single" {
		t.Fatalf("expected cmp.Err.Msg = single, got: %s", cmp.Err.Msg)
	}

	if len(cmp.Err.Meta) != 0 {
		t.Fatalf("expected cmp.Err.Meta length 0, got: %d", len(cmp.Err.Meta))
	}
}
