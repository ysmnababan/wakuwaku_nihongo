package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	echoSwagger "github.com/swaggo/echo-swagger"
	"wakuwaku_nihongo/config"
	"wakuwaku_nihongo/docs"
	"wakuwaku_nihongo/internals/app/example_feat"
	"wakuwaku_nihongo/internals/factory"
)

func Init(e *echo.Echo, f *factory.Factory) {
	cfg := config.Get()

	// index
	e.GET("/", func(c echo.Context) error {
		message := fmt.Sprintf("Welcome to %s", cfg.App.Name)
		return c.String(http.StatusOK, message)
	})

	// doc
	if config.Get().EnableSwagger {
		docs.SwaggerInfo.Title = cfg.App.Name
		docs.SwaggerInfo.Host = cfg.App.URL
		docs.SwaggerInfo.Schemes = []string{cfg.App.Schema, "https"}
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	// routes v1
	api := e.Group("/api/v1")

	example_feat.NewHandler(f).Route(api.Group("/users"))
}
