package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

type Token[T interface{}] struct {
	cfg       *JwtConfig
	token     string
	tokenType TokenType
}

func (token *Token[T]) Validate() (*ValidToken[T], error) {
	parsedToken, tokenParsingError := jwt.ParseWithClaims(token.token, &Claims[T]{}, func(t *jwt.Token) (any, error) {
		switch token.tokenType {
		case AccessToken:
			return []byte(token.cfg.AccessSecret), nil
		case RefreshToken:
			return []byte(token.cfg.RefreshSecret), nil
		}
		return nil, errors.New("can't determine type of the token")
	})

	if tokenParsingError != nil {
		return nil, tokenParsingError
	}

	validator := jwt.NewValidator(
		jwt.WithIssuer(token.cfg.Issuer),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
	)

	claims, ok := parsedToken.Claims.(*Claims[T])

	if !ok {
		return nil, errors.New("invalid claims type")
	}

	if validationError := validator.Validate(claims); validationError != nil {
		return nil, validationError
	}

	return &ValidToken[T]{
		token:     token.token,
		tokenType: token.tokenType,
		claims:    &claims.Data,
	}, nil
}

func CreateTokenFromHandlerAndString[T interface{}](jwtHandler IHandler[T], token string, tokenType TokenType) *Token[T] {
	return &Token[T]{
		cfg:       jwtHandler.getConfig(),
		token:     token,
		tokenType: tokenType,
	}
}
