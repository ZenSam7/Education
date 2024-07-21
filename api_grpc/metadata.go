package api_grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"strings"
)

const (
	// supportedAuthType Все типы авторизации которые поддерживает API
	supportedAuthType   = "Bearer"
	xForwardedForHeader = "x-forwarded-for"
	userAgentHeader     = "user-agent"
	authHeader          = "authorization"
)

type Metadata struct {
	ClientIP    string
	ClientAgent string
	AccessToken string
}

func (server *Server) extractMetadata(ctx context.Context) (*Metadata, error) {
	mtdt := &Metadata{}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("ошибка в извлечении метаданных из контекста")
	}

	if clientIPs := md.Get(xForwardedForHeader); len(clientIPs) != 0 {
		mtdt.ClientIP = clientIPs[0]
	}

	if clientAgents := md.Get(userAgentHeader); len(clientAgents) != 0 {
		mtdt.ClientAgent = clientAgents[0]
	}

	if tkn := md.Get(authHeader); len(tkn) > 0 {
		tokenInfo := strings.Fields(tkn[0])
		if len(tokenInfo) < 2 {
			return nil, fmt.Errorf("неправильный формат токена авторизации")
		} else if tokenInfo[0] != supportedAuthType {
			return nil, fmt.Errorf("неподдерживаемый тип токена: %s", tokenInfo[0])
		}

		mtdt.AccessToken = tokenInfo[1]
	}

	return mtdt, nil
}
