// Package errorprivacy wraps causes with caller-supplied fixed categories, not
// content sanitization. Wrap owns no handles/workers; causes remain borrowed and
// inspectable through Unwrap, including input/credential-bearing text.
package errorprivacy

// Wrap retains an inspectable cause without including its text in Error.
// Categories must be library-owned constants, never input or callback text.
func Wrap(category string, cause error) error {
	if cause == nil {
		return nil
	}
	return &failure{category: category, cause: cause}
}

type failure struct {
	category string
	cause    error
}

func (e *failure) Error() string { return e.category }
func (e *failure) Unwrap() error { return e.cause }
