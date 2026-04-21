package config

type RedisConfig struct {
	Address string
}

func NewRedisConfig(address string) RedisConfig {
	return RedisConfig{Address: address}
}
