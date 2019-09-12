package resp

import (
	. "vid/models"
)

type UserResp struct {
	User User          `json:"user"`
	Info UserExtraInfo `json:"info"`
}

type UserExtraInfo struct {
	Subscriber_cnt  int `json:"subscriber_cnt"`
	Subscribing_cnt int `json:"subscribing_cnt"`
	Video_cnt       int `json:"video_cnt"`
}
