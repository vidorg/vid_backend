package head

import (
	"encoding/json"
	"strings"
	"vid/config"
)

type RegLogReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *RegLogReq) Unmarshal(jsonBody string) bool {
	err := json.Unmarshal([]byte(jsonBody), r)
	if err != nil ||
		r.Username == "" || r.Password == "" ||
		strings.Index(r.Username, " ") != -1 || strings.Index(r.Password, " ") != -1 {
		return false
	}
	return true
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

func (p *PassReq) CheckFormat() bool {
	cfg := config.AppCfg
	return len(p.Password) >= cfg.FormatConfig.MinLen_Password &&
		len(p.Password) <= cfg.FormatConfig.MaxLen_Password
}

func (p *PassReq) Unmarshal(jsonBody string) bool {
	err := json.Unmarshal([]byte(jsonBody), p)
	if err != nil || p.Password == "" || strings.Index(p.Password, " ") != -1 {
		return false
	}
	return true
}
