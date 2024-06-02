package token

import (
	"github.com/google/uuid"
	"time"
)

// Payload структура полезной нагрузки токена
type Payload struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	// IssuedAt Когда токен создался
	IssuedAt time.Time `json:"issued_at"`
	// ExpiredAt Когда токен просрочится
	ExpiredAt time.Time `json:"expired_at"`
}

func (p *Payload) Valid() error {
	// Если время токена истекло
	if time.Now().After(p.ExpiredAt) {
		return ErrorExpiredToken
	}
	return nil
}

// NewPayload Просто создаём новый токен
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}
