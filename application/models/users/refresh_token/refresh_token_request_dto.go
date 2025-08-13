package refresh_token

type RefreshTokenRequestDto struct {
	RefreshToken string `json:"refresh_token" validator:"not_empty"`
}
