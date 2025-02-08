package apperrors

type Unauthenticated struct {
	Code    ErrCode
	Message string
	Wrapped error
}

func (e Unauthenticated) Error() string {
	return e.Message
}

func (e Unauthenticated) Unwrap() error {
	return e.Wrapped
}
