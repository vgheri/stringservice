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
	svc = appLoggingMiddleware{logger, svc}

	uppercase := makeUppercaseEndpoint(svc)
	uppercase = transportLoggingMiddleware(log.NewContext(logger).With("method", "uppercase"))(uppercase)

	count := makeCountEndpoint(svc)
	count = transportLoggingMiddleware(log.NewContext(logger).With("method", "count"))(count)
	//
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

	http.Handle("/uppercase", uppercaseHandler)
	http.Handle("/count", countHandler)
	logger.Log(http.ListenAndServe(":1337", nil))
}
