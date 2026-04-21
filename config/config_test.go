package config

import "testing"

func TestNewDatabaseConfig(t *testing.T) {
cfg := NewDatabaseConfig("postgres://localhost/web")
if cfg.DSN != "postgres://localhost/web" {
t.Fatalf("unexpected DSN: %s", cfg.DSN)
}
}

func TestNewRedisConfig(t *testing.T) {
cfg := NewRedisConfig("127.0.0.1:6379")
if cfg.Address != "127.0.0.1:6379" {
t.Fatalf("unexpected redis address: %s", cfg.Address)
}
}
