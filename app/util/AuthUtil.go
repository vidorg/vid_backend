package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
	"vid/app/middleware"

	"github.com/dgrijalva/jwt-go"
)

type authUtil struct{}

var AuthUtil = new(authUtil)

func (a *authUtil) MD5Check(content string, encrypted string) bool {
	return strings.EqualFold(a.MD5Encode(content), encrypted)
}

func (a *authUtil) MD5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func (a *authUtil) GenerateToken(uid int, ex int64) (string, error) {
	claims := middleware.Claims{
		UserID: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + ex,
			Issuer:    "",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(middleware.JwtSecret)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Bearer %s", tokenStr), nil
}

func (a *authUtil) ParseToken(tokenStr string) (*middleware.Claims, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return middleware.JwtSecret, nil
	}
	token, err := jwt.ParseWithClaims(tokenStr, &middleware.Claims{}, keyFunc)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*middleware.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ValidationError{}
}
