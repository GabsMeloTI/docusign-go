package cmd

import (
	"context"
	"docusign/config"
	_ "docusign/docs/app"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	echoSwagger "github.com/swaggo/echo-swagger"
	"time"
)

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

	docuSign := e.Group("/docusign")
	docuSign.POST("/send", container.HandlerContract.CreateContractHandler)

	e.Logger.Fatal(e.Start(container.Config.ServerPort))
}
