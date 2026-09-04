package properties

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	coreenv "goark.dev/goark/core/env"
	goarklog "goark.dev/log"
)

func optionalString(environment coreenv.Environment, key string) string {
	value, _ := environment.GetProperty(key)
	return strings.TrimSpace(value)
}

func stringWithDefault(environment coreenv.Environment, key string, fallback string) string {
	if value, found := environment.GetProperty(key); found {
		return strings.TrimSpace(value)
	}
	return fallback
}

func boolWithDefault(environment coreenv.Environment, key string, fallback bool) (bool, error) {
	value, found := environment.GetProperty(key)
	if !found {
		return fallback, nil
	}
	parsed, err := strconv.ParseBool(strings.TrimSpace(value))
	if err != nil {
		return false, fmt.Errorf("gbc-log: invalid boolean for %q: %w", key, err)
	}
	return parsed, nil
}

func optionalBool(environment coreenv.Environment, key string) (*bool, error) {
	value, found := environment.GetProperty(key)
	if !found {
		return nil, nil
	}
	parsed, err := strconv.ParseBool(strings.TrimSpace(value))
	if err != nil {
		return nil, fmt.Errorf("gbc-log: invalid boolean for %q: %w", key, err)
	}
	return &parsed, nil
}

func intWithDefault(environment coreenv.Environment, key string, fallback int) (int, error) {
	value, found := environment.GetProperty(key)
	if !found {
		return fallback, nil
	}
	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || parsed < 0 {
		return 0, fmt.Errorf("gbc-log: %q must be a non-negative integer", key)
	}
	return parsed, nil
}

func optionalInt(environment coreenv.Environment, key string) (*int, error) {
	value, found := environment.GetProperty(key)
	if !found {
		return nil, nil
	}
	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || parsed < 0 {
		return nil, fmt.Errorf("gbc-log: %q must be a non-negative integer", key)
	}
	return &parsed, nil
}

func dataSizeWithDefault(environment coreenv.Environment, key string, fallback int64) (int64, error) {
	value, found := environment.GetProperty(key)
	if !found {
		return fallback, nil
	}
	parsed, err := goarklog.ParseByteSize(value)
	if err != nil {
		return 0, fmt.Errorf("gbc-log: invalid data size for %q: %w", key, err)
	}
	return parsed, nil
}

func optionalLevel(environment coreenv.Environment, key string) (*slog.Level, error) {
	value, found := environment.GetProperty(key)
	if !found {
		return nil, nil
	}
	level, err := parseLevel(key, value)
	if err != nil {
		return nil, err
	}
	return &level, nil
}

func levelWithDefault(environment coreenv.Environment, key string, fallback slog.Level) (*slog.Level, error) {
	level, err := optionalLevel(environment, key)
	if err != nil {
		return nil, err
	}
	if level != nil {
		return level, nil
	}
	return levelPointer(fallback), nil
}

func levelPointer(level slog.Level) *slog.Level {
	copy := level
	return &copy
}

func parseLevel(key string, value string) (slog.Level, error) {
	level, err := goarklog.ParseLevel(strings.TrimSpace(value))
	if err != nil {
		return 0, fmt.Errorf("gbc-log: invalid level for %q: %w", key, err)
	}
	return level, nil
}
