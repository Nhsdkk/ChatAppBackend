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

func (claims *Claims[T]) appendMetadataToClaims(cfg *JwtConfig) *Claims[T] {
	expireTimeout, durationParseErr := time.ParseDuration(cfg.ExpireTimeout)
	if durationParseErr != nil {
		panic(fmt.Sprintf("can't parse duration from %s", cfg.ExpireTimeout))
	}

	claims.RegisteredClaims = jwt.RegisteredClaims{
		Issuer:    cfg.Issuer,
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(expireTimeout)},
		IssuedAt:  &jwt.NumericDate{Time: time.Now()},
		ID:        uuid.New().String(),
	}
	return claims
}

func CreateClaimsFromData[T interface{}](data T) (claims *Claims[T]) {
	claims = &Claims[T]{}
	claims.Data = data

	return claims
}
