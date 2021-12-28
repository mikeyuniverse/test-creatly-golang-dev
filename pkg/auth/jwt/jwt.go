package jwtauth

import (
	"creatly-task/internal/config"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTTokener struct {
	signinKey []byte
	tokenTTL  time.Duration
}

func New(config *config.JWT) *JWTTokener {
	return &JWTTokener{
		signinKey: []byte(config.SigninKey),
		tokenTTL:  time.Second * time.Duration(config.TokenTTL),
	}
}

func (j *JWTTokener) GenerateToken(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Second * j.tokenTTL).Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   userId,
	})

	tokenString, err := token.SignedString(j.signinKey)
	if err != nil {
		return "", fmt.Errorf("error with signing token - %s", err.Error())
	}

	return tokenString, nil
}

func (j *JWTTokener) ParseToken(token string) (string, error) {
	acceptedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.signinKey, nil
	})

	if err != nil {
		return "", err
	}

	if !acceptedToken.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := acceptedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims")
	}

	subject, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("invalid claims - subject")
	}

	return subject, nil
}
