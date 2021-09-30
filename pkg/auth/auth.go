package auth

var TokenFile string

const (
	EventsAuthUrlFlag = "events-auth-url"
	ApiAuthUrlFlag    = "api-auth-url"
)

func login() {
	Auth.Login()
}

func printAccessToken() {
	Auth.printAccessToken()
}
