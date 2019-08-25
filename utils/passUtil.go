package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type PassUtil struct{}

var JwtSecret []byte
var JwtTokenExpire int64

type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

func (p *PassUtil) MD5Check(content string, encrypted string) bool {
	return strings.EqualFold(p.MD5Encode(content), encrypted)
}

func (p *PassUtil) MD5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func (p *PassUtil) GenToken(id int) (string, error) {
	claims := Claims{
		id,
		jwt.StandardClaims{
			ExpiresAt: int64(time.Now().Unix() + JwtTokenExpire),
			Issuer:    "",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(JwtSecret)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Bearer %s", tokenStr), nil
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
