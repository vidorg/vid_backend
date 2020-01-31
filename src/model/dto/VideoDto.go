package dto

// @Model         _VideoDto
// @Description   视频信息
// @Property      vid         integer                            true * "视频id"
// @Property      title       string                             true * "视频标题"
// @Property      description string                             true * "标题简介"
// @Property      video_url   string                             true * "视频资源链接"
// @Property      cover_url   string                             true * "视频封面连接"
// @Property      upload_time string(format:2000-01-01 00:00:00) true * "视频上传时间"
// @Property      update_time string(format:2000-01-01 00:00:00) true * "视频修改时间"
// @Property      author      object(#_UserDto)                  true * "视频作者，用户id为-1表示用户已删除"
type VideoDto struct {
	Vid         int32    `json:"vid"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	VideoUrl    string   `json:"video_url"`
	CoverUrl    string   `json:"cover_url"`
	UploadTime  string   `json:"upload_time"`
	UpdateTime  string   `json:"update_time"`
	Author      *UserDto `json:"author"`
}
