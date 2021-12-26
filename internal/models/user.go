package models

type UserSignUpInput struct {
	Email    string `json:"email",bson:"email"`
	Password string `json:"password",bson:"password"`
}

type UserSignInInput struct {
	Email        string `json:"email"`
	PasswordHash string `json:"password"`
}
