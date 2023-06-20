package gapi

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type Metadata struct {
	ClientIp  string
	UserAgent string
}

const (
	grpcGatewayUserAgent = "grpcgateway-user-agent"
	xForwardedFor        = "x-forwarded-for"
	userAgentHeader      = "user-agent"
)

func (server *Server) extractMetada(ctx context.Context) *Metadata {
	meta_data := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		fmt.Printf("metadata [%v]", md)
		if userAgent := md.Get(userAgentHeader); len(userAgent) > 0 {
			meta_data.UserAgent = userAgent[0]
		}
		if userAgent := md.Get(grpcGatewayUserAgent); len(userAgent) > 0 {
			meta_data.UserAgent = userAgent[0]
		}
		if clientIp := md.Get(xForwardedFor); len(clientIp) > 0 {
			meta_data.ClientIp = clientIp[0]
		}
		// map[grpcgateway-accept:[*/*]
		//grpcgateway-content-type:[application/x-www-form-urlencoded]
		//grpcgateway-user-agent:[curl/7.87.0]
		//	x-forwarded-for:[127.0.0.1]
		//	x-forwarded-host:[127.0.0.1:8888]]
	}
	if p, ok := peer.FromContext(ctx); ok {
		meta_data.ClientIp = p.Addr.String()
	}
	return meta_data
}
