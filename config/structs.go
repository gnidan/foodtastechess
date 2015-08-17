package config

type QueriesCacheConfig struct {
	HostAddr string
	Port     string
	Database string
}

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

type DatabaseConfig struct {
	HostAddr string
	Port     string
	Username string
	Password string
	Database string
}
