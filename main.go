package main

import (
	"context"
	"docusign/cmd"
	"docusign/config"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	defer stop()

	loadingEnv := config.NewConfig()
	container := config.NewContainerDI(loadingEnv)

	cmd.StartAPI(ctx, container)
}
