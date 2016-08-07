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
	ctx       context.Context
	next      StringService
	lowercase endpoint.Endpoint
}

func (mw proxy) Uppercase(s string) (string, error) {
	return mw.next.Uppercase(s)
}

func (mw proxy) Count(s string) int {
	return mw.next.Count(s)
}

func (mw proxy) Lowercase(s string) (string, error) {
	response, err := mw.lowercase(mw.ctx, lowercaseRequest{S: s})
	if err != nil {
		return "", err
	}
	resp := response.(lowercaseResponse)
	if resp.Err != "" {
		return resp.S, errors.New(resp.Err)
	}
	return resp.S, nil
}

func proxyingMiddleware(proxyURL string, logger log.Logger, ctx context.Context) ServiceMiddleware {
	return func(next StringService) StringService {
		lowercaseEndpoint := makeLowercaseProxy(proxyURL, ctx)
		st := gobreaker.Settings{}
		st.Name = "StringServiceToLowercase"
		st.OnStateChange = func(name string, from gobreaker.State, to gobreaker.State) {
			logger.Log("CircuitBreaker", "from "+from.String()+" to "+to.String())
		}
		lowercaseEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(st))(lowercaseEndpoint)
		return proxy{ctx, next, lowercaseEndpoint}
	}
}

func makeLowercaseProxy(proxyURL string, ctx context.Context) endpoint.Endpoint {
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
