package req

import (
	"net/http"
	"strings"

	"vid/config"
)

type RegLogReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Form-Data
func (r *RegLogReq) FromFormData(request *http.Request) bool {
	username := request.FormValue("username")
	password := request.FormValue("password")
	if username == "" || strings.Index(username, " ") != -1 ||
		password == "" || strings.Index(password, " ") != -1 {
		return false
	} else {
		r.Username = username
		r.Password = password
		return true
	}
}

func (r *RegLogReq) CheckFormat() bool {
	cfg := config.AppCfg
	return len(r.Username) >= cfg.FormatConfig.MinLen_Username &&
		len(r.Username) <= cfg.FormatConfig.MaxLen_Username &&
		len(r.Password) >= cfg.FormatConfig.MinLen_Password &&
		len(r.Password) <= cfg.FormatConfig.MaxLen_Password
}

//////////////////////////

type PassReq struct {
	Password string `json:"password"`
}

// Form-Data
func (p *PassReq) FromFormData(request *http.Request) bool {
	password := request.FormValue("password")
	if password == "" || strings.Index(password, " ") != -1 {
		return false
	} else {
		p.Password = password
		return true
	}
}

func (p *PassReq) CheckFormat() bool {
	cfg := config.AppCfg
	return len(p.Password) >= cfg.FormatConfig.MinLen_Password &&
		len(p.Password) <= cfg.FormatConfig.MaxLen_Password
}
