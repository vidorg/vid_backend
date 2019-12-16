package controller
//
// import (
// 	"fmt"
// 	"net/http"
// 	"vid/app/controller/exception"
// 	"vid/app/database/dao"
// 	"vid/app/model/enum"
// 	po2 "vid/app/model/po"
// 	"vid/app/model/req"
// 	"vid/app/model/resp"
// 	"vid/app/util"
//
// 	"github.com/gin-gonic/gin"
// )
//
// type playlistCtrl struct{}
//
// var PlaylistCtrl = new(playlistCtrl)
//
// // GET /playlist/all (Auth) (Admin)
// func (p *playlistCtrl) GetAllPlaylists(c *gin.Context) {
// 	authusr, _ := c.Get("user")
// 	if authusr.(po2.User).Authority != enum.AuthAdmin {
// 		c.JSON(http.StatusUnauthorized, resp.Message{
// 			Message: exception.NeedAdminError.Error(),
// 		})
// 		return
// 	}
//
// 	c.JSON(http.StatusOK, dao.PlaylistDao.QueryAllPlaylists())
// }
//
// // GET /playlist/uid/:uid (Non-Auth)
// func (p *playlistCtrl) GetPlaylistsByUid(c *gin.Context) {
// 	uid, ok := util.ReqUtil.GetIntParam(c.Params, "uid")
// 	if !ok {
// 		c.JSON(http.StatusBadRequest, resp.Message{
// 			Message: fmt.Sprintf(exception.RouteParamError.Error(), "uid"),
// 		})
// 		return
// 	}
// 	query, err := dao.PlaylistDao.QueryPlaylistsByUid(uid)
// 	if err == nil {
// 		c.JSON(http.StatusOK, query)
// 	} else {
// 		c.JSON(http.StatusNotFound, resp.Message{
// 			Message: err.Error(),
// 		})
// 	}
// }
//
// // GET /playlist/gid/:gid (Non-Auth)
// func (p *playlistCtrl) GetPlaylistByGid(c *gin.Context) {
// 	vid, ok := util.ReqUtil.GetIntParam(c.Params, "gid")
// 	if !ok {
// 		c.JSON(http.StatusBadRequest, resp.Message{
// 			Message: fmt.Sprintf(exception.RouteParamError.Error(), "gid"),
// 		})
// 		return
// 	}
// 	query, ok := dao.PlaylistDao.QueryPlaylistByGid(vid)
// 	if ok {
// 		c.JSON(http.StatusOK, query)
// 	} else {
// 		c.JSON(http.StatusNotFound, resp.Message{
// 			Message: exception.PlaylistNotFoundError.Error(),
// 		})
// 	}
// }
//
// // POST /playlist/new (Auth)
// func (p *playlistCtrl) CreateNewPlaylist(c *gin.Context) {
// 	body := util.ReqUtil.GetBody(c.Request.Body)
// 	var playlist po2.Playlist
// 	if !playlist.Unmarshal(body, true) {
// 		c.JSON(http.StatusBadRequest, resp.Message{
// 			Message: exception.JsonParamError.Error(),
// 		})
// 		return
// 	}
//
// 	authusr, _ := c.Get("user")
// 	uid := authusr.(po2.User).Uid
// 	playlist.AuthorUid = uid
//
// 	query, err := dao.PlaylistDao.InsertPlaylist(&playlist)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, resp.Message{
// 			Message: err.Error(),
// 		})
// 	} else {
// 		c.JSON(http.StatusOK, query)
// 	}
// }
//
// // PUT /playlist/update (Auth)
// func (p *playlistCtrl) UpdatePlaylistInfo(c *gin.Context) {
// 	body := util.ReqUtil.GetBody(c.Request.Body)
// 	var playlist po2.Playlist
// 	if !playlist.Unmarshal(body, false) {
// 		c.JSON(http.StatusBadRequest, resp.Message{
// 			Message: exception.JsonParamError.Error(),
// 		})
// 		return
// 	}
//
// 	authusr, _ := c.Get("user")
// 	query, err := dao.PlaylistDao.UpdatePlaylist(&playlist, authusr.(po2.User).Uid)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, resp.Message{
// 			Message: err.Error(),
// 		})
// 	} else {
// 		c.JSON(http.StatusOK, query)
// 	}
// }
//
// // DELETE /playlist/delete?gid (Auth)
// func (p *playlistCtrl) DeletePlaylist(c *gin.Context) {
// 	gid, ok := util.ReqUtil.GetIntQuery(c, "gid")
// 	if !ok {
// 		c.JSON(http.StatusBadRequest, resp.Message{
// 			Message: fmt.Sprintf(exception.RouteParamError.Error(), "gid"),
// 		})
// 		return
// 	}
//
// 	authusr, _ := c.Get("user")
// 	query, err := dao.PlaylistDao.DeletePlaylist(gid, authusr.(po2.User).Uid)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, resp.Message{
// 			Message: err.Error(),
// 		})
// 	} else {
// 		c.JSON(http.StatusOK, query)
// 	}
// }
//
// // POST /playlist/add (Auth)
// func (p *playlistCtrl) AddVideosInList(c *gin.Context) {
// 	body := util.ReqUtil.GetBody(c.Request.Body)
// 	var vreq req.VideoinlistReq
// 	if !vreq.Unmarshal(body) {
// 		c.JSON(http.StatusBadRequest, resp.Message{
// 			Message: exception.JsonParamError.Error(),
// 		})
// 		return
// 	}
//
// 	authusr, _ := c.Get("user")
// 	query, err := dao.PlaylistDao.InsertVideosInList(vreq.Gid, vreq.Vids, authusr.(po2.User).Uid)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, resp.Message{
// 			Message: err.Error(),
// 		})
// 	} else {
// 		c.JSON(http.StatusOK, query)
// 	}
// }
//
// // DELETE /playlist/remove (Auth)
// func (p *playlistCtrl) RemoveVideosInList(c *gin.Context) {
// 	body := util.ReqUtil.GetBody(c.Request.Body)
// 	var vreq req.VideoinlistReq
// 	if !vreq.Unmarshal(body) {
// 		c.JSON(http.StatusBadRequest, resp.Message{
// 			Message: exception.JsonParamError.Error(),
// 		})
// 		return
// 	}
//
// 	authusr, _ := c.Get("user")
// 	query, err := dao.PlaylistDao.DeleteVideosInList(vreq.Gid, vreq.Vids, authusr.(po2.User).Uid)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, resp.Message{
// 			Message: err.Error(),
// 		})
// 	} else {
// 		c.JSON(http.StatusOK, query)
// 	}
// }
