package api

import (
	"errors"
	"github.com/ZenSam7/Education/token"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// authTypes Все типы авторизации которые поддерживает API
const (
	supportedAuthType = "bearer"
	authPayloadKey    = "authorization_payload"
)

// authMiddleware Создаём промежуточную функцию авторизации
func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if len(authHeader) == 0 {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				errorResponse(errors.New("authorization header is not provided")),
			)
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				errorResponse(errors.New("wrong authorization header")),
			)
			return
		}

		authType := strings.ToLower(fields[0])
		if authType != supportedAuthType {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				errorResponse(errors.New("unsupported authorization type")),
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
		ctx.Set(authPayloadKey, payload)
		ctx.Next()
	}
}
