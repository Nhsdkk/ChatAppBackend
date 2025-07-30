package middleware

import (
	"chat_app_backend/application/models/exception"
	"chat_app_backend/application/models/jwt_claims"
	"chat_app_backend/internal/jwt"
	"chat_app_backend/internal/sqlc/db"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"regexp"
)

const ClaimsKey = "Claims"

var authorizationHeaderRegexp = regexp.MustCompile("Bearer (?P<token>\\S+)")

func AuthorizationMiddleware(jwtHandler jwt.IHandler[jwt_claims.UserClaims], db db.IDbConnection) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader, exists := ctx.Request.Header["Authorization"]

		if !exists || !authorizationHeaderRegexp.MatchString(authorizationHeader[0]) {
			_ = ctx.Error(exception.UnauthorizedException{Err: errors.New("access token format does not match")})
			ctx.Next()
			return
		}

		groupIdx := authorizationHeaderRegexp.SubexpIndex("token")
		matches := authorizationHeaderRegexp.FindStringSubmatch(authorizationHeader[0])
		if groupIdx >= len(matches) {
			_ = ctx.Error(exception.ServerException{Err: errors.New("access token matched, but token group not")})
			ctx.Next()
			return
		}

		token := jwt.CreateTokenFromHandlerAndString(jwtHandler, matches[groupIdx], jwt.AccessToken)
		validToken, validationError := token.Validate()
		if validationError != nil {
			_ = ctx.Error(exception.UnauthorizedException{Err: validationError})
			ctx.Next()
			return
		}

		user, userExistenceError := db.GetQueries().GetUserById(ctx, validToken.GetClaims().ID)

		if userExistenceError != nil || !validToken.GetClaims().Equals(&user) {
			_ = ctx.Error(
				exception.UnauthorizedException{
					Err: errors.New(
						fmt.Sprintf(
							"user with id: %s no longer exists or related data does not match",
							user.ID,
						),
					),
				},
			)
			ctx.Next()
			return
		}

		ctx.Set(ClaimsKey, &user)
		ctx.Next()
		return
	}
}
