package dto

type CreateUserRequestBody struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=3,max=32"`
	Password string `json:"password" validate:"required,min=3,max=32"`
}

type GetUserByIdRequestParam struct {
	ID string `json:"id"`
}

type AuthenticateUserRequestBody struct {
	Username string `json:"username" validate:"required, email"`
	Password string `json:"password" validate:"required min=3,max=32"`
}
