package jwt

type JwtConfig struct {
	AccessSecret         string `env:"access_secret"`
	RefreshSecret        string `env:"refresh_secret"`
	Issuer               string `env:"issuer"`
	ExpireTimeoutAccess  string `env:"expire_timeout_access"`
	ExpireTimeoutRefresh string `env:"expire_timeout_refresh"`
}
