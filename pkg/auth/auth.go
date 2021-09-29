package auth

var TokenFile string

const (
	EventAuthHostFlag = "event-auth-host"
	ApiAuthHostFlag   = "api-auth-host"
)

func login() {
	Auth.Login()
}

func printAccessToken() {
	Auth.printAccessToken()
}
