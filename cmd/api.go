package cmd

import (
	"context"
	"docusign/config"
	_ "docusign/docs/app"
	_middleware "docusign/infra/middleware"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
	"os/exec"
	"strings"

	"time"
)

func startNgrok() (string, error) {
	cmd := exec.Command("ngrok", "http", "8080")

	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("erro ao iniciar o ngrok: %v", err)
	}

	time.Sleep(2 * time.Second)

	out, err := exec.Command("curl", "-s", "http://127.0.0.1:4040/api/tunnels").Output()
	if err != nil {
		return "", fmt.Errorf("erro ao obter informações do ngrok: %v", err)
	}

	output := string(out)
	startIndex := strings.Index(output, "https://")
	endIndex := strings.Index(output[startIndex:], "\"")
	if startIndex == -1 || endIndex == -1 {
		return "", fmt.Errorf("não foi possível encontrar o endpoint público do ngrok")
	}
	return output[startIndex : startIndex+endIndex], nil
}

// @title Swagger First Access API
// @version 1.0
// @description Document API
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func StartAPI(ctx context.Context, container *config.ContainerDI) {

	e := echo.New()

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	go func() {
		for {
			select {
			case <-ctx.Done():
				if err := e.Shutdown(ctx); err != nil {
					panic(err)
				}
				return
			default:
				time.Sleep(1 * time.Second)
			}
		}
	}()

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	docuSign := e.Group("/docusign", _middleware.CheckAuthorization)
	docuSign.POST("/send", container.HandlerContract.CreateContractHandler)

	webhookHandler := container.WebhookHandler
	e.POST("/webhook", webhookHandler.HandleWebhook)

	ngrokURL, err := startNgrok()
	if err != nil {
		e.Logger.Fatalf("Erro ao iniciar o ngrok: %v", err)
	}
	e.Logger.Infof("ngrok iniciado no endpoint: %s", ngrokURL)

	go func() {
		client := &http.Client{Timeout: 5 * time.Second}
		appRunnerURL := "https://qmdt6r8upv.us-east-1.awsapprunner.com/webhook"
		req, _ := http.NewRequest("POST", appRunnerURL, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Webhook-URL", appRunnerURL)

		resp, err := client.Do(req)
		if err != nil {
			e.Logger.Errorf("Erro ao registrar o webhook no App Runner: %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			e.Logger.Errorf("Falha ao registrar webhook, status: %d", resp.StatusCode)
		} else {
			e.Logger.Info("Webhook registrado com sucesso no App Runner.")
		}
	}()

	e.Logger.Fatal(e.Start(container.Config.ServerPort))
}
