package util

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/shomali11/util/xconditions"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/model/dto/common"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type authUtil struct{}

var AuthUtil = new(authUtil)

func (a *authUtil) EncryptPassword(password string) (string, error) {
	pwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(pwd), nil
}

func (a *authUtil) CheckPassword(password string, encrypted string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encrypted), []byte(password)) == nil
}

func (a *authUtil) GenerateToken(uid int, ex int64, config *config.JwtConfig) (string, error) {
	claims := &common.Claims{
		UserID: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + ex,
			Issuer:    config.Issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Secret))
}

func (a *authUtil) ParseToken(signedToken string, config *config.JwtConfig) (*common.Claims, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Secret), nil
	}
	token, err := jwt.ParseWithClaims(signedToken, &common.Claims{}, keyFunc)
	if err != nil || !token.Valid {
		err = xconditions.IfThenElse(err == nil, jwt.ValidationError{}, err).(error)
		return nil, err
	}

	claims, ok := token.Claims.(*common.Claims)
	if !ok {
		return nil, jwt.ValidationError{}
	}
	return claims, nil
}

func (a *authUtil) IsTokenExpireError(err error) bool {
	str := "token is expired by"
	return err.Error()[:len(str)] == str
}
