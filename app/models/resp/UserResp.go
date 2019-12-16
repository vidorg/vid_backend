package resp

import (
	po2 "vid/app/models/po"
)

type UserResp struct {
	User po2.User      `json:"user"`
	Info UserExtraInfo `json:"info"`
}

type UserExtraInfo struct {
	PhoneNumber     int `json:"phone_number,omitempty"`
	Subscriber_cnt  int `json:"subscriber_cnt"`
	Subscribing_cnt int `json:"subscribing_cnt"`
	Video_cnt       int `json:"video_cnt"`
	Playlist_cnt    int `json:"playlist_cnt"`
}
