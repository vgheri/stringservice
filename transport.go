package main

import (
	"bytes"
	"encoding/json"
	"net/http"

	"io/ioutil"

	"github.com/go-kit/kit/endpoint"
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

/// Endpoints definition
func makeUppercaseEndpoint(svc StringService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(uppercaseRequest)
		res, err := svc.Uppercase(req.S)
		if err != nil {
			return uppercaseResponse{req.S, err.Error()}, nil
		}
		return uppercaseResponse{res, ""}, nil
	}
}

func makeCountEndpoint(svc StringService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(countRequest)
		res := svc.Count(req.S)
		return countResponse{res}, nil
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

func decodeUppercaseResponse(_ context.Context, w *http.Response) (interface{}, error) {
	var response uppercaseResponse
	if err := json.NewDecoder(w.Body).Decode(response); err != nil {
		return nil, err
	}
	return response, nil
}

func decodeCountRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	var request countRequest
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, res interface{}) error {
	return json.NewEncoder(w).Encode(res)
}
