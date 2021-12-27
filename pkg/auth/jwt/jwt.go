package jwtauth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTTokener struct {
	signinKey []byte
	tokenTTL  time.Duration
}

func New(signinKey string, tokenTTL time.Duration) *JWTTokener {
	return &JWTTokener{
		signinKey: []byte(signinKey),
		tokenTTL:  tokenTTL,
	}
}

func (j *JWTTokener) GenerateToken(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Second * j.tokenTTL).Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   userId,
	})

	fmt.Printf("SIGNING KEY - '%s'\n", j.signinKey)
	tokenString, err := token.SignedString(j.signinKey)
	if err != nil {
		return "", fmt.Errorf("error with signing token - %s", err.Error())
	}

	fmt.Println("Token: success generated, Token - ", tokenString)
	return tokenString, nil
}
