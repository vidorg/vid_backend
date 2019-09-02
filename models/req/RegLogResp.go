package head

import (
	"strings"
	"vid/config"
)

type RegLogReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// @override
func (r *RegLogReq) CheckValid() bool {
	return r.Username != "" && r.Password != "" &&
		strings.Index(r.Username, " ") == -1
}

func (r *RegLogReq) CheckFormat() bool {
	cfg := config.AppCfg
	return len(r.Username) >= cfg.FormatConfig.MinLen_Username &&
		len(r.Username) <= cfg.FormatConfig.MaxLen_Username &&
		len(r.Password) >= cfg.FormatConfig.MinLen_Password &&
		len(r.Password) <= cfg.FormatConfig.MaxLen_Password
}

type PassReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// @override
func (p *PassReq) CheckValid() bool {
	return p.Password != ""
}

func (p *PassReq) CheckFormat() bool {
	cfg := config.AppCfg
	return len(p.Password) >= cfg.FormatConfig.MinLen_Password &&
		len(p.Password) <= cfg.FormatConfig.MaxLen_Password
}
