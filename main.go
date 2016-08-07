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
	var svc StringService
	svc = stringService{}
	svc = proxyingMiddleware("http://lowercase:80/", logger, ctx)(svc)
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
	)

	countHandler := httptransport.NewServer(
		ctx,
		count,
		decodeCountRequest,
		encodeResponse,
	)

	lowercaseHandler := httptransport.NewServer(
		ctx,
		lowercase,
		decodeLowercaseRequest,
		encodeResponse,
	)

	http.Handle("/uppercase", uppercaseHandler)
	http.Handle("/count", countHandler)
	http.Handle("/lowercase", lowercaseHandler)
	logger.Log(http.ListenAndServe(":1337", nil))
}
