package main

import "flag"

var (
	Config = flag.String("config", "./config.yaml", "配置文件路径")
	Help   = flag.Bool("h", false, "show help")

	survivalTimeout = int(3e9)
)
