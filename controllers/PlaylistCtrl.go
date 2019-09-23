package controllers

import (
	"fmt"
	"net/http"

	. "vid/database"
	. "vid/exceptions"
	. "vid/models"
	. "vid/models/resp"
	. "vid/utils"

	"github.com/gin-gonic/gin"
)

type playlistCtrl struct{}

var PlaylistCtrl = new(playlistCtrl)

// GET /playlist/all (Auth) (Admin)
func (p *playlistCtrl) GetAllPlaylists(c *gin.Context) {
	authusr, _ := c.Get("user")
	if authusr.(User).Authority != AuthAdmin {
		c.JSON(http.StatusUnauthorized, Message{
			Message: NeedAdminException.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, PlaylistDao.QueryAllPlaylists())
}

// GET /playlist/uid/:uid (Non-Auth)
func (p *playlistCtrl) GetPlaylistsByUid(c *gin.Context) {
	uid, ok := ReqUtil.GetIntParam(c.Params, "uid")
	if !ok {
		c.JSON(http.StatusBadRequest, Message{
			Message: fmt.Sprintf(RouteParamError.Error(), "uid"),
		})
		return
	}
	query, err := PlaylistDao.QueryPlaylistsByUid(uid)
	if err == nil {
		c.JSON(http.StatusOK, query)
	} else {
		c.JSON(http.StatusNotFound, Message{
			Message: err.Error(),
		})
	}
}

// GET /playlist/gid/:gid (Non-Auth)
func (p *playlistCtrl) GetPlaylistByGid(c *gin.Context) {
	vid, ok := ReqUtil.GetIntParam(c.Params, "gid")
	if !ok {
		c.JSON(http.StatusBadRequest, Message{
			Message: fmt.Sprintf(RouteParamError.Error(), "gid"),
		})
		return
	}
	query, ok := PlaylistDao.QueryPlaylistByGid(vid)
	if ok {
		c.JSON(http.StatusOK, query)
	} else {
		c.JSON(http.StatusNotFound, Message{
			Message: PlaylistNotExistException.Error(),
		})
	}
}

// POST /playlist/new (Auth)
func (p *playlistCtrl) CreateNewPlaylist(c *gin.Context) {
	body := ReqUtil.GetBody(c.Request.Body)
	var playlist Playlist
	if !playlist.Unmarshal(body, true) {
		c.JSON(http.StatusBadRequest, Message{
			Message: RequestBodyError.Error(),
		})
		return
	}

	authusr, _ := c.Get("user")
	uid := authusr.(User).Uid
	playlist.AuthorUid = uid

	query, err := PlaylistDao.InsertPlaylist(&playlist)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Message{
			Message: err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, query)
	}
}

// POST /playlist/update (Auth)
func (p *playlistCtrl) UpdatePlaylistInfo(c *gin.Context) {
	body := ReqUtil.GetBody(c.Request.Body)
	var playlist Playlist
	if !playlist.Unmarshal(body, false) {
		c.JSON(http.StatusBadRequest, Message{
			Message: RequestBodyError.Error(),
		})
		return
	}

	authusr, _ := c.Get("user")
	uid := authusr.(User).Uid
	playlist.AuthorUid = uid

	query, err := PlaylistDao.UpdatePlaylist(&playlist)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Message{
			Message: err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, query)
	}
}

// DELETE /playlist/delete?gid (Auth)
func (p *playlistCtrl) DeletePlaylist(c *gin.Context) {
	gid, ok := ReqUtil.GetIntQuery(c, "gid")
	if !ok {
		c.JSON(http.StatusBadRequest, Message{
			Message: fmt.Sprintf(RouteParamError.Error(), "gid"),
		})
		return
	}

	query, err := PlaylistDao.DeletePlaylist(gid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Message{
			Message: err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, query)
	}
}
