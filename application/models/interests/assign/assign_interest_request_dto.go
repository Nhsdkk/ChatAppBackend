package assign

import "chat_app_backend/internal/extensions"

type AssignInterestRequestDto struct {
	UserID      extensions.UUID   `json:"user_id" validator:"not_empty"`
	InterestIds []extensions.UUID `json:"interest_ids" validator:"not_empty"`
}
