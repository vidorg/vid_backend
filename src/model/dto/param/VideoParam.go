package param

type VideoParam struct {
	Title       string  `form:"title"       json:"title"       binding:"required,min=1,max=100"`
	Description *string `form:"description" json:"description" binding:"required,min=0,max=1024"`
	VideoUrl    string  `form:"video_url"   json:"video_uil"   binding:"required"`
}
