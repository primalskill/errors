package errors

import "strings"

// HasMessage reports whether any error in err's chain matches msg. It performs a case-insensitive matching.
//
// This is helpful when the error type's in err's chain is unknown and a string matching is preferred.
func HasMessage(err error, msg string) bool {
	errs := Flatten(err)

	for _, e := range errs {
		if strings.EqualFold(e.Msg, msg) {
			return true
		}
	}

	return false
}

// ContainsMessage reports whether any error in err's chain contains msg. It performs a case-insensitive matching.
//
// This is helpful when the error type's in err's chain is unknown and a string matching is preferred.
func ContainsMessage(err error, msg string) bool {
	errs := Flatten(err)
	lowerMsg := strings.ToLower(msg)

	for _, e := range errs {
		if strings.Contains(strings.ToLower(e.Msg), lowerMsg) {
			return true
		}
	}

	return false
}
