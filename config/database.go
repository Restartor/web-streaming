package config

type DatabaseConfig struct {
	DSN string
}

func NewDatabaseConfig(dsn string) DatabaseConfig {
	return DatabaseConfig{DSN: dsn}
}
