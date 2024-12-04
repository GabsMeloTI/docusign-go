package main

import (
	"context"
	"docusign/cmd"
	"docusign/config"
	"docusign/pkg"
	"github.com/labstack/gommon/log"
	"net/http"
	"os/signal"
	"syscall"
)

func main() {
	http.HandleFunc("/webhook", pkg.WebhookHandler)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	defer stop()

	loadingEnv := config.NewConfig()
	container := config.NewContainerDI(loadingEnv)

	cmd.StartAPI(ctx, container)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
