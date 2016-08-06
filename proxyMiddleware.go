package main

import (
	"errors"
	"net/url"
	"strings"

	"github.com/go-kit/kit/endpoint"

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

func proxyingMiddleware(proxyURL string, ctx context.Context) ServiceMiddleware {
	return func(next StringService) StringService {
		return proxy{ctx, next, makeLowercaseProxy(proxyURL, ctx)}
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
