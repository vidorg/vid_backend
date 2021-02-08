package jwt

import (
	"github.com/kataras/jwt"
	"time"
)

// sharedKey
var (
	sharedKey = []byte("sercrethatmaycontainch@r$32chars")
	issuer    = "seefs"
)

func SetMeta(secret, publisher string) {
	sharedKey = []byte(secret)
	issuer = publisher
}

// UserClaims ...
type UserClaims struct {
	UID int `json:"uid"`
}

// GenerateToken generate token by userID
func GenerateToken(uid int, expire time.Duration) ([]byte, error) {
	userClaims := UserClaims{
		UID: uid,
	}

	now := time.Now()
	standardClaims := jwt.Claims{
		Expiry:   now.Add(expire).Unix(),
		IssuedAt: now.Unix(),
		Issuer:   issuer,
	}

	token, err := jwt.Sign(jwt.HS256, sharedKey, userClaims, standardClaims)
	if err != nil {
		return nil, err
	}
	return token, err
}

// GenerateTokenWithoutExpire generate token by userID without expire time
func GenerateTokenWithoutExpire(uid int) ([]byte, error) {
	userClaims := UserClaims{
		UID: uid,
	}

	token, err := jwt.Sign(jwt.HS256, sharedKey, userClaims)
	if err != nil {
		return nil, err
	}

	return token, err
}

// ParseToken parse token
func ParseToken(token []byte) (UserClaims, error) {
	verify, err := jwt.Verify(jwt.HS256, sharedKey, token)
	if err != nil {
		return UserClaims{}, err
	}

	var claims UserClaims
	err = verify.Claims(&claims)

	return claims, err
}
