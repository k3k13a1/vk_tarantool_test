package services

import (
	"context"
	"log/slog"
	"sync"

	"github.com/k3k13a1/vk_tarantool_test/internal/core/jwt"
	"github.com/k3k13a1/vk_tarantool_test/internal/core/models"
	tt "github.com/k3k13a1/vk_tarantool_test/internal/storage/tarantool"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	UserProvider
	TTProvider
}

type UserProvider interface {
	User(
		ctx context.Context,
		login string,
	) (models.User, error)
}

type TTProvider interface {
	Read(
		ctx context.Context,
		key string,
	) (interface{}, error)

	Write(
		ctx context.Context,
		key string,
		value interface{},
	) error
}

func New(
	storage *tt.Storage,
) *Service {
	return &Service{
		UserProvider: storage,
		TTProvider:   storage,
	}
}

func (s *Service) Login(ctx context.Context, login, password string) (string, error) {
	const op = "services.Login"

	log := slog.With(
		slog.String("op", op),
		slog.String("login", login),
	)

	log.Info("attempting to login user")

	user, err := s.UserProvider.User(ctx, login)
	if err != nil {
		log.Error("user not found", slog.String("error", err.Error()))
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Error("incorrect password", slog.String("error", err.Error()))
		return "", err
	}

	log.Info("successfully logged in user")

	token, err := jwt.GenerateToken(user)
	if err != nil {
		log.Error("failed to generate token", slog.String("error", err.Error()))
	}

	return token, nil
}

func (s *Service) Write(ctx context.Context, data map[string]interface{}) error {
	const op = "services.Write"

	log := slog.With(
		slog.String("op", op),
	)

	var wg sync.WaitGroup
	for key, value := range data {
		wg.Add(1)
		go func(k string, v interface{}) {
			defer wg.Done()
			log.Info("writing data", slog.String("key", k))

			s.TTProvider.Write(ctx, k, v)

		}(key, value)
	}

	// Ожидание завершения всех горутин
	wg.Wait()

	return nil
}

func (s *Service) Read(ctx context.Context, data []string) (map[string]interface{}, error) {
	const op = "services.Read"

	log := slog.With(
		slog.String("op", op),
	)

	// Используем sync.Map для параллельного чтения
	var wg sync.WaitGroup
	results := sync.Map{}

	for _, key := range data {
		wg.Add(1)
		go func(k string) {
			defer wg.Done()

			res, err := s.TTProvider.Read(ctx, k)
			if err != nil {
				log.Info("Ошибка чтения данных для ключа %s: %v", k, err)
				return
			}

			log.Info("Запись данных для ключа %s: %v", k, res)
			results.Store(k, res)

		}(key)
	}

	// Ожидание завершения всех горутин
	wg.Wait()

	// Формирование ответа
	result := make(map[string]interface{})
	results.Range(func(key, value interface{}) bool {
		result[key.(string)] = value
		return true
	})

	return result, nil
}
