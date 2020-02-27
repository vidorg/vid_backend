package server

import (
	"github.com/Aoi-hosizora/ahlib-gin-gorm/xdatetime"
	"github.com/Aoi-hosizora/ahlib-gin-gorm/xgin"
)

func BindValidation() {
	xgin.SetupRegexBinding()

	xgin.SetupDateTimeBinding("date", xdatetime.ISO8601DateFormat)
	xgin.SetupDateTimeBinding("datetime", xdatetime.ISO8601DateTimeFormat)

	xgin.SetupSpecificRegexpBinding("name", "^[a-zA-Z0-9\u4E00-\u9FBF\u3040-\u30FF\\-_]+$")              // alphabet number character kana - _
	xgin.SetupSpecificRegexpBinding("pwd", "^.+$")                                                       // all
	xgin.SetupSpecificRegexpBinding("phone", "^(13[0-9]|15[012356789]|17[678]|18[0-9]|14[57])[0-9]{8}$") // 11
}
