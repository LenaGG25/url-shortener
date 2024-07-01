package config

// Config is config entity
type Config struct {
	MaxShortKeySize int
	Port            int
	PostgresHost    string
	RedisHost       string
	Expiration      string
}
