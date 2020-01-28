package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/vidorg/vid_backend/src/config"
	"github.com/vidorg/vid_backend/src/server"
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

// @Title                       vid backend
// @Version                     1.1
// @Description                 Backend of repo https://github.com/vidorg/vid_vue
// @TermsOfService              https://github.com/vidorg
// @Host                        localhost:3344
// @BasePath                    /
// @License.Name                MIT
// @License.Url                 https://github.com/vidorg/vid_backend/blob/master/LICENSE

// @DemoModel                   ./docs/demo.json
// @Template Auth.Param         Authorization header string true "用户登录令牌"
// @Template Auth.ErrorCode     401 unauthorized user
// @Template Auth.ErrorCode     401 token has expired
// @Template Auth.ErrorCode     401 authorized user not found
// @Template Admin.ErrorCode    401 need admin authority

func run() {
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalln("Failed to load yaml config file:", err)
	}

	s := server.InitServer(cfg)

	fmt.Println()
	log.Printf("Server init on port :%d\n\n", cfg.HTTPConfig.Port)
	if err := s.ListenAndServe(); err != nil {
		log.Fatalln("Failed to listen and serve:", err)
	}
}
