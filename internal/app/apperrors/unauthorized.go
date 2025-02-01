package apperrors

type Unauthorized struct {
	Message           string
	Wrapped           error
	MissingPermission string
}

func (e Unauthorized) Error() string {
	return e.Message
}

func (e Unauthorized) Unwrap() error {
	return e.Wrapped
}
