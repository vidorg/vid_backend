package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	. "vid/database"
	. "vid/exceptions"
	. "vid/models"
	. "vid/models/resp"
	. "vid/utils"

	"github.com/gin-gonic/gin"
)

type searchCtrl struct{}

var SearchCtrl = new(searchCtrl)

// GET /search/user?keyword (Non-Auth)
func (s *searchCtrl) SearchUser(c *gin.Context) {
	keyWord, ok := ReqUtil.GetStrQuery(c, "keyword")
	if !ok {
		c.JSON(http.StatusOK, Message{
			Message: fmt.Sprintf(QueryParamError.Error(), "keyword"),
		})
		return
	}
	ret := _searchUser(keyWord)
	if ret == nil {
		c.JSON(http.StatusOK, gin.H{})
	} else {
		c.JSON(http.StatusOK, ret)
	}
}

// GET /search/video?keyword (Non-Auth)
func (s *searchCtrl) SearchVideo(c *gin.Context) {
	keyWord, ok := ReqUtil.GetStrQuery(c, "keyword")
	if !ok {
		c.JSON(http.StatusOK, Message{
			Message: fmt.Sprintf(QueryParamError.Error(), "keyword"),
		})
		return
	}
	ret := _searchVideo(keyWord)
	if ret == nil {
		c.JSON(http.StatusOK, gin.H{})
	} else {
		c.JSON(http.StatusOK, ret)
	}
}

/////////////////////////////////////////////////////////////////////////////

// SearchUser
func _searchUser(keyWord string) interface{} {
	if strings.HasPrefix(keyWord, "uid:") {
		// uid 搜索
		uid, err := strconv.Atoi(strings.TrimLeft(keyWord, "uid:"))
		if err == nil {
			user, ok := UserDao.QueryUserByUid(uid)
			if ok && user != nil {
				return [...]User{*user}
			} else {
				return nil
			}
		}
	}
	// 用户名 搜索
	nameSp := strings.Split(keyWord, " ")
	query, ret := make([]User, 0), make([]User, 0)
	for _, v := range nameSp {
		if v != "" {
			query = append(query, UserDao.SearchByUserName(v)...)
		}
	}
	// 去重
	for _, v := range query {
		ok := true
		for _, v2 := range ret {
			if v.Uid == v2.Uid {
				ok = false
				break
			}
		}
		if ok {
			ret = append(ret, v)
		}
	}
	return ret
}

// SearchVideo
func _searchVideo(keyWord string) interface{} {
	if strings.HasPrefix(keyWord, "vid:") {
		// vid 搜索
		vid, err := strconv.Atoi(strings.TrimLeft(keyWord, "vid:"))
		if err == nil {
			video, ok := VideoDao.QueryVideoByVid(vid)
			if ok && video != nil {
				return [...]Video{*video}
			} else {
				return nil
			}
		}
	}
	// 标题 搜索
	titleSp := strings.Split(keyWord, " ")
	query, ret := make([]Video, 0), make([]Video, 0)
	for _, v := range titleSp {
		if v != "" {
			query = append(query, VideoDao.SearchByVideoTitle(v)...)
		}
	}
	// 去重
	for _, v := range query {
		ok := true
		for _, v2 := range ret {
			if v.Vid == v2.Vid {
				ok = false
				break
			}
		}
		if ok {
			ret = append(ret, v)
		}
	}
	return ret
}
