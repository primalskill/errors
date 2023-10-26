package errors

import (
	"errors"
	"reflect"
	"testing"
)

func TestErrors(t *testing.T) {
	t.Run("it should store msg", storeMsg)
	t.Run("it should store meta", storeMeta)
	t.Run("it should store stack", storeStack)
	t.Run("it should wrap errors", wrapErrors)
	t.Run("it should add args on existing error", mirrorErrors)
	t.Run("M() should return a new error value ", mirrorErrorReturnNew)
	t.Run("it should get meta from error", getMetaFromError)
	t.Run("it should return empty Meta if error is not goerror.Error", returnEmptyMeta)
	t.Run("it should merge Meta to error", mergeMetaToError)
	t.Run("it should merge Meta to error with existing Meta", mergeMetaToErrorExistingMeta)
	t.Run("it should fail merge Meta on regular error", mergeMetaToRegularError)
	t.Run("it should flatten all embedded errors", flattenAllErrors)
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

	if len(string(ee.Stack)) == 0 {
		t.Fatalf("E() should store the stack, got empty string.")
	}
}

func wrapErrors(t *testing.T) {
	e1 := E("error 1")
	e2 := E("error 2", e1)

	ee2 := e2.(*Error)

	if ee2.err == nil {
		t.Fatalf("E() should wrap errors, the Err field is nil")
	}
}

func mirrorErrors(t *testing.T) {
	e0 := E("error 0")
	e1 := E("error 1")

	e2 := M(e1, e0)

	var ee2 *Error
	As(e2, &ee2)

	if ee2.Msg != "error 1" {
		t.Fatalf("e2 should have message 'error 1', got: %s", ee2.Msg)
	}

	ee0 := ee2.Unwrap()

	if ee0.Error() != "error 0" {
		t.Fatalf("ee0 sould have message 'error 0', got: %s", ee0.Error())
	}

	// Merges Meta to existing error
	e0 = E("error 0", WithMeta("k1", "v1"))
	e1 = M(e0, WithMeta("k2", "v2"))

	retMeta, _ := GetMeta(e1)
	cmpMeta := WithMeta("k1", "v1", "k2", "v2")

	if reflect.DeepEqual(retMeta, cmpMeta) == false {
		t.Fatalf("M() should merge Meta to existing error, expected: %+v, got: %+v", cmpMeta, retMeta)
	}

	// It should convert regular error to errors.Error
	regErr := errors.New("regular error")
	someErr := E("some error")

	wRegErr := M(regErr, someErr)

	var ee_wRegErr *Error
	As(wRegErr, &ee_wRegErr)

	if ee_wRegErr.Msg != "regular error" {
		t.Fatalf("M() should convert regular error to errors.Error")
	}

	if ee_wRegErr.err == nil {
		t.Fatalf("M() should overwrite arguments on converted regular error")
	}
}

func mirrorErrorReturnNew(t *testing.T) {
	m0 := WithMeta("e0k1", "e0v1")
	m1 := WithMeta("e1k1", "e1v1", "e1k2", "e1v2")

	e0 := E("e0 error", m0)
	e1 := M(e0, m1)

	if Is(e1, e0) == false {
		t.Fatalf("e1 and e0 should be identical when using M()")
	}

	uErr := errors.New("some error")
	if Is(e1, uErr) == true {
		t.Fatalf("e1 and uErr shouldn't be identical when using M()")
	}

	if Is(e0, uErr) == true {
		t.Fatalf("e0 and uErr shouldn't be identical when using M()")
	}

	// Preload regular error
	mre1 := M(uErr)
	if Is(mre1, uErr) == false {
		t.Fatalf("mre1 and uErr should be identical when using M()")
	}

	mre2 := E("test err", mre1)
	if Is(mre2, uErr) == false {
		t.Fatalf("mre2 and uErr should be identical when using M()")
	}
}

func getMetaFromError(t *testing.T) {
	e := E("test error", WithMeta("key1", "val1"))
	m, _ := GetMeta(e)

	v, has := m["key1"]
	if has != true {
		t.Fatalf("GetMeta() returned Meta doesn't contain key1.")
	}

	if v != "val1" {
		t.Fatalf("GetMeta() returned Meta key1 value mismatch, expected: 'val1', got: %+v", v)
	}
}

func returnEmptyMeta(t *testing.T) {
	e := errors.New("not goerror")
	m, _ := GetMeta(e)

	if len(m) != 0 {
		t.Fatalf("GetMeta() returned Meta map is not empty.")
	}
}

func mergeMetaToError(t *testing.T) {
	err := E("test error")
	m := WithMeta("metaKey1", "metaVal1")

	MergeMeta(err, m)

	mCmp, has := GetMeta(err)
	if has == false {
		t.Fatalf("err should have a Meta map, got: %+v", mCmp)
	}

	if reflect.DeepEqual(m, mCmp) == false {
		t.Fatalf("err should have a merged Meta map\n - expected: %+v\n - got: %+v", m, mCmp)
	}
}

func mergeMetaToErrorExistingMeta(t *testing.T) {
	err := E("test error with meta", WithMeta("key1", "val1"))
	m := WithMeta("metaKey1", "metaVal1")

	MergeMeta(err, m)

	mCmp, has := GetMeta(err)
	if has == false {
		t.Fatalf("err with meta should have a Meta map, got: %+v", mCmp)
	}

	mRet := WithMeta("key1", "val1", "metaKey1", "metaVal1")

	if reflect.DeepEqual(mRet, mCmp) == false {
		t.Fatalf("err with meta should have a merged Meta map\n - expected: %+v\n - got: %+v", mRet, mCmp)
	}
}

func mergeMetaToRegularError(t *testing.T) {
	err := errors.New("regular error")
	m := WithMeta("metaKey1", "metaVal1")

	_, is := MergeMeta(err, m)
	if is == true {
		t.Fatalf("MergeMeta should fail when err is a regular error")
	}
}

func flattenAllErrors(t *testing.T) {

	// nil error
	errs := Flatten(nil)

	if len(errs) > 0 {
		t.Fatalf("It should return an empty slice on nil error, got: %+v", errs)
	}

	m1 := WithMeta("key1", "val1")
	m2 := WithMeta("key2", "val2")

	e0 := errors.New("not goerror")
	e1 := E("e1", e0, m1)
	e2 := E("e2", e1, m2)

	errs = Flatten(e2)

	for i, err := range errs {

		// Metas
		if i == 0 && reflect.DeepEqual(err.Meta, m2) == false {
			t.Fatalf("Meta is not preserved, expected: %+v, got: %+v", m2, err.Meta)
		} else if i == 1 && reflect.DeepEqual(err.Meta, m1) == false {
			t.Fatalf("Meta is not preserved, expected: %+v, got: %+v", m1, err.Meta)
		} else if i == 2 && len(err.Meta) > 0 {
			t.Fatalf("Meta length on last error is not empty, got: %+v", err.Meta)
		}

		// Messages
		if i == 0 && err.Error() != "e2" {
			t.Fatalf("Error message mismatch, expected: e2, got: %s", err.Error())
		} else if i == 1 && err.Error() != "e1" {
			t.Fatalf("Error message mismatch, expected: e1, got: %s", err.Error())
		} else if i == 3 && err.Error() != "not goerror" {
			t.Fatalf("Error message mismatch, expected: not goerror, got: %s", err.Error())
		}
	}
}
