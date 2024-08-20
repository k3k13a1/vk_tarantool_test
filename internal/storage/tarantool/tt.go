package tt

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/k3k13a1/vk_tarantool_test/internal/core/models"
	"github.com/k3k13a1/vk_tarantool_test/internal/storage"
	"github.com/tarantool/go-tarantool/v2"
	"golang.org/x/crypto/bcrypt"
)

type Storage struct {
	db *tarantool.Connection
}

func New(login, password, host string, port int) (*Storage, error) {
	const op = "storage.tarantool.New"

	log := slog.With(
		slog.String("op", op),
	)

	addrStr := fmt.Sprintf("%s:%d", host, port)

	dialer := tarantool.NetDialer{
		Address:  addrStr,
		User:     login,
		Password: password,
	}
	opts := tarantool.Opts{
		Timeout: time.Second,
	}

	conn, err := tarantool.Connect(context.Background(), dialer, opts)
	if err != nil {
		log.Debug("Ошибка подключения к Tarantool: %v", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: conn}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) User(ctx context.Context, login string) (models.User, error) {
	const op = "storage.tarantool.Login"

	log := slog.With(
		slog.String("op", op),
	)

	// Заглушка для тестового пользователя
	if login == storage.TestUserName {
		hash, err := bcrypt.GenerateFromPassword([]byte(storage.TestUserPass), bcrypt.DefaultCost)
		if err != nil {
			log.Error("Не удалось создать хэш пароля для пользователя %s: %v", login, err)
			return models.User{}, fmt.Errorf("%s: %w", op, err)
		}
		return models.User{Username: login, PassHash: hash}, nil
	} else {
		return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrInvalidCredentials)
	}
}

func (s *Storage) Write(ctx context.Context, key string, value interface{}) error {
	const op = "storage.tarantool.Write"

	log := slog.With(
		slog.String("op", op),
	)

	// _, err := s.db.Do(
	// 	tarantool.NewCallRequest("vshard.router.callrw").
	// 		Args([]interface{}{"vk_test", "upsert", []interface{}{key, value}}),
	// ).Get()

	_, err := s.db.Do(
		tarantool.NewUpsertRequest("vk_test").Tuple([]interface{}{key, value}).Operations(tarantool.NewOperations().Assign(1, value)),
	).Get()

	if err != nil {
		log.Error("Ошибка записи данных %s: %v", key, err)
	}

	log.Info("Запись данных %s: %v", key, value)

	return nil
}

func (s *Storage) Read(ctx context.Context, key string) (interface{}, error) {
	const op = "storage.tarantool.Read"

	log := slog.With(
		slog.String("op", op),
	)

	data, err := s.db.Do(
		tarantool.NewSelectRequest("vk_test").Key([]interface{}{key}),
	).Get()

	// data, err := s.db.Do(
	// 	tarantool.NewCallRequest("vshard.router.callro").
	// 		Args([]interface{}{"vk_test", "select", []interface{}{key}}),
	// ).Get()

	if err != nil {
		log.Error("Ошибка чтения данных %s: %v", key, err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Чтение данных %s: %s", key, key)

	log.Info(fmt.Sprintf("%v", data[0].([]interface{})[1]))

	return data[0].([]interface{})[1], nil
}
