package server

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xstring"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
)

func SetupLogger() {
	logFile, err := os.Create(fmt.Sprintf("./log/log-%s.log", xstring.CurrentTimeUuid(14)))
	if err != nil {
		log.Fatalln("Failed to create log file:", err)
	}

	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)
}