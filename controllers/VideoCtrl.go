package controllers

import (
	"net/http"
	// "time"

	"github.com/gin-gonic/gin"

	. "vid/models"
	. "vid/database"
)

type VideoCtrl struct{}

func (v *VideoCtrl) GetVideo(c *gin.Context) {
	// user := User{
	// 	Username: "12345678",
	// 	RegisterTime: time.Now(),
	// }
	// DB.Create(&user)
	// video := Video{
	// 	UploadTime: time.Now(),
	// 	Title: "123",
	// 	AuthorUid: user.Uid,
	// }
	// DB.Create(&video)
	// video = Video{
	// 	UploadTime: time.Now(),
	// 	Title: "456",
	// 	AuthorUid: user.Uid,
	// }
	// DB.Create(&video)

	/////////////////////////////////////

	video := Video{
		Vid : 1,
	}
	DB.Find(&video)
	user , _:=userDao.QueryUserByUid(video.AuthorUid)
	video.Author = user

	c.JSON(http.StatusOK, video)
}

func (v *VideoCtrl) GetUserVideos(c *gin.Context) {
	quser, _ := userDao.QueryUserByUid(1)
	videos := make([]*Video, 1, 1)
	DB.Model(quser).Related(&videos, "Videos")
	c.JSON(http.StatusOK, videos)
}
