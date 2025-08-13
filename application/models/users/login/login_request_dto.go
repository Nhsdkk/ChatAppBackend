package login

type LoginRequestDto struct {
	Email    string `validator:"not_empty" json:"email"`
	Password string `validator:"not_empty" json:"password"`
}
