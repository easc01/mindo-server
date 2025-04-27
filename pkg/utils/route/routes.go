package route

const (
	Api       = "/api"
	User      = "/users"
	Admin     = "/admins"
	Auth      = "/auth"
	Interest  = "/interest"
	Refresh   = "/refresh"
	Google    = "/google"
	Playlists = "/playlists"
	SignIn    = "/sign-in"
	SignUp    = "/sign-up"
)

func GetRefreshRoute() string {
	return Api + Auth + Refresh
}
