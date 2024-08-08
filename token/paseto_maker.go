package token

import (
	"github.com/ZenSam7/Education/tools"
	"github.com/o1egl/paseto"
	"github.com/rs/zerolog/log"
	"time"
)

// PasetoMaker Реализыем то же самое что и JWTMaker, но на paseto
type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(secretKey string) Maker {
	if len(secretKey) < minSecretKeySize {
		secretKey = tools.GetRandomString(minSecretKeySize)[:minSecretKeySize]
		log.Warn().Msgf("длина secretKey < %d, secretKey заменён на случайную строку", minSecretKeySize)
	}

	newPaseto := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(secretKey),
	}
	return newPaseto
}

func (maker *PasetoMaker) CreateToken(idUser int32, role string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(idUser, role, duration)
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
