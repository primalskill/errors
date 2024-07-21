package errors

import (
	"testing"
)

func TestMatchers(t *testing.T) {
	t.Run("it should match message", itShouldMatchMessage)
	t.Run("it should match containing message", itShouldMatchContainingMessage)
}

func itShouldMatchMessage(t *testing.T) {
	err1 := E("this is an error", WithMeta("key1", "val1"))
	err2 := E("this is error2", err1, WithMeta("key2", "val2"))

	has := HasMessage(err2, "this is an error")
	if has == false {
		t.Fatalf("HasMessage() return expected true, got false")
	}

	has = HasMessage(err1, "this is an error")
	if has == false {
		t.Fatalf("HasMessage() should match a single err chain")
	}

	has = HasMessage(err2, "this is error2")
	if has == false {
		t.Fatalf("HasMessage() should match the top-most wrapped error in the chain")
	}

	has = HasMessage(err2, "ThIs Is An ErRor")
	if has == false {
		t.Fatalf("HasMessage() match should be case-insensitive")
	}

	has = HasMessage(err2, "error2")
	if has == true {
		t.Fatalf("HasMessage() should match exact messages")
	}
}

func itShouldMatchContainingMessage(t *testing.T) {
	err1 := E("this is an error", WithMeta("key1", "val1"))
	err2 := E("this is error2", err1, WithMeta("key2", "val2"))

	has := ContainsMessage(err2, "an error")
	if has == false {
		t.Fatalf("ContainsMessage() should match partial messages")
	}

	has = ContainsMessage(err2, "an erROR")
	if has == false {
		t.Fatalf("ContainsMessage() should match case-insensitive partial messages")
	}

	has = ContainsMessage(err1, "an erROR")
	if has == false {
		t.Fatalf("ContainsMessage() should match a single err chain")
	}

	has = ContainsMessage(err2, "erROR2")
	if has == false {
		t.Fatalf("ContainsMessage() should match the top-most wrapped error in the chain")
	}
}
