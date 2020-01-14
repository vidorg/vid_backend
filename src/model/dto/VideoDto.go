package dto

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
