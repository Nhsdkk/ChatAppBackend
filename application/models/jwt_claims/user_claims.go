package jwt_claims

import (
	"chat_app_backend/internal/extensions"
	"chat_app_backend/internal/sqlc/db_queries"
)

type UserClaims struct {
	ID            extensions.UUID     `json:"id"`
	FullName      string              `json:"full_name"`
	Email         string              `json:"email"`
	Role          db_queries.RoleType `json:"role"`
	EmailVerified bool                `json:"email_verified"`
}

func (uc *UserClaims) Equals(user *db_queries.User) bool {
	return uc.ID == user.ID &&
		uc.FullName == user.FullName &&
		uc.Email == user.Email &&
		uc.EmailVerified == user.EmailVerified &&
		uc.Role == user.Role
}
