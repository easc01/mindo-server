package route

const (
	Api       = "/api"
	User      = "/users"
	Admin     = "/admins"
	Auth      = "/auth"
	Refresh   = "/refresh"
	Google    = "/google"
	Playlists = "/playlists"
)

func GetRefreshRoute() string {
	return Auth + Refresh
}
