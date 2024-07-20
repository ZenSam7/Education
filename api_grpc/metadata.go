package api_grpc

import (
	"context"
	"fmt"
	"github.com/ZenSam7/Education/token"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"strings"
)

const (
	xForwardedForHeader = "x-forwarded-for"
	authHeader          = "authorization"
)

type Metadata struct {
	ClientIP    string
	AccessToken string
}

func (server *Server) extractMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return mtdt
	}

	// Запасной вариант для получения ip
	if p, ok := peer.FromContext(ctx); ok {
		mtdt.ClientIP = p.Addr.String()
	}

	if clientIPs := md.Get(xForwardedForHeader); len(clientIPs) != 0 {
		mtdt.ClientIP = clientIPs[0]
	}

	if tkn := md.Get(authHeader); len(tkn) != 0 {
		mtdt.AccessToken = strings.Split(tkn[0], " ")[1]
	}

	return mtdt
}

func (server *Server) authUser(ctx context.Context) (*token.Payload, error) {
	info := server.extractMetadata(ctx)
	if len(info.AccessToken) == 0 {
		return nil, fmt.Errorf("не указан токен авторизации")
	}

	payload, err := server.tokenMaker.VerifyToken(info.AccessToken)
	if err != nil {
		return nil, err
	}

	return payload, nil
}
