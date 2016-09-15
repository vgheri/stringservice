package main

import (
	"errors"
	"strings"

	"golang.org/x/net/context"
)

// ServiceMiddleware Chainable service behaviours
type ServiceMiddleware func(StringService) StringService

// StringService models our service
type StringService interface {
	Uppercase(context.Context, string) (string, error)
	Count(context.Context, string) int
	Lowercase(context.Context, string) (string, error)
}

type stringService struct{}

// ErrEmpty is returned when input string is empty
var ErrEmpty = errors.New("Empty string")

/// Business logic
func (stringService) Uppercase(_ context.Context, s string) (string, error) {
	if len(s) == 0 {
		return "", ErrEmpty
	}
	return strings.ToUpper(s), nil
}

func (stringService) Count(_ context.Context, s string) int {
	return len(s)
}

// Placeholder, we won't use it
func (stringService) Lowercase(_ context.Context, s string) (string, error) {
	return "", nil
}
