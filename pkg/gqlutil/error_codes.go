package gqlutil

const (
	ErrCodeUnauthenticated ErrCode = "auth.not_authenticated"
	ErrCodeInvalidCredentials ErrCode = "auth.invalid_credentials"
	ErrCodeUserNotSignedUp ErrCode = "auth.user_not_signed_up"
	ErrCodeInternalError ErrCode = "internal"
	ErrCodeInvalidDateTimeFormat ErrCode = "validation.datetime_format"
	ErrCodeCouldNotSaveProfile ErrCode = "profile.could_not_save"
)
