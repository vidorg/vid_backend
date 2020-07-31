package main

import (
	"flag"
	"github.com/Aoi-hosizora/goapidoc"
	"github.com/vidorg/vid_backend/src/provide"
	"github.com/vidorg/vid_backend/src/server"
	"log"
)

var (
	fConfig = flag.String("config", "./config.yaml", "change config path")
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

func run() {
	_, err := goapidoc.GenerateJsonWithSwagger2("./docs/doc.json")
	if err != nil {
		log.Fatalln("Failed to generate swagger:", err)
	}

	provide.Provide(*fConfig)
	s := server.NewServer()

	err = s.Serve()
	if err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}
