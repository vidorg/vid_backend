package utils

import (
	jwt "github.com/dgrijalva/jwt-go"
)

type PassUtil struct{}

var JwtSecret []byte
var JwtTokenExpire int64

type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

func (p *PassUtil) GenToken(id int) (string, error) {
	claims := Claims{
		id,
		jwt.StandardClaims{
			ExpiresAt: JwtTokenExpire,
			Issuer:    "",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(JwtSecret)

	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (p *PassUtil) ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtSecret, nil
	})

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
