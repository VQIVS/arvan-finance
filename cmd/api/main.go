package main

import (
	"finance/config"
	"finance/internal/api/handlers/http"
	"finance/internal/app"
	"flag"
	"log"
	"os"
)

var configPath = flag.String("config", "config.yaml", "service configuration file")

func main() {
	flag.Parse()

	if v := os.Getenv("CONFIG_PATH"); len(v) > 0 {
		*configPath = v
	}
	c := config.MustReadConfig(*configPath)

	appContainer := app.NewMustApp(c)

	log.Fatal(http.Run(appContainer, c.Server))
}
