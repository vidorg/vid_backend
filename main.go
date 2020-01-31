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

// @Title            vid backend
// @Version          1.1
// @Description      Backend of repo https://github.com/vidorg/vid_vue
// @TermsOfService   https://github.com/vidorg
// @Host             127.0.0.1:3344
// @BasePath         /
// @License.Name     MIT
// @License.Url      https://github.com/vidorg/vid_backend/blob/master/LICENSE

// @Tag              Ping           "Ping"
// @Tag              Authorization  "Auth-Controller"
// @Tag              User           "User-Controller"
// @Tag              Subscribe      "Sub-Controller"
// @Tag              Video          "Video-Controller"
// @Tag              Raw            "Raw-Controller"
// @Tag              Administration "*-Controller"
// @GlobalSecurity   Jwt Authorization header
// @DemoModel        ./docs/demo.json

// @Template Page.Param            page query integer false false "当前页" 1
// @Template Auth.ResponseDesc     401 "unauthorized user"
// @Template Auth.ResponseDesc     401 "token has expired"
// @Template Auth.ResponseDesc     401 "authorized user not found"
// @Template Admin.ResponseDesc    401 "need admin authority"
// @Template Param.ResponseDesc    400 "request param error"
// @Template Param.ResponseDesc    400 "request format error"
// @Template ParamA.ResponseDesc   400 "request param error"

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
