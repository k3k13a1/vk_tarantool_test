package restapp

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/k3k13a1/vk_tarantool_test/internal/handlers"
	"github.com/k3k13a1/vk_tarantool_test/internal/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type App struct {
	RESTServer *echo.Echo
	Host       string
	Port       int
}

func New(
	ttService services.Service,
	host string,
	port int,
) *App {
	e := echo.New()

	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.Secure())

	e.POST("/api/login", func(c echo.Context) error {
		return handlers.Login(c, &ttService)
	})

	e.POST("/api/write", func(c echo.Context) error {
		return handlers.Write(c, &ttService)
	})

	e.POST("/api/read", func(c echo.Context) error {
		return handlers.Read(c, &ttService)
	})

	return &App{
		RESTServer: e,
		Host:       host,
		Port:       port,
	}
}

func (a *App) Run() error {
	const op = "restapp.Run"

	slog.Info("Starting rest application on", slog.String("host", a.Host), slog.Int("port", a.Port), slog.String("op", op))

	return a.RESTServer.Start(fmt.Sprintf("%s:%d", a.Host, a.Port))
}

func (a *App) Stop() error {
	return a.RESTServer.Shutdown(context.Background())
}
