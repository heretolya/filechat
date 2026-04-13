package main

import (
	"context"
	"os"
	"os/signal"

	"filechat/internal/cli"

	"github.com/charmbracelet/fang"
)

func init() {
	cli.InitChat()
	cli.InitDrop()
}

func main() {
	ctx := context.Background()
	sigs := []os.Signal{os.Interrupt, os.Kill}
	ctx, stop := signal.NotifyContext(ctx, sigs...)
	defer stop()

	cli.ChatCmd.AddCommand(cli.DropCmd)
	err := fang.Execute(ctx, cli.ChatCmd)
	if err != nil {
		os.Exit(1)
	}

	stop()
	<-ctx.Done()
}
