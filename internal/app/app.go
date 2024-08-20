package app

import (
	"log/slog"

	restapp "github.com/k3k13a1/vk_tarantool_test/internal/app/rest"
	"github.com/k3k13a1/vk_tarantool_test/internal/config"
	"github.com/k3k13a1/vk_tarantool_test/internal/services"
	tt "github.com/k3k13a1/vk_tarantool_test/internal/storage/tarantool"
)

type App struct {
	RESTSrv *restapp.App
}

func New(
	cfg config.Config,
) *App {
	ttStorage, err := tt.New(cfg.Tarantool.User, cfg.Tarantool.Pass, cfg.Tarantool.Host, cfg.Tarantool.Port)
	if err != nil {
		slog.Error("Can't create tarantool storage", slog.String("error", err.Error()))
		panic(err)
	}

	ttService := services.New(ttStorage)

	restApp := restapp.New(*ttService, cfg.Srv.Host, cfg.Srv.Port)

	return &App{
		RESTSrv: restApp,
	}
}
