package handlers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/k3k13a1/vk_tarantool_test/internal/core/jwt"
	"github.com/k3k13a1/vk_tarantool_test/internal/services"
	"github.com/labstack/echo/v4"
)

type TTHandlers interface {
	Login(
		ctx context.Context,
		login string,
		pass string,
	) (string, error)
	Write(
		ctx context.Context,
		data map[string]interface{},
	) error
	Read(
		ctx context.Context,
		data []string,
	) (map[string]interface{}, error)
}

func checkAuthorization(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "заголовок Authorization отсутствует")
	}

	// Проверяем, начинается ли заголовок с "Bearer "
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		return echo.NewHTTPError(http.StatusUnauthorized, "неправильный формат токена")
	}

	// Валидация на пустой токен.
	if token == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "токен отсутствует")
	}

	if jwt.VerifyToken(token) != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "неправильный токен")
	}

	return nil
}

func Login(c echo.Context, a TTHandlers) error {
	const op = "handlers.Login"

	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "" || password == "" {
		slog.Error(fmt.Sprintf("%s: empty username or password", op))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "имя пользователя и пароль обязательны"})
	}

	token, err := a.Login(c.Request().Context(), username, password)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: %v", op, err))
		if errors.As(err, &services.ErrInvalidCredentials) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "неверное имя пользователя или пароль"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, token)
}

func Write(c echo.Context, a TTHandlers) error {
	const op = "handlers.Write"

	if err := checkAuthorization(c); err != nil {
		return err
	}

	type reqData struct {
		Data map[string]interface{} `json:"data"`
	}

	req := new(reqData)
	if err := c.Bind(req); err != nil {
		slog.Error(fmt.Sprintf("%s: %v", op, err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := a.Write(c.Request().Context(), req.Data); err != nil {
		slog.Error(fmt.Sprintf("%s: %v", op, err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "success"})
}

func Read(c echo.Context, a TTHandlers) error {
	const op = "handlers.Read"

	if err := checkAuthorization(c); err != nil {
		return err
	}

	// Структура для данных запроса
	type reqData struct {
		Keys []string `json:"keys"`
	}

	req := new(reqData)
	if err := c.Bind(req); err != nil {
		slog.Error(fmt.Sprintf("%s: %v", op, err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	slog.Info(fmt.Sprintf("%s: %v", op, req.Keys))

	data, err := a.Read(c.Request().Context(), req.Keys)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: %v", op, err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	type ResData struct {
		Data map[string]interface{} `json:"data"`
	}

	return c.JSON(http.StatusOK, ResData{Data: data})
}
