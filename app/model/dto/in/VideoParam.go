package in

type VideoParam struct {
	Title       string `form:"title"       json:"title"       binding:"required"`
	Description string `form:"description" json:"description" binding:"required"`
	VideoUrl    string `form:"video_url"   json:"video_uil"   binding:"required"`
}
