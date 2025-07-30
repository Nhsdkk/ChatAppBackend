package jwt

type TokenType int

const (
	AccessToken TokenType = iota
	RefreshToken
)
