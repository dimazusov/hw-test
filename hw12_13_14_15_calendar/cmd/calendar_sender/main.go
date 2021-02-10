package main

import (
	"flag"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", `confisdgs/config.yaml`, "Path to configuration file")
	flag.Parse()
}

func main() {
	// sender

	// read from rabbit mq
}
