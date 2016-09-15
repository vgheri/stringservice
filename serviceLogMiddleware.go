package main

import (
	"time"

	"golang.org/x/net/context"

	log "github.com/go-kit/kit/log"
)

type appLoggingMiddleware struct {
	logger log.Logger
	next   StringService
}

func loggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next StringService) StringService {
		return appLoggingMiddleware{logger, next}
	}
}

func (mw appLoggingMiddleware) Uppercase(ctx context.Context, s string) (output string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "uppercase",
			"requestID", ctx.Value("requestID"),
			"input", s,
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	output, err = mw.next.Uppercase(ctx, s)
	return
}

func (mw appLoggingMiddleware) Count(ctx context.Context, s string) (n int) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "count",
			"input", "s",
			"count", n,
			"took", time.Since(begin),
		)
	}(time.Now())
	n = mw.next.Count(ctx, s)
	return
}

func (mw appLoggingMiddleware) Lowercase(ctx context.Context, s string) (output string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "lowercase",
			"input", s,
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	output, err = mw.next.Lowercase(ctx, s)
	return
}
