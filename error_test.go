package errors

import (
	"errors"
	"testing"
	"reflect"
)

func TestErrors(t *testing.T) {
	t.Run("it should store msg", storeMsg)
	t.Run("it should store meta", storeMeta)
	t.Run("it should store stack", storeStack)
	t.Run("it should wrap errors", wrapErrors)
	t.Run("it should get meta from error", getMetaFromError)
	t.Run("it should test if Meta is stored on error", hasMetaOnError)
	t.Run("it should return empty Meta if error is not goerror.Error", returnEmptyMeta)
	t.Run("it should unwrap all embedded errors", unwrapAllErrors)
}

func storeMsg(t *testing.T) {
	e := E("test error")
	ee := e.(*Error)

	if ee.Msg != "test error" {
		t.Fatalf("E() should store Msg, expected: 'test error', got: %s", ee.Msg)
	}
}

func storeMeta(t *testing.T) {
	e := E("test error", WithMeta("metaKey", "meta value"))
	ee := e.(*Error)

	if len(ee.Meta) != 1 {
		t.Fatalf("E() should store the passed Meta map.")
	}
}

func storeStack(t *testing.T) {
	e := E("test error")
	ee := e.(*Error)

	if len(ee.Stack.FilePath) == 0 {
		t.Fatalf("E() should store the stack, FilePath is empty.")
	}

	if len(ee.Stack.FuncName) == 0 {
		t.Fatalf("E() should store the stack, FuncName is empty.")
	}

	if ee.Stack.Line == 0 {
		t.Fatalf("E() should store the stack, Line is zero.")
	}
}

func wrapErrors(t *testing.T) {
	e1 := E("error 1")
	e2 := E("error 2", e1)

	ee2 := e2.(*Error)

	if ee2.Err == nil {
		t.Fatalf("E() should wrap errors, the Err field is nil")
	}
}

func getMetaFromError(t *testing.T) {
	e := E("test error", WithMeta("key1", "val1"))
	m := GetMeta(e)

	v, has := m["key1"]
	if has != true {
		t.Fatalf("GetMeta() returned Meta doesn't contain key1.")
	}

	if v != "val1" {
		t.Fatalf("GetMeta() returned Meta key1 value mismatch, expected: 'val1', got: %+v", v)
	}
}

func hasMetaOnError(t *testing.T) {
	e := E("test error", WithMeta("key1", "val1"))
	has := HasMeta(e)

	if has == false {
		t.Fatalf("HasMeta returns wrong bool, expected: true, got: false")
	}

	e = E("test error")
	has = HasMeta(e)

	if has == true {
		t.Fatalf("HasMeta returns wrong bool, expected: false, got: true")
	}
}


func returnEmptyMeta(t *testing.T) {
	e := errors.New("not goerror")
	m := GetMeta(e)

	if len(m) != 0 {
		t.Fatalf("GetMeta() returned Meta map is not empty.")
	}
}

func unwrapAllErrors(t *testing.T) {
	m1 := WithMeta("key1", "val1")
	m2 := WithMeta("key2", "val2")
	m3 := WithMeta("key3", "val3")

	e0 := errors.New("not goerror")
	e1 := E("e1", e0, m1)
	e2 := E("e2", e1, m2)
	e3 := E("e3", e2, m3)
	errs := UnwrapAll(e3)

	if len(errs) != 4 {
		t.Fatalf("Returned errs length mismatch, expected: 4, got: %d", len(errs))
	}

	for i, err := range errs {
		m := GetMeta(err)

		if i == 0 && reflect.DeepEqual(m, m3) == false {
			t.Fatalf("Meta is not preserved, expected: %+v, got: %+v", m3, m)
		} else if i == 1 && reflect.DeepEqual(m, m2) == false {
			t.Fatalf("Meta is not preserved, expected: %+v, got: %+v", m2, m)
		} else if i == 2 && reflect.DeepEqual(m, m1) == false {
			t.Fatalf("Meta is not preserved, expected: %+v, got: %+v", m1, m)
		} else if i == 3 && HasMeta(err) == true {
			t.Fatalf("Last error shouldn't have Meta, got: %+v", GetMeta(err))
		}

		if i == 0 && err.Error() != "e3" {
			t.Fatalf("Error message mismatch, expected: e3, got: %s", err.Error())
		} else if i == 1 && err.Error() != "e2" {
			t.Fatalf("Error message mismatch, expected: e2, got: %s", err.Error())
		} else if i == 2 && err.Error() != "e1" {
			t.Fatalf("Error message mismatch, expected: e1, got: %s", err.Error())
		} else if i == 3 && err.Error() != "not goerror" {
			t.Fatalf("Error message mismatch, expected: not goerror, got: %s", err.Error())
		}
	}


	errs = UnwrapAll(nil)
	if len(errs) > 0 {
		t.Fatalf("It should return empty slice on nil error, got: %+v", errs)	
	}
	
}
