package token

import (
	"fmt"
	"time"
)

var (
	ErrorInvalidToken      = fmt.Errorf("токен недействителен")
	ErrorExpiredToken      = fmt.Errorf("срок действия токена истёк")
	ErrorInvalidPermission = fmt.Errorf("остутствуют нужные права")
)

// Maker интерфейс для управления токенами
type Maker interface {
	CreateToken(idUser int32, role string, duration time.Duration) (string, *Payload, error)
	VerifyToken(Token string) (*Payload, error) /* "Token" обязательно с заглавной буквы
	(ну, или любое другое название кроме "token", т.к. в mock'е это название будет заменено
	на пакет, и, короче, моки перестанут работать) */
}
