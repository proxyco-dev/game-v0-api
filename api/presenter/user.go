package presenter

type SignInRequest struct {
	Email    string `json:"email" column:"email"`
	Password string `json:"password"`
}

type SignInResponse struct {
	Token string `json:"token"`
}

type MeResponse struct {
	Email string `json:"email"`
}

type SignUpRequest struct {
	Username string `json:"username" column:"username" validate:"required,min=4,max=32"`
	Email    string `json:"email" column:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type SignUpResponse struct {
	Message string `json:"message"`
	User    struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Id       string `json:"id"`
	} `json:"user"`
}
