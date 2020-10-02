package service

import (
	"github.com/Aoi-hosizora/ahlib/xproperty"
	"github.com/vidorg/vid_backend/lib/xgorm"
)

type OrderbyService struct {
	userDict  xproperty.PropertyDict
	videoDict xproperty.PropertyDict
}

func NewOrderbyService() *OrderbyService {
	return &OrderbyService{
		userDict: xproperty.PropertyDict{
			"nickname":      xproperty.NewValue(false, "nickname"),
			"register_time": xproperty.NewValue(false, "created_at"),
		},
		videoDict: xproperty.PropertyDict{
			"title":       xproperty.NewValue(false, "title"),
			"upload_time": xproperty.NewValue(false, "create_at"),
			"author_uid":  xproperty.NewValue(false, "author_uid"),
		},
	}
}

func (o *OrderbyService) User(s string) string {
	return xgorm.OrderByFunc(o.userDict)(s)
}

func (o *OrderbyService) Video(s string) string {
	return xgorm.OrderByFunc(o.videoDict)(s)
}
