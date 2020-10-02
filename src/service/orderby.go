package service

import (
	"github.com/Aoi-hosizora/ahlib/xproperty"
	"github.com/vidorg/vid_backend/lib/xgorm"
)

type OrderbyService struct {
	_userDict             xproperty.PropertyDict
	_subscribeForUserDict xproperty.PropertyDict
	_videoDict            xproperty.PropertyDict
	_favoriteForVideoDict xproperty.PropertyDict
	_favoriteForUserDict  xproperty.PropertyDict
}

func NewOrderbyService() *OrderbyService {
	return &OrderbyService{
		_userDict: xproperty.PropertyDict{
			"uid":           xproperty.NewValue(false, "uid"),
			"nickname":      xproperty.NewValue(false, "nickname"),
			"register_time": xproperty.NewValue(false, "created_at"),
		},
		_subscribeForUserDict: xproperty.PropertyDict{
			"uid":            xproperty.NewValue(false, "uid"),
			"nickname":       xproperty.NewValue(false, "nickname"),
			"register_time":  xproperty.NewValue(false, "created_at"),
			"subscribe_time": xproperty.NewValue(false, "tbl_subscribe.created_at"),
		},
		_videoDict: xproperty.PropertyDict{
			"vid":         xproperty.NewValue(false, "vid"),
			"title":       xproperty.NewValue(false, "title"),
			"upload_time": xproperty.NewValue(false, "create_at"),
			"author_uid":  xproperty.NewValue(false, "author_uid"),
		},
		_favoriteForVideoDict: xproperty.PropertyDict{
			"vid":           xproperty.NewValue(false, "vid"),
			"title":         xproperty.NewValue(false, "title"),
			"upload_time":   xproperty.NewValue(false, "create_at"),
			"author_uid":    xproperty.NewValue(false, "author_uid"),
			"favorite_time": xproperty.NewValue(false, "tbl_favorite.created_at"),
		},
		_favoriteForUserDict: xproperty.PropertyDict{
			"uid":           xproperty.NewValue(false, "uid"),
			"nickname":      xproperty.NewValue(false, "nickname"),
			"register_time": xproperty.NewValue(false, "created_at"),
			"favorite_time": xproperty.NewValue(false, "tbl_favorite.created_at"),
		},
	}
}

// tbl_user
func (o *OrderbyService) User(s string) interface{} {
	order := xgorm.OrderByFunc(o._userDict)(s)
	if order == "" {
		return "uid"
	}
	return order
}

// tbl_subscribe
func (o *OrderbyService) SubscribeForUser(s string) interface{} {
	order := xgorm.OrderByFunc(o._subscribeForUserDict)(s)
	if order == "" {
		return "tbl_subscribe.created_at"
	}
	return order
}

// tbl_video
func (o *OrderbyService) Video(s string) string {
	order := xgorm.OrderByFunc(o._videoDict)(s)
	if order == "" {
		return "vid"
	}
	return order
}

// tbl_favorite
func (o *OrderbyService) FavoriteForVideo(s string) string {
	order := xgorm.OrderByFunc(o._favoriteForVideoDict)(s)
	if order == "" {
		return "tbl_favorite.created_at"
	}
	return order
}

// tbl_favorite
func (o *OrderbyService) FavoriteForUser(s string) string {
	order := xgorm.OrderByFunc(o._favoriteForUserDict)(s)
	if order == "" {
		return "tbl_favorite.created_at"
	}
	return order
}
