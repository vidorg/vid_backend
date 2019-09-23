package database

import (
	"time"

	"vid/config"
	. "vid/exceptions"
	. "vid/models"
)

type playlistDao struct{}

var PlaylistDao = new(playlistDao)

const (
	col_playlist_gid         = "gid"
	col_playlist_groupname   = "groupname"
	col_playlist_description = "description"
	col_playlist_create_time = "create_time"
	col_playlist_author_uid  = "author_uid"
)

const (
	col_videoinlist_gid = "gid"
	col_videoinlist_vid = "vid"
)

// db 查询所有视频列表和作者
//
// @return `[]Playlist`
func (p *playlistDao) QueryAllPlaylists() (playlist []Playlist) {
	DB.Find(&playlist)
	for k, v := range playlist {
		user, ok := UserDao.QueryUserByUid(v.AuthorUid)
		if ok {
			playlist[k].Author = user
		}
	}
	return playlist
}

// db 查询用户视频列表和作者
//
// @return `[]Playlist` `err`
//
// @error `UserNotExistException`
func (p *playlistDao) QueryPlaylistsByUid(uid int) ([]Playlist, error) {
	var playlists []Playlist
	if _, ok := UserDao.QueryUserByUid(uid); !ok {
		return nil, UserNotExistException
	}
	DB.Where(col_playlist_author_uid+" = ?", uid).Find(&playlists)
	user, ok := UserDao.QueryUserByUid(uid)
	if ok {
		for k, _ := range playlists {
			playlists[k].Author = user
		}
	}
	return playlists, nil
}

// db 查询用户视频列表数
//
// @return `playlist_cnt` `err`
//
// @error `UserNotExistException`
func (p *playlistDao) QueryUserPlaylistCnt(uid int) (int, error) {
	if _, ok := UserDao.QueryUserByUid(uid); !ok {
		return 0, UserNotExistException
	}
	var playlists []Playlist
	DB.Where(col_playlist_author_uid+" = ?", uid).Find(&playlists)
	return len(playlists), nil
}

// db 查询 gid 视频列表和作者
//
// @return `*Playlist` `isExist`
func (p *playlistDao) QueryPlaylistByGid(gid int) (*Playlist, bool) {
	var playlist Playlist
	nf := DB.Where(col_playlist_gid+" = ?", gid).Find(&playlist).RecordNotFound()
	if nf {
		return nil, false
	} else {
		user, ok := UserDao.QueryUserByUid(playlist.AuthorUid)
		if ok {
			playlist.Author = user
		}
		return &playlist, true
	}
}

// db 查询 uid + name 视频列表和作者
//
// @return `*Playlist` `isExist`
func (p *playlistDao) QueryPlaylistByUidName(uid int, gname string) (*Playlist, bool) {
	var playlist Playlist
	nf := DB.Where(col_playlist_author_uid+" = ?", uid).
		Where(col_playlist_groupname+" = ?", gname).
		Find(&playlist).RecordNotFound()

	if nf {
		return nil, false
	} else {
		user, ok := UserDao.QueryUserByUid(playlist.AuthorUid)
		if ok {
			playlist.Author = user
		}
		return &playlist, true
	}
}

// db 创建新列表
//
// @return `*Playlist` `err`
//
// @error `PlaylistNameUsedException` `CreatePlaylistException`
func (p *playlistDao) InsertPlaylist(playlist *Playlist) (*Playlist, error) {
	if _, ok := p.QueryPlaylistByUidName(playlist.AuthorUid, playlist.Groupname); ok {
		return nil, PlaylistNameUsedException
	}
	playlist.CreateTime = time.Now()
	DB.Create(playlist)
	query, ok := p.QueryPlaylistByGid(playlist.Gid)
	if !ok {
		return nil, CreatePlaylistException
	} else {
		return query, nil
	}
}

// db 更新旧视频列表
//
// @return `*Playlist` `err`
//
// @error `PlaylistNotExistException` `PlaylistNameUsedException` `NotUpdatePlaylistException`

func (p *playlistDao) UpdatePlaylist(playlist *Playlist) (*Playlist, error) {
	old, ok := p.QueryPlaylistByGid(playlist.Gid)
	if !ok {
		return nil, PlaylistNotExistException
	}

	// 更新空字段
	if playlist.Groupname == "" {
		playlist.Groupname = old.Groupname
	}
	if playlist.Description == config.AppCfg.MagicToken {
		playlist.Description = old.Description
	}

	// 检查同名
	if _, ok := p.QueryPlaylistByUidName(playlist.AuthorUid, playlist.Groupname); ok && playlist.Groupname != old.Groupname {
		return nil, PlaylistNameUsedException
	}

	DB.Model(playlist).Updates(map[string]interface{}{
		col_playlist_groupname:   playlist.Groupname,
		col_playlist_description: playlist.Description,
	})
	after, _ := p.QueryPlaylistByGid(playlist.Gid)
	if old.Equals(after) {
		return after, NotUpdatePlaylistException
	} else {
		return after, nil
	}
}

// db 删除视频列表
//
// @return `*video` `err`
//
// @error `PlaylistNotExistException` `DeletePlaylistException`
func (p *playlistDao) DeletePlaylist(gid int) (*Playlist, error) {
	query, ok := p.QueryPlaylistByGid(gid)
	if !ok {
		return nil, PlaylistNotExistException
	}
	if DB.Delete(query).RowsAffected != 1 {
		return nil, DeletePlaylistException
	} else {
		return query, nil
	}
}

///////////////////////////////////////////////////////////////////////////////////////////
