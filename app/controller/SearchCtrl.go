package controller
//
// import (
// 	"fmt"
// 	"net/http"
// 	"strconv"
// 	"strings"
// 	"vid/app/controller/exception"
// 	"vid/app/database/dao"
// 	po2 "vid/app/model/po"
// 	"vid/app/model/resp"
// 	"vid/app/util"
//
// 	"github.com/gin-gonic/gin"
// )
//
// type searchCtrl struct{}
//
// var SearchCtrl = new(searchCtrl)
//
// // GET /search/user?keyword (Non-Auth)
// func (s *searchCtrl) SearchUser(c *gin.Context) {
// 	keyWord, ok := util.ReqUtil.GetStrQuery(c, "keyword")
// 	if !ok {
// 		c.JSON(http.StatusOK, resp.Message{
// 			Message: fmt.Sprintf(exception.QueryParamError.Error(), "keyword"),
// 		})
// 		return
// 	}
// 	ret := _searchUser(keyWord)
// 	if ret == nil {
// 		c.JSON(http.StatusOK, gin.H{})
// 	} else {
// 		c.JSON(http.StatusOK, ret)
// 	}
// }
//
// // GET /search/video?keyword (Non-Auth)
// func (s *searchCtrl) SearchVideo(c *gin.Context) {
// 	keyWord, ok := util.ReqUtil.GetStrQuery(c, "keyword")
// 	if !ok {
// 		c.JSON(http.StatusOK, resp.Message{
// 			Message: fmt.Sprintf(exception.QueryParamError.Error(), "keyword"),
// 		})
// 		return
// 	}
// 	ret := _searchVideo(keyWord)
// 	if ret == nil {
// 		c.JSON(http.StatusOK, gin.H{})
// 	} else {
// 		c.JSON(http.StatusOK, ret)
// 	}
// }
//
// // GET /search/playlist?keyword (Non-Auth)
// func (s *searchCtrl) SearchPlaylist(c *gin.Context) {
// 	keyWord, ok := util.ReqUtil.GetStrQuery(c, "keyword")
// 	if !ok {
// 		c.JSON(http.StatusOK, resp.Message{
// 			Message: fmt.Sprintf(exception.QueryParamError.Error(), "keyword"),
// 		})
// 		return
// 	}
// 	ret := _searchPlaylist(keyWord)
// 	if ret == nil {
// 		c.JSON(http.StatusOK, gin.H{})
// 	} else {
// 		c.JSON(http.StatusOK, ret)
// 	}
// }
//
// /////////////////////////////////////////////////////////////////////////////
//
// // SearchUser
// func _searchUser(keyWord string) interface{} {
// 	if strings.HasPrefix(keyWord, "uid:") {
// 		// uid 搜索
// 		uid, err := strconv.Atoi(strings.TrimLeft(keyWord, "uid:"))
// 		if err == nil {
// 			user, ok := dao.UserDao.QueryUserByUid(uid)
// 			if ok && user != nil {
// 				return [...]po2.User{*user}
// 			} else {
// 				return nil
// 			}
// 		}
// 	}
// 	// 用户名 搜索
// 	nameSp := strings.Split(keyWord, " ")
// 	query, ret := make([]po2.User, 0), make([]po2.User, 0)
// 	for _, v := range nameSp {
// 		if v != "" {
// 			query = append(query, dao.UserDao.SearchByUserName(v)...)
// 		}
// 	}
// 	// 去重
// 	for _, v := range query {
// 		ok := true
// 		for _, v2 := range ret {
// 			if v.Uid == v2.Uid {
// 				ok = false
// 				break
// 			}
// 		}
// 		if ok {
// 			ret = append(ret, v)
// 		}
// 	}
// 	return ret
// }
//
// // SearchVideo
// func _searchVideo(keyWord string) interface{} {
// 	if strings.HasPrefix(keyWord, "vid:") {
// 		// vid 搜索
// 		vid, err := strconv.Atoi(strings.TrimLeft(keyWord, "vid:"))
// 		if err == nil {
// 			video, ok := dao.VideoDao.QueryVideoByVid(vid)
// 			if ok && video != nil {
// 				return [...]po2.Video{*video}
// 			} else {
// 				return nil
// 			}
// 		}
// 	}
// 	// 标题 搜索
// 	titleSp := strings.Split(keyWord, " ")
// 	query, ret := make([]po2.Video, 0), make([]po2.Video, 0)
// 	for _, v := range titleSp {
// 		if v != "" {
// 			query = append(query, dao.VideoDao.SearchByVideoTitle(v)...)
// 		}
// 	}
// 	// 去重
// 	for _, v := range query {
// 		ok := true
// 		for _, v2 := range ret {
// 			if v.Vid == v2.Vid {
// 				ok = false
// 				break
// 			}
// 		}
// 		if ok {
// 			ret = append(ret, v)
// 		}
// 	}
// 	return ret
// }
//
// // SearchPlaylist
// func _searchPlaylist(keyWord string) interface{} {
// 	if strings.HasPrefix(keyWord, "gid:") {
// 		// gid 搜索
// 		gid, err := strconv.Atoi(strings.TrimLeft(keyWord, "gid:"))
// 		if err == nil {
// 			playlist, ok := dao.PlaylistDao.QueryPlaylistByGid(gid)
// 			if ok && playlist != nil {
// 				return [...]po2.Playlist{*playlist}
// 			} else {
// 				return nil
// 			}
// 		}
// 	}
// 	// 标题 搜索
// 	titleSp := strings.Split(keyWord, " ")
// 	query, ret := make([]po2.Playlist, 0), make([]po2.Playlist, 0)
// 	for _, v := range titleSp {
// 		if v != "" {
// 			query = append(query, dao.PlaylistDao.SearchByPlaylistTitle(v)...)
// 		}
// 	}
// 	// 去重
// 	for _, v := range query {
// 		ok := true
// 		for _, v2 := range ret {
// 			if v.Gid == v2.Gid {
// 				ok = false
// 				break
// 			}
// 		}
// 		if ok {
// 			ret = append(ret, v)
// 		}
// 	}
// 	return ret
// }
