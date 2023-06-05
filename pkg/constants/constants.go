package constants

import "time"

const (
	SessionCookieTtl = 3600
	SessionUserKey   = "user"
	SessionIdKey     = "id"
	SessionStoreName = "session"

	CookieUser = "user"

	PasswordRecoveryTokenExpiration = time.Minute * 5

	CsrfKey = "csrf"
)
