package main

import (
	"errors"
	"net/url"
	"strings"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/sony/gobreaker"

	log "github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"golang.org/x/net/context"
)

type proxy struct {
	next      StringService
	lowercase endpoint.Endpoint
}

func (mw proxy) Uppercase(ctx context.Context, s string) (string, error) {
	return mw.next.Uppercase(ctx, s)
}

func (mw proxy) Count(ctx context.Context, s string) int {
	return mw.next.Count(ctx, s)
}

func (mw proxy) Lowercase(ctx context.Context, s string) (string, error) {
	response, err := mw.lowercase(ctx, lowercaseRequest{S: s})
	if err != nil {
		return "", err
	}
	resp := response.(lowercaseResponse)
	if resp.Err != "" {
		return resp.S, errors.New(resp.Err)
	}
	return resp.S, nil
}

func proxyingMiddleware(proxyURL string, logger log.Logger) ServiceMiddleware {
	return func(next StringService) StringService {
		lowercaseEndpoint := makeLowercaseProxy(proxyURL)
		st := gobreaker.Settings{}
		st.Name = "StringServiceToLowercase"
		st.OnStateChange = func(name string, from gobreaker.State, to gobreaker.State) {
			logger.Log("CircuitBreaker", "from "+from.String()+" to "+to.String())
		}
		lowercaseEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(st))(lowercaseEndpoint)
		return proxy{next, lowercaseEndpoint}
	}
}

func makeLowercaseProxy(proxyURL string) endpoint.Endpoint {
	if !strings.HasPrefix(proxyURL, "http") {
		proxyURL = "http://" + proxyURL
	}
	u, err := url.Parse(proxyURL)
	if err != nil {
		panic(err)
	}
	return httptransport.NewClient(
		"GET",
		u,
		encodeRequest,
		decodeLowercaseResponse,
	).Endpoint()
}
