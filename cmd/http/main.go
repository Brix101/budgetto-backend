package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/Brix101/budgetto-backend/internal/cmd"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	ret := cmd.Execute(ctx)
	os.Exit(ret)
}
