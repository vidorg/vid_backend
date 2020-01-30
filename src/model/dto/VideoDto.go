package dto

// @Model      VideoDtoResult "返回视频信息"
// @Property   code    integer            true false "返回响应码"
// @Property   message string             true false "返回信息"
// @Property   data    object(#_VideoDto) true false "返回数据"

// @Model      VideoDtoPageResult "返回视频分页信息"
// @Property   code      integer                true false "返回响应码"
// @Property   message   string                 true false "返回信息"
// @Property   data      object(#_VideoDtoPage) true false "返回数据"

// @Model      _VideoDtoPage "视频分页信息"
// @Property   total   integer           true false "数据总量"
// @Property   message string            true false "当前页"
// @Property   data    array(#_VideoDto) true false "返回数据"

// @Model      _VideoDto "视频信息"
// @Property   vid         integer           true false "视频id"
// @Property   title       string            true false "视频标题"
// @Property   description string            true false "标题简介"
// @Property   video_url   string            true false "视频资源链接"
// @Property   cover_url   string            true false "视频封面连接"
// @Property   upload_time string            true false "视频上传时间，固定格式为 2000-01-01 00:00:00"
// @Property   update_time string            true false "视频修改时间，固定格式为 2000-01-01 00:00:00"
// @Property   author      object(#_UserDto) true false "视频作者，用户id为-1表示用户已删除"
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
