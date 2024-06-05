package token

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrorInvalidToken = fmt.Errorf("токен недействителен")
	ErrorExpiredToken = errors.New("срок действия токена истёк")
)

// Maker интерфейс для управления токенами
type Maker interface {
	CreateToken(IDUser int32, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}
