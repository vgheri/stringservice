package main

import (
	"errors"
	"strings"
)

// ServiceMiddleware Chainable service behaviours
type ServiceMiddleware func(StringService) StringService

// StringService models our service
type StringService interface {
	Uppercase(string) (string, error)
	Count(string) int
}

type stringService struct{}

// ErrEmpty is returned when input string is empty
var ErrEmpty = errors.New("Empty string")

/// Business logic
func (stringService) Uppercase(s string) (string, error) {
	if len(s) == 0 {
		return "", ErrEmpty
	}
	return strings.ToUpper(s), nil
}

func (stringService) Count(s string) int {
	return len(s)
}
