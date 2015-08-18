package config

import (
	"os"
)

type DatabaseConfig struct {
	HostAddr string
	Port     string
	Username string
	Password string
	Database string
}

func NewMariaDockerComposeConfig() DatabaseConfig {
	return DatabaseConfig{
		HostAddr: os.Getenv("MARIADB_PORT_3306_TCP_ADDR"),
		Port:     os.Getenv("MARIADB_PORT_3306_TCP_PORT"),
		Username: os.Getenv("MARIADB_ENV_MYSQL_USER"),
		Password: os.Getenv("MARIADB_ENV_MYSQL_PASSWORD"),
		Database: os.Getenv("MARIADB_ENV_MYSQL_DATABASE"),
	}
}
