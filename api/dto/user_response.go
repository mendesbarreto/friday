package dto

type UserResponseBody struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}


type AuthResponseBody struct {
	Token string `json:"token"`
}
