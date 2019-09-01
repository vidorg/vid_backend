package head

import (
	"strings"
	"vid/config"
)

type RegLogHead struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// @override
func (r *RegLogHead) CheckValid() bool {
	return r.Username != "" && r.Password != "" &&
		strings.Index(r.Username, " ") == -1
}

func (r *RegLogHead) CheckFormat() bool {
	cfg := config.AppCfg
	return len(r.Username) >= cfg.FormatConfig.MinLen_Username &&
		len(r.Username) <= cfg.FormatConfig.MaxLen_Username &&
		len(r.Password) >= cfg.FormatConfig.MinLen_Password &&
		len(r.Password) <= cfg.FormatConfig.MaxLen_Password
}
