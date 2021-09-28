package auth

var TokenFile string

const (
	EventAuthHostFlag = "event-auth-host"
	ApiAuthUrlFlag    = "api-auth-url"
)

func login() {
	Auth.Login()
}

func printAccessToken() {
	Auth.printAccessToken()
}
