package main

import (
	"billing-service/api/handler/consumer"
	"billing-service/app"
	"billing-service/config"
	"context"
	"log"
	"os"
	"os/signal"
)

func main() {
	cfg := config.MustReadConfig("config.json")

	a, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	h := consumer.New(a)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	if err := h.Start(ctx); err != nil {
		log.Printf("consumer stopped with error: %v", err)
		os.Exit(1)
	}
}
