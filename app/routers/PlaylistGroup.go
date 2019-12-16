package routers

import (
	. "vid/app/controllers"
	"vid/app/middleware"

	"github.com/gin-gonic/gin"
)

func SetupPlaylistGroup(router *gin.Engine) {

	jwt := middleware.JWTMiddleware()

	playlistGroup := router.Group("/playlist")
	{
		// Admin
		playlistGroup.Use(jwt).GET("/all", PlaylistCtrl.GetAllPlaylists)

		// Public
		playlistGroup.GET("/uid/:uid", PlaylistCtrl.GetPlaylistsByUid)
		playlistGroup.GET("/gid/:gid", PlaylistCtrl.GetPlaylistByGid)

		// Auth
		playlistGroup.Use(jwt).POST("/new", PlaylistCtrl.CreateNewPlaylist)
		playlistGroup.Use(jwt).PUT("/update", PlaylistCtrl.UpdatePlaylistInfo)
		playlistGroup.Use(jwt).DELETE("/delete", PlaylistCtrl.DeletePlaylist)

		playlistGroup.Use(jwt).POST("/add", PlaylistCtrl.AddVideosInList)
		playlistGroup.Use(jwt).DELETE("/remove", PlaylistCtrl.RemoveVideosInList)
	}
}
