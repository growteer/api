package apperrors

type Unauthorized struct {
	Code    ErrCode
	Message string
	Wrapped error
}

func (e Unauthorized) Error() string {
	return e.Message
}

func (e Unauthorized) Unwrap() error {
	return e.Wrapped
}
