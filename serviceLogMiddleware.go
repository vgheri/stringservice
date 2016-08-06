package main

import (
	"time"

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

func (mw appLoggingMiddleware) Uppercase(s string) (output string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "uppercase",
			"input", s,
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	output, err = mw.next.Uppercase(s)
	return
}

func (mw appLoggingMiddleware) Count(s string) (n int) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "count",
			"input", "s",
			"count", n,
			"took", time.Since(begin),
		)
	}(time.Now())
	n = mw.next.Count(s)
	return
}

func (mw appLoggingMiddleware) Lowercase(s string) (output string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "lowercase",
			"input", s,
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	output, err = mw.next.Lowercase(s)
	return
}
