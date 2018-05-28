package backups

const baseURL = "https://cloudapi.acquia.com/v1"

// Client for interacting with Acquias CloudAPI.
type Client struct {
	site string
	env  string
	db   string
	user string
	pass string
}

// New client for interacting with Acquias CloudAPI.
func New(user, pass, site, env, db string) Client {
	return Client{
		site: site,
		env:  env,
		db:   db,
		user: user,
		pass: pass,
	}
}
