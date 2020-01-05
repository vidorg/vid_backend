package param

type VideoParam struct {
	Title       string `form:"title"       json:"title"       binding:"required,gte=5,lte=100"`
	Description string `form:"description" json:"description" binding:"required,gte=0,lte=255"`
	VideoUrl    string `form:"video_url"   json:"video_uil"   binding:"required"`
}
