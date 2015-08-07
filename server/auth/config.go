package auth

type AuthConfig struct {
	GoogleKey    string
	GoogleSecret string
	CallbackUrl  string
	SessionKey   string
}

const ContextKey string = "user"
