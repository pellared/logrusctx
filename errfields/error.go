package errfields

import "errors"

// Error contains custom fields and wrapped error.
type Error struct {
	Err    error
	Fields map[string]interface{}
}

// Unwrap returns the wrapped error.
func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) Error() string {
	return e.Err.Error()
}

// Add adds a field to *errfields.Error in the error chain.
// Returns a new *errfields.Error if it does not exits in the chain.
func Add(err error, key string, value interface{}) error {
	var e *Error
	if errors.As(err, &e) {
		e.Fields[key] = value
		return err
	}
	return &Error{
		Err:    err,
		Fields: map[string]interface{}{key: value},
	}
}
