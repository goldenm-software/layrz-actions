package env

import (
	"os"
	"strconv"
)

func Str(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func Bool(key string) bool {
	return os.Getenv(key) == "true"
}

func Int(key string) (int, error) {
	return strconv.Atoi(os.Getenv(key))
}
