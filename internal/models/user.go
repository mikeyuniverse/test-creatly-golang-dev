package models

type UserSignUpInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserSignInInput struct {
	Email        string `json:"email"`
	PasswordHash string `json:"password"`
}
