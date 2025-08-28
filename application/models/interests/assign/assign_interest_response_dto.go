package assign

import "chat_app_backend/application/models/interests/get"

type AssignInterestResponseDto struct {
	Interests []get.GetInterestResponseDto `json:"interests"`
}
