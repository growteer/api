package apperrors

type BadInput struct {
	Wrapped error
	Message string
	Field   string
}

func (e BadInput) Error() string {
	return e.Message
}

func (e BadInput) Unwrap() error {
	return e.Wrapped
}
