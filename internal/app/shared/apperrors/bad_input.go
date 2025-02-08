package apperrors

type BadInput struct {
	Code    ErrCode
	Message string
	Wrapped error
}

func (e BadInput) Error() string {
	return e.Message
}

func (e BadInput) Unwrap() error {
	return e.Wrapped
}
