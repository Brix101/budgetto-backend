package main

import (
	"context"
	"os"
	"os/signal"

	budgettocmd "github.com/Brix101/budgetto-backend/cmd"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	ret := budgettocmd.Execute(ctx)
	os.Exit(ret)
}
