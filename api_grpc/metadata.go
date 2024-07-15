package api_grpc

import (
	"context"
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

	if token := md.Get(authHeader); len(token) != 0 {
		mtdt.AccessToken = strings.Split(token[0], " ")[1]
	}

	return mtdt
}
