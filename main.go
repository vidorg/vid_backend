package main

import (
	"flag"
	"github.com/vidorg/vid_backend/src/provide"
	"github.com/vidorg/vid_backend/src/server"
	"log"
)

var (
	fConfig = flag.String("config", "./config/config.yaml", "change config path")
	fHelp   = flag.Bool("h", false, "show help")
)

func main() {
	flag.Parse()
	if *fHelp {
		flag.Usage()
	} else {
		run()
	}
}

// @Title            vid backend
// @Version          1.1
// @Description      Backend of repo https://github.com/vidorg/vid_vue
// @TermsOfService   https://github.com/vidorg/vid_backend/issues
// @Host             127.0.0.1:3344
// @BasePath         /
// @License.Name     MIT
// @License.Url      https://github.com/vidorg/vid_backend/blob/master/LICENSE
// @Contact.Name     vidorg
// @Contact.Url      https://github.com/vidorg

// @Tag              Ping           "Ping"
// @Tag              Authorization  "Auth-Controller"
// @Tag              User           "User-Controller"
// @Tag              Subscribe      "Sub-Controller"
// @Tag              Video          "Video-Controller"
// @Tag              Raw            "Raw-Controller"
// @Tag              Policy         "Policy-Controller"
// @Tag              Search         "Search-Controller"
// @Tag              Administration "*-Controller"
// @GlobalSecurity   Jwt Authorization header

// @Template Page.Param    page  query integer false "当前页" (default:1)
// @Template Page.Param    limit query integer false "页大小" (default:10)
// @Template Order.Param   order query string  false "排序字符串"

func run() {
	provide.Provide(*fConfig)

	s := server.NewServer()
	err := s.Serve()
	if err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}
