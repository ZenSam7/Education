package api

import (
	"context"
	"fmt"
	"github.com/ZenSam7/Education/token"
)

func (server *Server) authUser(ctx context.Context) (*token.Payload, error) {
	mtdt, err := server.extractMetadata(ctx)
	if err != nil {
		return nil, err
	}

	if len(mtdt.AccessToken) == 0 {
		return nil, fmt.Errorf("не указан токен авторизации")
	}

	payload, err := server.tokenMaker.VerifyToken(mtdt.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("указан неправильный токен авторизации: %s", err)
	}

	return payload, nil
}
