package refresh_token

type RefreshTokenRequestDto struct {
	RefreshToken string `json:"refresh_token" binder:"body,refresh_token" validator:"not_empty"`
}
