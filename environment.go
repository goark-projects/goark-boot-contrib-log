package gbclog

import (
	"fmt"
	"log/slog"
	"sort"
	"strings"

	coreenv "goark.dev/goark/core/env"
	goarklog "goark.dev/log"
)

type loggingProperties struct {
	rootLevel        *slog.Level
	loggerLevels     map[string]slog.Level
	groups           map[string][]string
	consolePattern   string
	filePattern      string
	fileName         string
	filePath         string
	consoleThreshold *slog.Level
	fileThreshold    *slog.Level
}

func readLoggingProperties(environment coreenv.Environment) (loggingProperties, error) {
	properties := loggingProperties{
		loggerLevels: make(map[string]slog.Level),
		groups:       make(map[string][]string),
	}
	if environment == nil {
		return properties, nil
	}
	properties.consolePattern, _ = environment.GetProperty(PropertyConsolePattern)
	properties.filePattern, _ = environment.GetProperty(PropertyFilePattern)
	properties.fileName, _ = environment.GetProperty(PropertyFileName)
	properties.filePath, _ = environment.GetProperty(PropertyFilePath)

	var err error
	if properties.rootLevel, err = readOptionalLevel(environment, PropertyRootLevel); err != nil {
		return loggingProperties{}, err
	}
	if properties.consoleThreshold, err = readOptionalLevel(environment, PropertyConsoleThreshold); err != nil {
		return loggingProperties{}, err
	}
	if properties.fileThreshold, err = readOptionalLevel(environment, PropertyFileThreshold); err != nil {
		return loggingProperties{}, err
	}
	if err := properties.readGroups(environment); err != nil {
		return loggingProperties{}, err
	}
	if err := properties.readLevels(environment); err != nil {
		return loggingProperties{}, err
	}
	return properties, nil
}

func (p *loggingProperties) readGroups(environment coreenv.Environment) error {
	values, found, err := coreenv.GetPropertyMapAsValue[string](environment, strings.TrimSuffix(PropertyGroupPrefix, "."))
	if err != nil {
		return fmt.Errorf("gbc-log: read logging groups: %w", err)
	}
	if !found {
		return nil
	}
	for name, value := range values {
		name = strings.TrimSpace(name)
		members := splitNonEmpty(value)
		if name == "" || len(members) == 0 {
			return fmt.Errorf("gbc-log: logging group %q has no members", name)
		}
		p.groups[name] = members
	}
	return nil
}

func (p *loggingProperties) readLevels(environment coreenv.Environment) error {
	values, found, err := coreenv.GetPropertyMapAsValue[string](environment, strings.TrimSuffix(PropertyLevelPrefix, "."))
	if err != nil {
		return fmt.Errorf("gbc-log: read logger levels: %w", err)
	}
	if !found {
		return nil
	}
	groupLevels := make(map[string]slog.Level)
	directLevels := make(map[string]slog.Level)
	for name, value := range values {
		if name == "root" {
			continue
		}
		level, parseErr := parseLevel(PropertyLevelPrefix+name, value)
		if parseErr != nil {
			return parseErr
		}
		if _, group := p.groups[name]; group {
			groupLevels[name] = level
		} else {
			directLevels[name] = level
		}
	}
	for _, name := range sortedLevelNames(groupLevels) {
		for _, member := range p.groups[name] {
			p.loggerLevels[member] = groupLevels[name]
		}
	}
	for name, level := range directLevels {
		p.loggerLevels[name] = level
	}
	return nil
}

func readOptionalLevel(environment coreenv.Environment, key string) (*slog.Level, error) {
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

func parseLevel(key string, value string) (slog.Level, error) {
	level, err := goarklog.ParseLevel(strings.TrimSpace(value))
	if err != nil {
		return 0, fmt.Errorf("gbc-log: invalid level for %q: %w", key, err)
	}
	return level, nil
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

func sortedLevelNames(levels map[string]slog.Level) []string {
	names := make([]string, 0, len(levels))
	for name := range levels {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
