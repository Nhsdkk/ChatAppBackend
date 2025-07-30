package jwt

type JwtConfig struct {
	AccessSecret  string `env:"access_secret"`
	RefreshSecret string `env:"refresh_secret"`
	Issuer        string `env:"issuer"`
	ExpireTimeout string `env:"expire_timeout"`
}
