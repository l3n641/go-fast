package tools

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

type JWT struct {
	SigningKey []byte
}

func (j *JWT) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

func (j *JWT) ParseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.SigningKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
