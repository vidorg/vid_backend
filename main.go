package main

import (
	"flag"
	"fmt"
	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/router"
	"log"
)

var (
	help       bool
	configPath string
)

func init() {
	flag.BoolVar(&help, "h", false, "show help")
	flag.StringVar(&configPath, "config", "./src/config/Config.yaml", "change the config path")
}

func main() {
	flag.Parse()
	if help {
		flag.Usage()
	} else {
		run()
	}
}

// @title 					vid backend
// @version 				1.1
// @description 			Backend of repo https://github.com/vidorg/vid_vue
// @termsOfService 			https://github.com/vidorg
// @host 					localhost:3344
// @basePath 				/
// @license.name 			MIT
// @license.url 			https://github.com/vidorg/vid_backend/blob/master/LICENSE
// @swagger 				2.0

// @response.demopath		./src/model/dto/demo.json
// @authorization.param 	Authorization header string true "用户登录令牌"
// @authorization.error		401 authorization failed
// @authorization.error		401 token has expired

func run() {
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalln("Failed to load yaml config file:", err)
	}

	server := router.InitServer(cfg)

	fmt.Println()
	log.Printf("Server init on port :%d\n", cfg.HTTPConfig.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalln("Failed to listen and serve:", err)
	}
}
