package service

import (
	"github.com/Aoi-hosizora/ahlib-more/xjwt"
	"github.com/Aoi-hosizora/ahlib/xdi"
	"github.com/Aoi-hosizora/ahlib/xstatus"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/vidorg/vid_backend/src/common/exception"
	"github.com/vidorg/vid_backend/src/common/result"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/model/po"
	"github.com/vidorg/vid_backend/src/provide/sn"
	"time"
)

type JwtService struct {
	config       *config.Config
	userService  *UserService
	tokenService *TokenService
}

func NewJwtService() *JwtService {
	return &JwtService{
		config:       xdi.GetByNameForce(sn.SConfig).(*config.Config),
		userService:  xdi.GetByNameForce(sn.SUserService).(*UserService),
		tokenService: xdi.GetByNameForce(sn.STokenService).(*TokenService),
	}
}

func (j *JwtService) UserKey() string {
	return "user"
}

func (j *JwtService) GetToken(c *gin.Context) string {
	token := c.GetHeader("Authorization")
	if token != "" {
		return token
	}
	token = c.DefaultQuery("Authorization", "")
	if token != "" {
		return token
	}
	return c.DefaultQuery("token", "")
}

func (j *JwtService) JwtCheck(token string) (*po.User, xstatus.JwtStatus, error) {
	if token == "" {
		return nil, xstatus.JwtBlank, nil
	}

	// parse
	uid, err := j.ParseToken(token)
	if err != nil {
		if xjwt.TokenExpired(err) {
			return nil, xstatus.JwtExpired, nil
		} else if xjwt.TokenIssuerInvalid(err) {
			return nil, xstatus.JwtIssuer, nil
		}
		return nil, xstatus.JwtInvalid, nil
	}

	// redis
	ok, err := j.tokenService.Query(token)
	if err != nil {
		return nil, xstatus.JwtFailed, err
	} else if !ok {
		return nil, xstatus.JwtNotFound, nil
	}

	// mysql
	user, err := j.userService.QueryByUid(uid)
	if err != nil {
		return nil, xstatus.JwtFailed, err
	} else if user == nil {
		return nil, xstatus.JwtUserErr, nil
	}

	return user, xstatus.JwtSuccess, nil
}

func (j *JwtService) GetContextUser(c *gin.Context) *po.User {
	itf, exist := c.Get(j.UserKey())
	if exist {
		user, ok := itf.(*po.User)
		if !ok {
			result.Error(exception.UnAuthorizedError).JSON(c)
			c.Abort()
			return nil
		}
		return user
	}

	token := j.GetToken(c)
	user, status, err := j.JwtCheck(token)
	if err != nil || status != xstatus.JwtSuccess {
		return nil
	}

	c.Set(j.UserKey(), user)
	return user
}

type userClaims struct {
	Uid uint64
	jwt.StandardClaims
}

func (j *JwtService) GenerateToken(uid uint64) (string, error) {
	claims := &userClaims{
		Uid: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + j.config.Jwt.Expire,
			Issuer:    j.config.Jwt.Issuer,
		},
	}
	token, err := xjwt.GenerateToken(claims, []byte(j.config.Jwt.Secret))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (j *JwtService) ParseToken(signedToken string) (uint64, error) {
	claims, err := xjwt.ParseToken(signedToken, []byte(j.config.Jwt.Secret), &userClaims{})
	if err != nil {
		return 0, err
	}

	c, ok := claims.(*userClaims)
	if !ok {
		return 0, xjwt.DefaultValidationError
	} else if c.Issuer != j.config.Jwt.Issuer {
		return 0, jwt.NewValidationError("unknown token issuer", jwt.ValidationErrorIssuer)
	}

	return c.Uid, nil
}
