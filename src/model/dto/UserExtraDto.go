package dto

type UserExtraInfo struct {
	SubscribingCount int `json:"subscribing_cnt"`
	SubscriberCount  int `json:"subscriber_cnt"`
	VideoCount       int `json:"video_cnt"`
}
