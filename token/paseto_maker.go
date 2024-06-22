package token

import (
	"fmt"
	"github.com/o1egl/paseto"
	"time"
)

// PasetoMaker Реализыем то же самое что и JWTMaker, но на paseto
type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("секретный ключ должен содержать не менее %d", minSecretKeySize)
	}

	newPaseto := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(secretKey),
	}
	return newPaseto, nil
}

func (maker *PasetoMaker) CreateToken(IDUser int32, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(IDUser, duration)
	if err != nil {
		return "", payload, err
	}

	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	return token, payload, err
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	// Очень легко и просто проверяем токен
	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrorInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
