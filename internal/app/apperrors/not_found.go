package apperrors

type NotFound struct {
	Message string
	Wrapped error
}

func (e NotFound) Error() string {
	return e.Message
}

func (e NotFound) Unwrap() error {
	return e.Wrapped
}
