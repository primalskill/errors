package errors

import (
	"fmt"
)

// Meta holds extra meta data around an error. Try adding simple values to the Meta map. Key order is not guaranteed.
type Meta map[string]any

// WithMeta accepts an even number of arguments representing key/value pairs. The first argument "firstKey" forces
// the compiler to fail if the first argument is not a string. In "args" every odd argument must be of type string
// which will be used as the Meta map key. If an odd argument is not a string that pair will be skipped.
func WithMeta(firstKey string, args ...any) Meta {
	m := make(Meta, len(args)+1)

	// No arguments are passed, return an empty Meta.
	if len(args) == 0 {
		m[firstKey] = ""

		return m
	}

	// Set the firstKey, if only one args is present return early.
	m[firstKey] = args[0]
	if len(args) == 1 {
		return m
	}

	// Pop the first argument as this is already set.
	args = args[1:]

	// If args is not even add an !BADVALUE string as a last value to make it even and to let the user know something's
	// wrong.
	if len(args)%2 != 0 {
		args = append(args, "!BADVALUE")
	}

	// Loop over the rest of the args and set the key/value pairs in m.
	for i := 0; i < len(args); i = i + 2 {

		// If the even args are not string, replace it with !BADKEY<index> to let the user know it's settings the Meta values
		// wrong.
		strKey, ok := args[i].(string)
		if ok == false {
			strKey = fmt.Sprintf("!BADKEY%d", i+2)
		}

		m[strKey] = args[i+1]
	}

	return m
}

// Set will set key to value and returns Meta. Same keys will be overwritten.
func (p Meta) Set(key string, value any) (m Meta) {
	p[key] = value
	return p
}

// Merge combines the arguments to an existing Meta and returns it. Existing keys will be overwritten.
func (p Meta) Merge(firstKey string, args ...any) (m Meta) {
	nm := WithMeta(firstKey, args...)

	for k, v := range nm {
		p.Set(k, v)
	}

	return p
}
