package main

import (
	"context"
	"github.com/vasch3nko/songlibrary/internal/app"
	"log"
	"os"
	"os/signal"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// Context initialization
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Starting app
	if err := app.Run(ctx); err != nil {
		return err
	}

	return nil
}
