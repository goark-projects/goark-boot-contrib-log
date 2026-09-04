package properties

import (
	"fmt"
	"log/slog"
	"sort"
	"strings"

	coreenv "goark.dev/goark/core/env"
)

var builtInGroups = map[string][]string{
	"sql": {"goark.dev.database", "goark.dev.orm"},
	"web": {"goark.dev.arkarta", "goark.dev.arkhos", "goark.dev.goark.web"},
}

func readGroups(environment coreenv.Environment) (map[string][]string, error) {
	groups := cloneGroups(builtInGroups)
	values, found, err := coreenv.GetPropertyMapAsValue[string](environment, strings.TrimSuffix(GroupPrefix, "."))
	if err != nil {
		return nil, fmt.Errorf("gbc-log: read logging groups: %w", err)
	}
	if !found {
		return groups, nil
	}
	for name, value := range values {
		name = strings.TrimSpace(name)
		members := splitNonEmpty(value)
		if name == "" || len(members) == 0 {
			return nil, fmt.Errorf("gbc-log: logging group %q has no members", name)
		}
		groups[name] = members
	}
	return groups, nil
}

func readLevels(environment coreenv.Environment, groups map[string][]string) (map[string]slog.Level, error) {
	values, found, err := coreenv.GetPropertyMapAsValue[string](environment, strings.TrimSuffix(LevelPrefix, "."))
	if err != nil {
		return nil, fmt.Errorf("gbc-log: read logger levels: %w", err)
	}
	levels := make(map[string]slog.Level)
	if !found {
		return levels, nil
	}
	groupLevels := make(map[string]slog.Level)
	directLevels := make(map[string]slog.Level)
	for name, value := range values {
		name = strings.TrimSpace(name)
		if name == "" || strings.EqualFold(name, "root") {
			continue
		}
		level, parseErr := parseLevel(LevelPrefix+name, value)
		if parseErr != nil {
			return nil, parseErr
		}
		if _, group := groups[name]; group {
			groupLevels[name] = level
		} else {
			directLevels[name] = level
		}
	}
	for _, name := range sortedLevelNames(groupLevels) {
		for _, member := range groups[name] {
			levels[member] = groupLevels[name]
		}
	}
	for name, level := range directLevels {
		levels[name] = level
	}
	return levels, nil
}

func readStringMap(environment coreenv.Environment, prefix string) (map[string]string, error) {
	values, found, err := coreenv.GetPropertyMapAsValue[string](environment, prefix)
	if err != nil {
		return nil, fmt.Errorf("gbc-log: read %s: %w", prefix, err)
	}
	if !found {
		return map[string]string{}, nil
	}
	result := make(map[string]string, len(values))
	for key, value := range values {
		key = strings.TrimSpace(key)
		if key == "" {
			return nil, fmt.Errorf("gbc-log: %s contains an empty key", prefix)
		}
		result[key] = value
	}
	return result, nil
}

func cloneGroups(source map[string][]string) map[string][]string {
	result := make(map[string][]string, len(source))
	for name, members := range source {
		result[name] = append([]string(nil), members...)
	}
	return result
}

func splitNonEmpty(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if part = strings.TrimSpace(part); part != "" {
			result = append(result, part)
		}
	}
	return result
}

// SortedLevelNames 返回稳定排序后的 Logger 名称。
func SortedLevelNames(levels map[string]slog.Level) []string {
	return sortedLevelNames(levels)
}

func sortedLevelNames(levels map[string]slog.Level) []string {
	names := make([]string, 0, len(levels))
	for name := range levels {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
