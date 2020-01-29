package server

import (
	"github.com/Aoi-hosizora/ahlib-gin-gorm/xdatetime"
	"github.com/Aoi-hosizora/ahlib-gin-gorm/xgin"
	"time"
)

func SetupDefinedValidation() {
	xgin.SetupRegexBinding()

	shanghaiLoc, _ := time.LoadLocation("Asia/Shanghai")
	xgin.SetupDateTimeBinding("date", xdatetime.DateFormat, shanghaiLoc)
	xgin.SetupDateTimeBinding("time", xdatetime.TimeFormat, shanghaiLoc)
	xgin.SetupDateTimeBinding("datetime", xdatetime.DateTimeFormat, shanghaiLoc)

	xgin.SetupSpecificRegexpBinding("name", "^[^'`\"\\\\]+$")          // ' ` " \
	xgin.SetupSpecificRegexpBinding("pwd", "^[a-zA-Z0-9+\\-*/.=_~]+$") // + - * / . = _ ~

	xgin.SetupSpecificRegexpBinding("phone", "^(13[0-9]|15[012356789]|17[678]|18[0-9]|14[57])[0-9]{8}$") // 11
}
