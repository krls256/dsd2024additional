package main

import (
	"context"
	chatDI "github.com/krls256/dsd2024additional/internal/chat/di"
	"github.com/krls256/dsd2024additional/pkg/constants"
	pkgDI "github.com/krls256/dsd2024additional/pkg/di"
	"github.com/krls256/dsd2024additional/pkg/pgsql"
	"github.com/krls256/dsd2024additional/pkg/transport/http"
	"github.com/krls256/dsd2024additional/utils"

	"go.uber.org/zap"
	"time"
)

func main() {
	now := time.Now()

	defs, pgSQLMigrations := chatDI.Defs(), []pgsql.SmartEmbed{}

	ctn, err := pkgDI.Build("./config.chat.yml", pgSQLMigrations, defs...)
	if err != nil {
		panic(err)
	}

	logger, ok := ctn.Get(constants.LoggerName).(*zap.Logger)
	if !ok {
		panic("no logger available")
	}

	logger.Info("Starting application...")

	server, ok := ctn.Get(constants.HTTPServerName).(*http.Server)
	if !ok {
		panic("no server available")
	}

	server.AsyncRun()

	zap.S().Infof("Up and running (%s)", time.Since(now))
	zap.S().Infof("Got %s signal. Shutting down...", <-utils.WaitTermSignal())

	if err := server.Shutdown(context.Background()); err != nil {
		zap.S().Errorf("Error stopping server: %s", err)
	}

	zap.S().Info("Service stopped.")
}
