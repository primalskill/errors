package errors

import (
	"fmt"
)

// Meta holds extra meta data around an error. Try to add simple values to the meta. Key order is not guaranteed.
type Meta map[string]interface{}

// WithMeta accepts an even number of arguments representing key, value pairs. The first argument `firstKey` forces
// the compiler to fail if the first argument is not a string. In `args` every odd argument must be of type string
// which will be used as the meta map key. If an odd argument is not a string that pair will be skipped
// (see example below).
//
// Examples:
//
// This is valid:
//   WithMeta("key1", 158, "key2", "value string", "myStruct", myStruct)
//
// These are invalid:
//   WithMeta(15) <-- will result in compile error
//   WithMeta("strKey1") <-- meta will be empty because there's no value to add to strKey1
//
// In this case the arguments 16 and "variable" will be skipped.
//   WithMeta("key1", 15, 16, "key3", "key4", "variable")
//
// The final output will be: key1 = 15, key3 = key4
func WithMeta(firstKey string, args ...interface{}) (m Meta) {
	if len(args) == 0 {
		return
	}

	m = make(map[string]interface{}, len(args)+1)

	// Set the firstKey to the first value
	m[firstKey] = args[0]

	if len(args) <= 2 {
		return
	}

	args = args[1:]

	// If args is not even add an empty string as a last value to make it even
	if len(args)%2 != 0 {
		args = append(args, "")
	}

	for i := 0; i < len(args); i = i + 2 {
		// Skip if the key is not string
		t := fmt.Sprintf("%T", args[i])
		if t != "string" {
			continue
		}

		m[args[i].(string)] = args[i+1]
	}

	return
}

// Set will set key to value and returns Meta. If key exists it will be overwritten.
func (p Meta) Set(key string, value interface{}) (m Meta) {
	p[key] = value
	return p
}

// Merge combines the arguments to an existing Meta and returns it. Existing keys will be overwritten.
func (p Meta) Merge(firstKey string, args ...interface{}) (m Meta) {
	nm := WithMeta(firstKey, args...)

	for k, v := range nm {
		p.Set(k, v)
	}

	return p
}
