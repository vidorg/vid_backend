package service

import (
	"github.com/Aoi-hosizora/ahlib/xproperty"
	"github.com/vidorg/vid_backend/lib/xgorm"
)

type OrderbyService struct {
	_userDict                xproperty.PropertyDict
	_followForUserDict       xproperty.PropertyDict
	_channelDict             xproperty.PropertyDict
	_subscribeForUserDict    xproperty.PropertyDict
	_subscribeForChannelDict xproperty.PropertyDict
	_videoDict               xproperty.PropertyDict
	_favoriteForVideoDict    xproperty.PropertyDict
	_favoriteForUserDict     xproperty.PropertyDict
}

func NewOrderbyService() *OrderbyService {
	return &OrderbyService{
		_userDict: xproperty.PropertyDict{
			"uid":           xproperty.NewValue(false, "uid"),
			"nickname":      xproperty.NewValue(false, "nickname"),
			"register_time": xproperty.NewValue(false, "created_at"),
		},
		_followForUserDict: xproperty.PropertyDict{
			"uid":           xproperty.NewValue(false, "uid"),
			"nickname":      xproperty.NewValue(false, "nickname"),
			"register_time": xproperty.NewValue(false, "created_at"),
			"follow_time":   xproperty.NewValue(false, "tbl_follow.created_at"),
		},
		_channelDict: xproperty.PropertyDict{
			"cid":         xproperty.NewValue(false, "cid"),
			"name":        xproperty.NewValue(false, "name"),
			"author_uid":  xproperty.NewValue(false, "author_uid"),
			"create_time": xproperty.NewValue(false, "created_at"),
		},
		_subscribeForUserDict: xproperty.PropertyDict{
			"uid":            xproperty.NewValue(false, "uid"),
			"nickname":       xproperty.NewValue(false, "nickname"),
			"register_time":  xproperty.NewValue(false, "created_at"),
			"subscribe_time": xproperty.NewValue(false, "tbl_subscribe.created_at"),
		},
		_subscribeForChannelDict: xproperty.PropertyDict{
			"cid":            xproperty.NewValue(false, "cid"),
			"name":           xproperty.NewValue(false, "name"),
			"author_uid":     xproperty.NewValue(false, "author_uid"),
			"subscribe_time": xproperty.NewValue(false, "tbl_subscribe.created_at"),
		},
		_videoDict: xproperty.PropertyDict{
			"vid":         xproperty.NewValue(false, "vid"),
			"title":       xproperty.NewValue(false, "title"),
			"upload_time": xproperty.NewValue(false, "created_at"),
			"author_uid":  xproperty.NewValue(false, "author_uid"),
		},
		_favoriteForUserDict: xproperty.PropertyDict{
			"uid":           xproperty.NewValue(false, "uid"),
			"nickname":      xproperty.NewValue(false, "nickname"),
			"register_time": xproperty.NewValue(false, "created_at"),
			"favorite_time": xproperty.NewValue(false, "tbl_favorite.created_at"),
		},
		_favoriteForVideoDict: xproperty.PropertyDict{
			"vid":           xproperty.NewValue(false, "vid"),
			"title":         xproperty.NewValue(false, "title"),
			"upload_time":   xproperty.NewValue(false, "created_at"),
			"author_uid":    xproperty.NewValue(false, "author_uid"),
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

// tbl_follow
func (o *OrderbyService) FollowForUser(s string) interface{} {
	order := xgorm.OrderByFunc(o._followForUserDict)(s)
	if order == "" {
		return "tbl_follow.created_at"
	}
	return order
}

// tbl_channel
func (o *OrderbyService) Channel(s string) string {
	order := xgorm.OrderByFunc(o._channelDict)(s)
	if order == "" {
		return "cid"
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

// tbl_subscribe
func (o *OrderbyService) SubscribeForChannel(s string) interface{} {
	order := xgorm.OrderByFunc(o._subscribeForChannelDict)(s)
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
func (o *OrderbyService) FavoriteForUser(s string) string {
	order := xgorm.OrderByFunc(o._favoriteForUserDict)(s)
	if order == "" {
		return "tbl_favorite.created_at"
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
