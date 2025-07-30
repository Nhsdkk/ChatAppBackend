package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

type IHandler[T interface{}] interface {
	GenerateJwtPair(data T) (accessToken *ValidToken[T], refreshToken *ValidToken[T], err error)
	getConfig() *JwtConfig
	generateSingleToken(claims *Claims[T], tokenType TokenType) (*ValidToken[T], error)
}

type Handler[T interface{}] struct {
	cfg *JwtConfig
}

func (handler *Handler[T]) generateSingleToken(claims *Claims[T], tokenType TokenType) (*ValidToken[T], error) {
	var secret []byte

	switch tokenType {
	case RefreshToken:
		secret = []byte(handler.cfg.RefreshSecret)
	case AccessToken:
		secret = []byte(handler.cfg.AccessSecret)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims.appendMetadataToClaims(handler.cfg))
	tokenString, generationError := token.SignedString(secret)
	if generationError != nil {
		return nil, generationError
	}
	return &ValidToken[T]{tokenString, tokenType, &claims.Data}, nil
}

func (handler *Handler[T]) getConfig() *JwtConfig {
	return handler.cfg
}

func (handler *Handler[T]) GenerateJwtPair(data T) (*ValidToken[T], *ValidToken[T], error) {
	claims := CreateClaimsFromData(data)

	accessToken, accessTokenGenerationError := handler.generateSingleToken(claims, AccessToken)
	if accessTokenGenerationError != nil {
		return nil, nil, accessTokenGenerationError
	}

	refreshToken, refreshTokenGenerationError := handler.generateSingleToken(claims, RefreshToken)
	if refreshTokenGenerationError != nil {
		return nil, nil, refreshTokenGenerationError
	}

	return accessToken, refreshToken, nil
}

func CreateJwtHandler[T interface{}](config *JwtConfig) (handler *Handler[T], error error) {
	handler = &Handler[T]{}
	handler.cfg = config
	return handler, nil
}
