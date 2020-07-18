package util

import (
	"github.com/Aoi-hosizora/ahlib/xcondition"
	"github.com/dgrijalva/jwt-go"
	"github.com/vidorg/vid_backend/src/config"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type _AuthUtil struct{}

var AuthUtil = &_AuthUtil{}

type userClaims struct {
	UserId int32
	jwt.StandardClaims
}

func (a *_AuthUtil) EncryptPassword(password string) (string, error) {
	pwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(pwd), nil
}

func (a *_AuthUtil) CheckPassword(password string, encrypted string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encrypted), []byte(password)) == nil
}

func (a *_AuthUtil) GenerateToken(uid int32, ex int64, config *config.JwtConfig) (string, error) {
	claims := &userClaims{
		UserId: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + ex,
			Issuer:    config.Issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(config.Secret))
	if err != nil {
		return "", err
	}
	return "Bearer " + t, nil
}

func (a *_AuthUtil) ParseToken(signedToken string, config *config.JwtConfig) (*userClaims, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Secret), nil
	}
	signedToken = strings.TrimPrefix(signedToken, "Bearer ")
	token, err := jwt.ParseWithClaims(signedToken, &userClaims{}, keyFunc)
	if err != nil || !token.Valid {
		err = xcondition.IfThenElse(err == nil, jwt.ValidationError{}, err).(error)
		return nil, err
	}

	claims, ok := token.Claims.(*userClaims)
	if !ok {
		return nil, jwt.ValidationError{}
	}
	return claims, nil
}

func (a *_AuthUtil) IsTokenExpireError(err error) bool {
	str := "token is expired by"
	return err.Error()[:len(str)] == str
}
