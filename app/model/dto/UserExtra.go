package dto

type UserExtraInfo struct {
	PhoneNumber      string `json:"phone_number,omitempty"`
	SubscribingCount int    `json:"subscribing_cnt"`
	SubscriberCount  int    `json:"subscriber_cnt"`
	VideoCount       int    `json:"video_cnt"`
	PlaylistCount    int    `json:"playlist_cnt"`
}
