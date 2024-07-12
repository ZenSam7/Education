package api

import (
	"errors"
	"github.com/ZenSam7/Education/token"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	// SupportedAuthType Все типы авторизации которые поддерживает API
	SupportedAuthType = "bearer"
	// AuthPayloadKey Ключ для взятия токена
	AuthPayloadKey = "authorization_payload"
)

// authMiddleware Создаём промежуточную функцию авторизации
func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if len(authHeader) == 0 {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				errorResponse(errors.New("заголовок авторизации не предоставлен")),
			)
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				errorResponse(errors.New("неправильный заголовок токена")),
			)
			return
		}

		authType := strings.ToLower(fields[0])
		if authType != SupportedAuthType {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				errorResponse(errors.New("тип авторизации пока не поддерживается")),
			)
		}

		authToken := fields[1]
		payload, err := tokenMaker.VerifyToken(authToken)
		if err != nil {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				errorResponse(err),
			)
			return
		}

		// Всё проверели, значит теперь токен действителен
		ctx.Set(AuthPayloadKey, payload)
		ctx.Next()
	}
}
