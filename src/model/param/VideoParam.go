package param

// @Model         VideoParam
// @Description   视频请求参数
// @Property      title       string true "视频标题，长度在 [1, 100] 之间"
// @Property      description string true "视频简介，长度在 [0, 255] 之间" (empty:true)
// @Property      cover_url   string true "视频封面链接"
// @Property      video_url   string true "视频资源链接"
type VideoParam struct {
	Title       string  `form:"title"       json:"title"       binding:"required,min=1,max=100"`
	Description *string `form:"description" json:"description" binding:"required,min=0,max=255"`
	CoverUrl    string  `form:"cover_url"   json:"cover_url"   binding:"required,url"`
	VideoUrl    string  `form:"video_url"   json:"video_url"   binding:"required"` // TODO url
}
