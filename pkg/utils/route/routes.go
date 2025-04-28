package route

const (
	Api       = "/api"
	User      = "/users"
	Admin     = "/admins"
	Auth      = "/auth"
	Interest  = "/interests"
	Refresh   = "/refresh"
	Google    = "/google"
	Playlists = "/playlists"
	Topics    = "/topics"
	SignIn    = "/sign-in"
	SignUp    = "/sign-up"
)

func GetRefreshRoute() string {
	return Api + Auth + Refresh
}
