package jwtauth

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTTokener struct {
	signinKey string
	tokenTTL  time.Duration
}

func New(signinKey string, tokenTTL time.Duration) *JWTTokener {
	return &JWTTokener{
		signinKey: signinKey,
		tokenTTL:  tokenTTL,
	}
}

func (j *JWTTokener) GenerateToken(userId int) (string, error) {
	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	claims := make(jwt.MapClaims)
	claims["userID"] = userId
	claims["exp"] = time.Now().Add(time.Second * j.tokenTTL).Unix()
	token.Claims = claims
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(j.signinKey)
	return tokenString, err
}
