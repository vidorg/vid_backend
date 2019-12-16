package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type passUtil struct{}

var PassUtil = new(passUtil)

var (
	JwtSecret []byte
	JwtExpire int64
)

type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

func (p *passUtil) MD5Check(content string, encrypted string) bool {
	return strings.EqualFold(p.MD5Encode(content), encrypted)
}

func (p *passUtil) MD5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func (p *passUtil) GenToken(uid int, ex int64) (string, error) {
	claims := Claims{
		uid,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + ex,
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

func (p *passUtil) ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtSecret, nil
	})

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
