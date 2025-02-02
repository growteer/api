package apperrors

type Internal struct {
	Code    ErrCode
	Message string
	Wrapped error
}

func (e Internal) Error() string {
	return e.Message
}

func (e Internal) Unwrap() error {
	return e.Wrapped
}
