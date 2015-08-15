package queries

import (
	"os"
)

type QueriesCacheConfig struct {
	HostAddr string
	Port     string
	Database string
}

func NewMongoDockerComposeConfig() QueriesCacheConfig {
	return QueriesCacheConfig{
		HostAddr: os.Getenv("MONGODB_PORT_27017_TCP_ADDR"),
		Port:     os.Getenv("MONGODB_PORT_27017_TCP_PORT"),
		Database: "test",
	}
}
