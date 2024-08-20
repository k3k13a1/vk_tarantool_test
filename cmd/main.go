package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/k3k13a1/vk_tarantool_test/internal/app"
	"github.com/k3k13a1/vk_tarantool_test/internal/config"
	"github.com/k3k13a1/vk_tarantool_test/internal/logger"
)

func main() {
	cfg := config.SetupConfig()

	logger.SetupLogger()

	slog.Info("starting app", slog.Int("port", cfg.Srv.Port))

	application := app.New(*cfg)

	go application.RESTSrv.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.RESTSrv.Stop()
	slog.Info("app stopped", slog.Int("port", cfg.Srv.Port))
}
