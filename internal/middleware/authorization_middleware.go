package middleware

import (
	"chat_app_backend/application/models/jwt_claims"
	"chat_app_backend/internal/exceptions"
	"chat_app_backend/internal/exceptions/common_exceptions"
	"chat_app_backend/internal/jwt"
	"chat_app_backend/internal/sqlc/db"
	"github.com/gin-gonic/gin"
	"regexp"
)

const ClaimsKey = "Claims"

var authorizationHeaderRegexp = regexp.MustCompile("Bearer (?P<token>\\S+)")

func AuthorizationMiddleware(jwtHandler jwt.IHandler[jwt_claims.UserClaims], db db.IDbConnection) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader, exists := ctx.Request.Header["Authorization"]

		if !exists || !authorizationHeaderRegexp.MatchString(authorizationHeader[0]) {
			_ = ctx.Error(
				common_exceptions.UnauthorizedException{
					BaseRestException: exceptions.BaseRestException{
						ITrackableException: exceptions.CreateTrackableExceptionFromStringF(
							"access token format does not match",
						),
						Message: "",
					},
				},
			)
			ctx.Next()
			return
		}

		groupIdx := authorizationHeaderRegexp.SubexpIndex("token")
		matches := authorizationHeaderRegexp.FindStringSubmatch(authorizationHeader[0])
		if groupIdx >= len(matches) {
			_ = ctx.Error(
				exceptions.CreateTrackableExceptionFromStringF(
					"access token matched, but token group not",
				),
			)
			ctx.Next()
			return
		}

		token := jwt.CreateTokenFromHandlerAndString(jwtHandler, matches[groupIdx], jwt.AccessToken)
		validToken, validationError := token.Validate()
		if validationError != nil {
			_ = ctx.Error(
				common_exceptions.UnauthorizedException{
					BaseRestException: exceptions.BaseRestException{
						ITrackableException: exceptions.WrapErrorWithTrackableException(validationError),
						Message:             "",
					},
				},
			)
			ctx.Next()
			return
		}

		user, userExistenceError := db.GetQueries().GetUserById(ctx, validToken.GetClaims().ID)

		if userExistenceError != nil || !validToken.GetClaims().Equals(&user) {
			_ = ctx.Error(
				common_exceptions.UnauthorizedException{
					BaseRestException: exceptions.BaseRestException{
						ITrackableException: exceptions.CreateTrackableExceptionFromStringF(
							"user with id: %s no longer exists or related data does not match",
							user.ID,
						),
						Message: "",
					},
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
