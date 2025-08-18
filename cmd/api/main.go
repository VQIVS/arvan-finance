package main

import (
	"billing-service/api/handler/http"
	"billing-service/app"
	"billing-service/config"
	"flag"
	"log"
	"os"
)

var configPath = flag.String("config", "config.json", "service configuration file")

func main() {
	flag.Parse()

	if v := os.Getenv("CONFIG_PATH"); len(v) > 0 {
		*configPath = v
	}

	c := config.MustReadConfig(*configPath)

	appContainer := app.NewMustApp(c)

	log.Fatal(http.Run(appContainer, c.Server))
}
