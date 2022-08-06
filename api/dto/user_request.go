package dto

type CreateUserRequestBody struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type GetUserByIdRequestParam struct {
	ID string `json:"id"`
}
