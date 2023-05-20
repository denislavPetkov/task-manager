package constants

import "time"

const (
	SessionCookieTtl = 60000
	SessionUserKey   = "user"
	SessionIdKey     = "id"
	SessionStoreName = "session"
	CookieUser       = "user"

	PasswordRecoveryTokenExpiration = time.Minute * 5

	CsrfKey = "csrf"
)
