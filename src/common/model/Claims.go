package model

import "github.com/dgrijalva/jwt-go"

type UserClaims struct {
	UserId int32
	jwt.StandardClaims
}
