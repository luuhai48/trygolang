package main

import "os"

func MustGetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	panic("Environment not found: " + key)
}
