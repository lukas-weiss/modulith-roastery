// Package check is the Go counterpart of Kotlin's require(): it reports a bad argument.
// The web package answers such errors with 400 Bad Request, as the Kotlin controllers do with
// IllegalArgumentException.
package check

// ArgumentError says that a caller passed a value that breaks a rule.
type ArgumentError struct {
	Message string
}

func (e *ArgumentError) Error() string { return e.Message }

// Require returns an ArgumentError with message if ok is false, otherwise nil.
func Require(ok bool, message string) error {
	if ok {
		return nil
	}
	return &ArgumentError{Message: message}
}
