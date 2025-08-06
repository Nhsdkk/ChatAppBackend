package redis

type RedisConfig struct {
	Host       string `env:"HOST"`
	Port       int    `env:"PORT"`
	Password   string `env:"PASSWORD"`
	Username   string `env:"USERNAME"`
	ClientName string `env:"CLIENT_NAME"`
	DB         int    `env:"DB"`
}
