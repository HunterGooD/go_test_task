package config

import "os"

type Config struct {
	// Host is host for api serving
	Host string
	// Port is port for api serving
	Port string
	// DSN is string for DB connection ${DBType}://${USERNAME}:${PASSWORD}@${HOST_DB}:{PORT_DB}}/${DB_NAME}?sslmode=disable
	DSN string
	// DBType is type db sqlite3 | postgres | mysql
	DBType string
}

func NewConfig() *Config {
	host := os.Getenv("HOST_API")
	port := os.Getenv("PORT_API")
	dsn := os.Getenv("DB_CONNECTION")
	DBtype := os.Getenv("DB_DRIVER")
	return &Config{host, port, dsn, DBtype}
}
