package server

import (
	"github.com/Aoi-hosizora/ahlib-gin-gorm/xdatetime"
	"github.com/Aoi-hosizora/ahlib-gin-gorm/xgin"
	"log"
	"time"
)

func BindValidation() {
	xgin.SetupRegexBinding()

	shanghaiLoc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Fatalln("Failed to load location:", err)
	}
	xgin.SetupDateTimeBinding("date", xdatetime.DateFormat, shanghaiLoc)
	xgin.SetupDateTimeBinding("time", xdatetime.TimeFormat, shanghaiLoc)
	xgin.SetupDateTimeBinding("datetime", xdatetime.DateTimeFormat, shanghaiLoc)

	xgin.SetupSpecificRegexpBinding("name", "^[a-zA-Z0-9\u4E00-\u9FBF\u3040-\u30FF\\-_]+$")              // alphabet number character kana - _
	xgin.SetupSpecificRegexpBinding("pwd", "^.+$")                                                       // all
	xgin.SetupSpecificRegexpBinding("phone", "^(13[0-9]|15[012356789]|17[678]|18[0-9]|14[57])[0-9]{8}$") // 11
}
