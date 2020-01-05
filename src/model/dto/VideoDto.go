package dto

import (
	"fmt"
	"github.com/vidorg/vid_backend/src/model/dto/common"
	"github.com/vidorg/vid_backend/src/model/enum"
	"github.com/vidorg/vid_backend/src/model/po"
	"strings"
)

type VideoDto struct {
	Vid         int      `json:"vid"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	VideoUrl    string   `json:"video_url"`
	CoverUrl    string   `json:"cover_url"`
	UploadTime  string   `json:"upload_time"`
	UpdateTime  string   `json:"update_time"`
	Author      *UserDto `json:"author"`
}

func (VideoDto) FromPo(video *po.Video) *VideoDto {
	if !strings.HasPrefix(video.CoverUrl, "http") {
		if video.CoverUrl == "" {
			video.CoverUrl = "http://localhost:3344/raw/image/default/cover.jpg"
		} else {
			video.CoverUrl = fmt.Sprintf("http://localhost:3344/raw/image/%d/%s", video.AuthorUid, video.CoverUrl)
		}
	}
	return &VideoDto{
		Vid:         video.Vid,
		Title:       video.Title,
		Description: video.Description,
		VideoUrl:    video.VideoUrl,
		CoverUrl:    video.CoverUrl,
		UploadTime:  video.UploadTime.String(),
		UpdateTime:  common.JsonDateTime(video.UpdatedAt).String(),
		Author:      UserDto{}.FromPo(video.Author, enum.DtoOptionNone),
	}
}

func (VideoDto) FromPos(videos []*po.Video) []*VideoDto {
	dtos := make([]*VideoDto, len(videos))
	for idx, video := range videos {
		dtos[idx] = VideoDto{}.FromPo(video)
	}
	return dtos
}
