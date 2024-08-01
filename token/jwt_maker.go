package token

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const minSecretKeySize = 32

// JWTMaker Просто реализуем интерфейс Maker
type JWTMaker struct {
	secretKey string
}

func (maker *JWTMaker) CreateToken(IDUser int32, role string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(IDUser, role, duration)
	if err != nil {
		return "", payload, err
	}

	// Создаём сам токен, подавая метод подписи и полезную нагрузку
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	return token, payload, err
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	// Проверяем заголовок чтобы алгоритм подписи был такой какой нужен
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		// Преобразуем метод подписи в HMSC (т.к. наш метод подписи (jwt.SigningMethodHS256) реализует SigningMethodHMAC)
		// (ok - обозначает что конвертирование удалось)
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			// Что-то не так с подписью
			return nil, ErrorInvalidToken
		}
		// Всё гуд, возвращаем ключ
		return []byte(maker.secretKey), nil
	}

	// Создаём пустой токен с функцией проверки метода подписи
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		return nil, err
	}

	// Получаем полезные данные из токена (преобраем в Payload)
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrorInvalidToken
	}

	return payload, nil
}

// NewJWTMaker Создаём сам JWT с секретным ключом
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("секретный ключ должен содержать не менее %d", minSecretKeySize)
	}
	return &JWTMaker{secretKey: secretKey}, nil
}
