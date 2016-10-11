package main

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"

	"io/ioutil"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/nu7hatch/gouuid"
	"golang.org/x/net/context"
)

/// Request and Response objects
type uppercaseRequest struct {
	S string `json:"s"`
}

type uppercaseResponse struct {
	S   string `json:"s"`
	Err string `json:"err, omitemtpy"`
}

type countRequest struct {
	S string `json:"s"`
}

type countResponse struct {
	V int `json:"v"`
}

type lowercaseRequest struct {
	S string `json:"s"`
}

type lowercaseResponse struct {
	S   string `json:"s"`
	Err string `json:"err, omitemtpy"`
}

/// Endpoints definition
func makeUppercaseEndpoint(svc StringService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(uppercaseRequest)
		res, err := svc.Uppercase(ctx, req.S)
		if err != nil {
			return uppercaseResponse{req.S, err.Error()}, nil
		}
		return uppercaseResponse{res, ""}, nil
	}
}

func makeCountEndpoint(svc StringService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(countRequest)
		res := svc.Count(ctx, req.S)
		return countResponse{res}, nil
	}
}

func makeLowercaseEndpoint(svc StringService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(lowercaseRequest)
		res, err := svc.Lowercase(ctx, req.S)
		if err != nil {
			return lowercaseResponse{req.S, err.Error()}, nil
		}
		return lowercaseResponse{res, ""}, nil
	}
}

func encodeRequest(ctx context.Context, r *http.Request, req interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(req); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

// Endpoints parameters encode/decode
func decodeUppercaseRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	var request uppercaseRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeCountRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	var request countRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeLowercaseRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	var request lowercaseRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeLowercaseResponse(_ context.Context, w *http.Response) (interface{}, error) {
	var response lowercaseResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, res interface{}) error {
	return json.NewEncoder(w).Encode(res)
}

/*
	Server before functions, executed on the HTTP request before req is decoded
*/

func setRequestIDInContext() httptransport.RequestFunc {
	return func(ctx context.Context, request *http.Request) context.Context {
		reqID := request.Header.Get("X-Request-ID")
		if reqID == "" {
			u, err := uuid.NewV4()
			if err == nil {
				reqID = u.String()
			}
		}
		request.Header.Set("X-Request-ID", reqID)
		return context.WithValue(ctx, "requestID", reqID)
	}
}

func setClientIPInContext() httptransport.RequestFunc {
	return func(ctx context.Context, request *http.Request) context.Context {
		ip := request.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip, _, _ = net.SplitHostPort(request.RemoteAddr)
		}
		return context.WithValue(ctx, "clientIP", ip)
	}
}

/*
	Client before functions, applied to the outgoing HTTP requests
*/

func addRequestIDtoOutgoingHTTPRequest() httptransport.RequestFunc {
	return func(ctx context.Context, request *http.Request) context.Context {
		request.Header.Set("X-Request-ID", ctx.Value("requestID").(string))
		return ctx
	}
}

func addHeaderToOutgoingHTTPRequest(key, value string) httptransport.RequestFunc {
	return func(ctx context.Context, request *http.Request) context.Context {
		request.Header.Set(key, value)
		return ctx
	}
}
