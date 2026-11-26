// Package utils provides common utility functions
package utils

import (
	"errors"
	"os"
	"strconv"
)

// GetEnv returns the value of the environment variable key, or fallback if it's not set
func GetEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// GetEnvInt returns the environment variable key as a positive int, or fallback if it's unset, malformed or not positive
func GetEnvInt(key string, fallback int) int {
	if n, err := strconv.Atoi(os.Getenv(key)); err == nil && n > 0 {
		return n
	}
	return fallback
}

// GetEnvOrError returns the value of the environment variable key, or an error if it's not set
func GetEnvOrError(key string) (string, error) {
	if value := os.Getenv(key); value == "" {
		return "", errors.New("environment variable " + key + " not set")
	}
	return os.Getenv(key), nil
}

// GetEnvOrErrorInProduction returns the value of the environment variable key.
// In production, it returns an error if not set. In dev, it returns devFallback.
func GetEnvOrErrorInProduction(key string, devFallback string) (string, error) {
	if !IsProduction() {
		return GetEnv(key, devFallback), nil
	}
	return GetEnvOrError(key)
}
