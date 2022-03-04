package errors

import (
	"fmt"
)

// Meta holds extra meta data around an error. Try adding simple values to the Meta map. Key order is not guaranteed.
type Meta map[string]interface{}

// WithMeta accepts an even number of arguments representing key/value pairs. The first argument "firstKey" forces
// the compiler to fail if the first argument is not a string. In "args" every odd argument must be of type string
// which will be used as the Meta map key. If an odd argument is not a string that pair will be skipped.
func WithMeta(firstKey string, args ...interface{}) (m Meta) {
	m = make(map[string]interface{}, len(args)+1)

	// Set the firstKey to the first value
	if len(args) == 0 {
		m[firstKey] = ""

		return
	}

	m[firstKey] = args[0]

	if len(args) == 1 {
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

// Set will set key to value and returns Meta. Same keys will be overwritten.
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
