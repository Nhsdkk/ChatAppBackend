package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type Claims[T interface{}] struct {
	Data T `json:"data"`
	jwt.RegisteredClaims
}

func (claims Claims[T]) appendMetadataToClaims(
	cfg *JwtConfig,
	tokenType TokenType,
) Claims[T] {
	var expireTimeoutString string

	switch tokenType {
	case RefreshToken:
		expireTimeoutString = cfg.ExpireTimeoutRefresh
	case AccessToken:
		expireTimeoutString = cfg.ExpireTimeoutAccess
	}

	expireTimeout, durationParseErr := time.ParseDuration(expireTimeoutString)
	if durationParseErr != nil {
		panic(fmt.Sprintf("can't parse duration from %s", cfg.ExpireTimeoutAccess))
	}

	claims.RegisteredClaims = jwt.RegisteredClaims{
		Issuer:    cfg.Issuer,
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(expireTimeout)},
		IssuedAt:  &jwt.NumericDate{Time: time.Now()},
		ID:        uuid.New().String(),
	}
	return claims
}

func CreateClaimsFromData[T interface{}](data T) Claims[T] {
	claims := Claims[T]{}
	claims.Data = data

	return claims
}
