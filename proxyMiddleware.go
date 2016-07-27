package main

import (
	"errors"
	"net/url"

	"github.com/go-kit/kit/endpoint"

	httptransport "github.com/go-kit/kit/transport/http"
	"golang.org/x/net/context"
)

type proxy struct {
	ctx       context.Context
	next      StringService
	uppercase endpoint.Endpoint
}

func (mw proxy) Uppercase(s string) (string, error) {
	response, err := mw.uppercase(mw.ctx, uppercaseRequest{s})
	if err != nil {
		return "", err
	}
	resp := response.(uppercaseResponse)
	if resp.Err != "" {
		return "", errors.New(resp.Err)
	}
	return resp.S, nil
}

func (mw proxy) Count(s string) int {
	return mw.next.Count(s)
}

func proxyingMiddleware(proxyURL string, ctx context.Context) ServiceMiddleware {
	return func(next StringService) StringService {
		return proxy{ctx, next, makeUppercaseProxy(proxyURL, ctx)}
	}
}

func makeUppercaseProxy(proxyURL string, ctx context.Context) endpoint.Endpoint {
	u, err := url.Parse(proxyURL)
	if err != nil {
		panic(err)
	}
	if u.Path == "" {
		u.Path = "/uppercase"
	}
	return httptransport.NewClient(
		"GET",
		u,
		encodeRequest,
		decodeUppercaseResponse,
	).Endpoint()
}
