package main

import (
	log "github.com/go-kit/kit/log"

	"golang.org/x/net/context"

	"github.com/go-kit/kit/endpoint"
)

// TransportMiddleware is an abstraction for transport level midlleware
type TransportMiddleware func(endpoint.Endpoint) endpoint.Endpoint

func transportLoggingMiddleware(logger log.Logger) TransportMiddleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			logger.Log("requestID", ctx.Value("requestID"))
			return next(ctx, request)
		}
	}
}
