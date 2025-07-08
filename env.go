package envreader

import (
	"fmt"
	"os"
	"strconv"
)

func ReadEnv[T any](key string, defaultValue T) (T, error) {
	envValue := os.Getenv(key)

	if envValue == "" {
		return defaultValue, nil
	}

	var result T
	switch any(result).(type) {
	case int:
		val, err := strconv.Atoi(envValue)
		if err != nil {
			return defaultValue, fmt.Errorf("failed to convert %q to int: %w", envValue, err)
		}
		return any(val).(T), nil
	case int64:
		val, err := strconv.ParseInt(envValue, 10, 64)
		if err != nil {
			return defaultValue, fmt.Errorf("failed to convert %q to int64: %w", envValue, err)
		}
		return any(val).(T), nil
	case string:
		return any(envValue).(T), nil
	case bool:
		val, err := strconv.ParseBool(envValue)
		if err != nil {
			return defaultValue, fmt.Errorf("failed to convert %q to bool: %w", envValue, err)
		}
		return any(val).(T), nil
	case float64:
		val, err := strconv.ParseFloat(envValue, 64)
		if err != nil {
			return defaultValue, fmt.Errorf("failed to convert %q to float64: %w", envValue, err)
		}
		return any(val).(T), nil
	}

	return defaultValue, fmt.Errorf("unsupported type for environment variable conversion: %T", defaultValue)
}
