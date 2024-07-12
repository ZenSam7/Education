package api_grpc

import (
	"context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	xForwardedForHeader = "x-forwarded-for"
)

type Metadata struct {
	ClientIP string
}

func (server *Server) extractMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return mtdt
	}

	if clientIPs := md.Get(xForwardedForHeader); len(clientIPs) > 0 {
		mtdt.ClientIP = clientIPs[0]
	}

	if p, ok := peer.FromContext(ctx); ok {
		mtdt.ClientIP = p.Addr.String()
	}

	return mtdt
}
