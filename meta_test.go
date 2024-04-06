package errors

import (
	"testing"
)

func TestMeta(t *testing.T) {
	t.Run("it should store meta map", storeMetaMap)
	t.Run("it should store all keys and values", storeAllKeysValues)
	t.Run("it should return empty value for first key", returnEmptyValueFirstKey)
	t.Run("it should append !BADVALUE string as value", badValueStringAsValue)
	t.Run("it should set !BADKEY string for non-string key in args", setBadKeyNonStringKey)
	t.Run("it should set a key/value pair in the map", setKeyValuePair)
	t.Run("it should merge map to existing map", mergeMaps)
}

func storeMetaMap(t *testing.T) {
	m := WithMeta("key1", "val1")

	v, has := m["key1"]
	if has != true {
		t.Fatalf("WithMeta() key not found.")
	}

	if v != "val1" {
		t.Fatalf("WithMeta(), value mismatch, expected: val1, got: %+v", v)
	}
}

func storeAllKeysValues(t *testing.T) {
	m := WithMeta("key1", "val1", "key2", "val2")

	v2, has := m["key2"]
	if has != true {
		t.Fatalf("WithMeta() key2 not found.")
	}

	if v2 != "val2" {
		t.Fatalf("WithMeta(), value mismatch, expected: val2, got: %+v", v2)
	}
}

func returnEmptyValueFirstKey(t *testing.T) {
	m := WithMeta("key1")

	v, ok := m["key1"]
	if ok == false {
		t.Fatalf("WithMeta() key1 not found.")
	}

	if len(v.(string)) > 0 {
		t.Fatalf("WithMeta(), key1 value is not empty.")
	}
}

func badValueStringAsValue(t *testing.T) {
	m := WithMeta("key1", "val1", "key2", "val2", "key3")

	v, has := m["key3"]
	if has != true {
		t.Fatalf("WithMeta(), key2 not found, map: %+v", m)
	}

	if v.(string) != "!BADVALUE" {
		t.Fatalf("WithMeta(), expected !BADVALUE as last value for odd key/value pairs, got: %+v", v)
	}
}

func setBadKeyNonStringKey(t *testing.T) {
	m := WithMeta("key1", "val1", 100, "val2", "key3", 55)

	if len(m) != 3 {
		t.Fatalf("WithMeta(), map is wrong length, expected 3, got: %+v", len(m))
	}

	// test key1 / val1
	mv1, mhas1 := m["key1"]
	if mhas1 != true {
		t.Fatalf("WithMeta(), key1 is not found.")
	}

	if mv1 != "val1" {
		t.Fatalf("WithMeta(), key1 expected value 'val1', got: %+v", mv1)
	}

	// test !BADKEY2 / val2
	mv2, mhas2 := m["!BADKEY2"]
	if mhas2 != true {
		t.Fatalf("WithMeta(), !BADKEY2 as key is not found.")
	}

	if mv2 != "val2" {
		t.Fatalf("WithMeta(), !BADKEY_AT_INDEX_2 expected value 'val2', got: %+v", mv2)
	}

	// test key3 / 55
	mv3, mhas3 := m["key3"]
	if mhas3 != true {
		t.Fatalf("WithMeta(), key3 is not found.")
	}

	if mv3 != 55 {
		t.Fatalf("WithMeta(), key3 expected value 55, got: %+v", mv3)
	}
}

func setKeyValuePair(t *testing.T) {
	m := WithMeta("key1", "val1")
	m.Set("key2", "val2")

	v, has := m["key2"]
	if has != true {
		t.Fatalf("Set() key2 is not found.")
	}

	if v != "val2" {
		t.Fatalf("Set(), wrong value for key2, expected: 'val2', got: %+v", v)
	}
}

func mergeMaps(t *testing.T) {
	m := WithMeta("key1", "val1")
	m.Merge("key2", "val2")

	v, has := m["key2"]
	if has != true {
		t.Fatalf("Merge() key2 is not found.")
	}

	if v != "val2" {
		t.Fatalf("Merge(), wrong value for key2, expected: 'val2', got: %+v", v)
	}
}
