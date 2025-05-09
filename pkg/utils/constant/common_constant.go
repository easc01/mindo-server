package constant

import "time"

const (
	AppName        = "Mindo2.0"
	Blank          = ""
	IdParam        = "/:id"
	Authorization  = "Authorization"
	Week           = 7 * 24 * time.Hour
	Month          = 30 * 24 * time.Hour
	RefreshToken   = "RefreshToken"
	UserContextKey = "userContext"
	TimeLayout     = "2006-01-02T15:04:05.999999Z"
)
