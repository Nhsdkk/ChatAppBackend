package jwt

type ValidToken[T interface{}] struct {
	token     string
	tokenType TokenType
	claims    *T
}

func (token *ValidToken[T]) GetClaims() *T {
	return token.claims
}

func (token *ValidToken[T]) GetToken() string {
	return token.token
}

func (token *ValidToken[T]) RefreshRelatedAccessToken(handler IHandler[T]) (*ValidToken[T], error) {
	claims := CreateClaimsFromData(*token.claims).
		appendMetadataToClaims(handler.getConfig(), AccessToken)

	switch token.tokenType {
	case AccessToken:
		panic("can't refresh access token using access token")
	case RefreshToken:
		return handler.generateSingleToken(&claims, AccessToken)
	}

	panic("Unreachable")
}
