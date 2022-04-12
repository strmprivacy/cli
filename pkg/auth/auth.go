package auth

var TokenFile string

const (
	EventsAuthUrlFlag = "events-auth-url"
	ApiAuthUrlFlag    = "api-auth-url"
)

func login() {
	Auth.login()
}

func revoke() {
	Auth.revoke()
}

func printAccessToken() {
	Auth.printAccessToken()
}
