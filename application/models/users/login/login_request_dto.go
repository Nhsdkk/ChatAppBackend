package login

type LoginRequestDto struct {
	Email    string `validator:"not_empty" json:"email" binder:"body,email"`
	Password string `validator:"not_empty" json:"password" binder:"body,password"`
}
