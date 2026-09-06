package gbclog

import (
	"context"
	"log/slog"
	"strings"

	internaloutput "goark.dev/gbc-log/internal/output"
	internalproperties "goark.dev/gbc-log/internal/properties"
	coreenv "goark.dev/goark/core/env"
	"goark.dev/log"
)

func loggingOptionsCustomizer(
	environment coreenv.Environment,
	customizerGroups ...[]log.StructuredJSONCustomizer,
) log.OptionsCustomizer {
	var customizers []log.StructuredJSONCustomizer
	if len(customizerGroups) > 0 {
		customizers = customizerGroups[0]
	}
	return func(
		_ context.Context,
		current log.Options,
		source *log.ConfigResult,
	) (log.Options, error) {
		properties, err := readLoggingProperties(environment)
		if err != nil {
			return current, err
		}
		if source != nil && source.Source == log.ConfigSourceDefault {
			if current, err = internaloutput.ApplyDefault(current, properties, customizers...); err != nil {
				return current, err
			}
		}
		applyLoggerLevels(&current, properties)
		applyAppenderThreshold(&current, "console", properties.ConsoleThreshold)
		applyAppenderThreshold(&current, "file", properties.FileThreshold)
		return current, nil
	}
}

func applyLoggerLevels(options *log.Options, properties loggingProperties) {
	if properties.RootLevel != nil {
		options.Root.Level = *properties.RootLevel
	}
	indices := make(map[string]int, len(options.Loggers))
	for index := range options.Loggers {
		indices[options.Loggers[index].Name] = index
	}
	for _, name := range internalproperties.SortedLevelNames(properties.LoggerLevels) {
		level := properties.LoggerLevels[name]
		if index, found := indices[name]; found {
			options.Loggers[index].Level = levelPointer(level)
			continue
		}
		options.Loggers = append(
			options.Loggers,
			log.LoggerRule{Name: name, Level: levelPointer(level)},
		)
	}
}

func applyAppenderThreshold(options *log.Options, wanted string, level *slog.Level) {
	if level == nil {
		return
	}
	name := findAppenderName(options.Appenders, wanted)
	if name == "" {
		return
	}
	if len(options.Root.AppenderRefs) == 0 && len(options.Root.AppenderRefControls) == 0 &&
		len(
			options.Appenders,
		) > 0 && options.Appenders[0] != nil && strings.EqualFold(options.Appenders[0].Name(), name) {
		options.Root.AppenderRefControls = append(options.Root.AppenderRefControls,
			log.NewAppenderRef(name, log.WithAppenderRefLevel(*level)))
	}
	applyReferenceThreshold(
		&options.Root.AppenderRefs,
		&options.Root.AppenderRefControls,
		name,
		*level,
	)
	for index := range options.Loggers {
		applyReferenceThreshold(
			&options.Loggers[index].AppenderRefs,
			&options.Loggers[index].AppenderRefControls,
			name,
			*level,
		)
	}
}

func applyReferenceThreshold(
	refs *[]string,
	controls *[]log.AppenderRef,
	name string,
	level slog.Level,
) {
	hadSimpleRef := containsAppenderRef(*refs, name)
	*refs = removeAppenderRef(*refs, name)
	for index := range *controls {
		if strings.EqualFold((*controls)[index].Ref, name) {
			(*controls)[index].Level = levelPointer(level)
			return
		}
	}
	if hadSimpleRef {
		*controls = append(*controls, log.NewAppenderRef(name, log.WithAppenderRefLevel(level)))
	}
}

func findAppenderName(appenders []log.Appender, wanted string) string {
	for _, appender := range appenders {
		if appender != nil && strings.EqualFold(appender.Name(), wanted) {
			return appender.Name()
		}
	}
	return ""
}

func removeAppenderRef(refs []string, name string) []string {
	result := refs[:0]
	for _, ref := range refs {
		if !strings.EqualFold(ref, name) {
			result = append(result, ref)
		}
	}
	return result
}

func containsAppenderRef(refs []string, name string) bool {
	for _, ref := range refs {
		if strings.EqualFold(ref, name) {
			return true
		}
	}
	return false
}

func levelPointer(level slog.Level) *slog.Level {
	copy := level
	return &copy
}
