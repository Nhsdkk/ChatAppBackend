package delete

import "github.com/google/uuid"

type DeleteUserRequestDto struct {
	ID uuid.UUID `json:"id" binder:"path,id" validator:"not_empty"`
}
