package apperrors

type Unauthenticated struct {
	Message string
	Wrapped error
}

// implement the error interface for the custom error type
func (e Unauthenticated) Error() string {
	return e.Message
}

func (e Unauthenticated) Unwrap() error {
	return e.Wrapped
}
