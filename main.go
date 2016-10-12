package main

import (
	"net/http"
	"os"

	log "github.com/go-kit/kit/log"

	httptransport "github.com/go-kit/kit/transport/http"
	"golang.org/x/net/context"
)

func main() {
	ctx := context.Background()
	logger := log.NewLogfmtLogger(os.Stderr)
	// Intercepting function which injects the request ID into the request context.
	option := httptransport.ServerBefore(setRequestIDInContext(),
		setClientIPInContext())

	var svc StringService
	svc = stringService{}
	svc = proxyingMiddleware("http://lowercase/", logger)(svc)
	svc = loggingMiddleware(logger)(svc)

	uppercase := makeUppercaseEndpoint(svc)
	uppercase = transportLoggingMiddleware(log.NewContext(logger).With("method", "uppercase"))(uppercase)

	count := makeCountEndpoint(svc)
	count = transportLoggingMiddleware(log.NewContext(logger).With("method", "count"))(count)

	lowercase := makeLowercaseEndpoint(svc)
	lowercase = transportLoggingMiddleware(log.NewContext(logger).With("method", "lowercase"))(lowercase)

	uppercaseHandler := httptransport.NewServer(
		ctx,
		uppercase,
		decodeUppercaseRequest,
		encodeResponse,
		option,
	)

	countHandler := httptransport.NewServer(
		ctx,
		count,
		decodeCountRequest,
		encodeResponse,
		option,
	)

	lowercaseHandler := httptransport.NewServer(
		ctx,
		lowercase,
		decodeLowercaseRequest,
		encodeResponse,
		option,
	)

	http.Handle("/uppercase", uppercaseHandler)
	http.Handle("/count", countHandler)
	http.Handle("/lowercase", lowercaseHandler)
	logger.Log(http.ListenAndServe(":1337", nil))
}
