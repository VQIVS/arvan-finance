package main

import (
	"flag"
	"os"
)

var configPath = flag.String("config", "config.json", "service configuration file")

func main() {
	flag.Parse()
	if v := os.Getenv("CONFIG_PATH"); len(v) > 0 {
		*configPath = v
	}

	// c := config.MustReadConfig(*configPath)

	//appContainer

	// logger

}
