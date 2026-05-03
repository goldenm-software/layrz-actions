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

func Float(key string, defaultValue float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return defaultValue
}
