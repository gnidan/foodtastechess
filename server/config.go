package server

type ServerConfig struct {
	BindAddress string
}

type AuthConfig struct {
	GoogleKey    string
	GoogleSecret string
	CallbackUrl  string
	SessionKey   string
}

type SessionConfig struct {
	SessionName string
	Secret      string
}

const authContextKey string = "user"
