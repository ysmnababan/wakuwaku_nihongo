package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"wakuwaku_nihongo/config"
	"wakuwaku_nihongo/internals/factory"
	middleware "wakuwaku_nihongo/internals/middleware"
	"wakuwaku_nihongo/internals/pkg/database"
	httpserver "wakuwaku_nihongo/internals/server"
	"wakuwaku_nihongo/internals/utils/env"
)

func init() {
	selectedEnv := config.Env()
	env := env.NewEnv()
	env.Load(`.env`)
	log.Info().Msg("Choosen environment " + selectedEnv)
}

// @title wakuwaku_nihongo-Project
// @version 0.0.1
// @description This is a doc for wakuwaku_nihongo-Project

// @securityDefinitions.apikey Authorization
// @in header
// @name Authorization
func main() {
	cfg := config.Get()

	port := cfg.App.Port

	logLevel, err := zerolog.ParseLevel(cfg.App.LogLevel)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(logLevel)

	database.Init("std")

	f := factory.NewFactory()

	e := echo.New()
	e.HideBanner = true
	e.IPExtractor = echo.ExtractIPDirect()
	middleware.Init(e, f.Redis)
	httpserver.Init(e, f)

	if err := e.Start(fmt.Sprintf(":%d", port)); err != nil && err != http.ErrServerClosed {
		e.Logger.Fatal("shutting down the server")
	}
}
