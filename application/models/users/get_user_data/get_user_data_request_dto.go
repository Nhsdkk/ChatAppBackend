package get_user_data

import "github.com/google/uuid"

type GetUserDataRequestDto struct {
	ID uuid.UUID `json:"id" binder:"path,id" validator:"not_empty"`
}
