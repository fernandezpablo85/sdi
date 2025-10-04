package env

import (
	"log/slog"
	"os"
	"strconv"
)

func GetOrElse(key, orElse string) string {
	param, found := os.LookupEnv(key)
	if found {
		return param
	}
	return orElse
}

func GetIntOrElse(key string, orElse int) int {
	param, found := os.LookupEnv(key)
	if !found {
		return orElse
	}

	intParam, err := strconv.Atoi(param)
	if err != nil {
		slog.Warn("error parsing int env var", "key", key, "value", param, "error", err)
		return orElse
	}

	return intParam
}
