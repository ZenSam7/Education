package token

import (
	"github.com/google/uuid"
	"time"
)

// Payload структура полезной нагрузки токена
type Payload struct {
	IDSession uuid.UUID `json:"id_session"`
	IDUser    int32     `json:"id_user"`
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
func NewPayload(IDUser int32, duration time.Duration) (*Payload, error) {
	payload := &Payload{
		IDSession: uuid.New(),
		IDUser:    IDUser,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}
