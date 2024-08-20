package jwt

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/k3k13a1/vk_tarantool_test/internal/core/models"
)

func GenerateToken(user models.User) (string, error) {
	// Чтение приватного ключа из файла
	privateKeyData, err := os.ReadFile("keys/private.pem")
	if err != nil {
		return "", err
	}

	// Парсинг приватного ключа
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		return "", err
	}

	// Создание JWT-токена с некоторыми пользовательскими данными (claims)
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"username": user.Username,                         // Имя пользователя
		"iat":      time.Now().Unix(),                     // Время выпуска токена
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Срок действия токена (24 часа)
	})

	// Подпись токена с использованием приватного ключа
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Функция для проверки JWT-токена с использованием публичного ключа
func VerifyToken(tokenString string) error {
	const op = "core.jwt.verifyToken"

	log := slog.With(
		slog.String("op", op),
	)

	// Чтение публичного ключа из файла
	publicKeyData, err := os.ReadFile("keys/public.pem")
	if err != nil {
		log.Error("can't read public key", slog.String("error", err.Error()))
		return err
	}

	// Парсинг публичного ключа
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		log.Error("can't parse public key", slog.String("error", err.Error()))
		return err
	}

	// Проверка и парсинг токена
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверка используемого метода подписи
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			log.Error("unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		log.Error("can't parse token", slog.String("error", err.Error()))
		return err
	}

	// Проверка валидности токена и извлечение claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("Token is valid. Claims:", claims)
	} else {
		log.Error("invalid token")
		return fmt.Errorf("invalid token")
	}

	return nil
}
