package apperrors

type ErrCode = string

const (
	ErrCodeUnauthenticated       ErrCode = "auth.not_authenticated"
	ErrCodeInvalidCredentials    ErrCode = "auth.invalid_credentials"
	ErrCodeUserNotOnboarded      ErrCode = "auth.user_not_onboarded"
	ErrCodeInternalError         ErrCode = "internal"
	ErrCodeInvalidDateTimeFormat ErrCode = "validation.datetime_format"
	ErrCodeCouldNotSaveProfile   ErrCode = "profile.could_not_save"
	ErrCodeNotFound              ErrCode = "no_found"
	ErrCodeInvalidInput          ErrCode = "invalid_input"
)
