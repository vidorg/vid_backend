package exception

import (
	"errors"
)

// Authorization
var (
	AuthorizationException = errors.New("Authorization failed")
	TokenExpiredException  = errors.New("Token has expired")
	TokenInvalidException  = errors.New("Token invalid")

	PasswordException        = errors.New("User password error")
	NeedAdminException       = errors.New("Action need admin authority")
	NoAuthorizationException = errors.New("Don't have authorization to action")
)

// Db
var (
	// exist
	UserExistException        = errors.New("User already existed") // R
	UserNotExistException     = errors.New("User not found")       // R
	VideoNotExistException    = errors.New("Video not found")      // R
	PlaylistNotExistException = errors.New("Playlist not found")   // R

	// user crud failed
	InsertUserException    = errors.New("User insert failed")           // C
	NotUpdateUserException = errors.New("User information not updated") // U
	DeleteUserException    = errors.New("User delete failed")           // D
	ModifyPassException    = errors.New("User password modify failed")  // U

	// video crud failed
	CreateVideoException    = errors.New("Video insert failed")           // C
	NotUpdateVideoException = errors.New("Video information not updated") // U
	DeleteVideoException    = errors.New("Video delete failed")           // D

	// playlist crud failed
	CreatePlaylistException    = errors.New("Playlist insert failed")           // C
	NotUpdatePlaylistException = errors.New("Playlist information not updated") // U
	DeletePlaylistException    = errors.New("Playlist delete failed")           // D
	DeleteVideoInListException = errors.New("Video in playlist delete failed")  // D

	// user other exception
	UserNameUsedException     = errors.New("Username has been used")
	UserInfoException         = errors.New("User information invalid")
	SubscribeOneSelfException = errors.New("Cound not subscribe to oneself")

	// video other exception
	VideoUrlUsedException = errors.New("Video resource url has been used")

	// playlist other exception
	PlaylistNameUsedException = errors.New("Playlist name duplicated")

	// raw exception
	ImageUploadException  = errors.New("Image upload failed")
	VideoUploadException  = errors.New("Video upload failed")
	FileExtException      = errors.New("Extension not supported")
	FileNotExistException = errors.New("File not exist")
)

// Ctrl
var (
	// request error
	RequestBodyError = errors.New("Request body error")
	QueryParamError  = errors.New("Query param '%s' not found or error")
	RouteParamError  = errors.New("Route param '%s' not found or error")

	// format error
	LoginFormatError    = errors.New("Login username or password format error")
	RegisterFormatError = errors.New("Register username or password format error")
)
