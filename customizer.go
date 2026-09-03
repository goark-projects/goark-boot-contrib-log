package gbclog

import (
	"context"
	"log/slog"
	"path/filepath"
	"strings"

	coreenv "goark.dev/goark/core/env"
	goarklog "goark.dev/log"
)

func loggingOptionsCustomizer(environment coreenv.Environment) goarklog.OptionsCustomizer {
	return func(_ context.Context, current goarklog.Options, source *goarklog.ConfigResult) (goarklog.Options, error) {
		properties, err := readLoggingProperties(environment)
		if err != nil {
			return current, err
		}
		if source != nil && source.Source == goarklog.ConfigSourceDefault {
			if current, err = applyDefaultOutputs(current, properties); err != nil {
				return current, err
			}
		}
		applyLoggerLevels(&current, properties)
		applyAppenderThreshold(&current, "console", properties.consoleThreshold)
		applyAppenderThreshold(&current, "file", properties.fileThreshold)
		return current, nil
	}
}

func applyDefaultOutputs(options goarklog.Options, properties loggingProperties) (goarklog.Options, error) {
	consoleOptions := make([]goarklog.ConsoleOption, 0, 1)
	if strings.TrimSpace(properties.consolePattern) != "" {
		layout, err := goarklog.NewPatternLayout(properties.consolePattern)
		if err != nil {
			return options, err
		}
		consoleOptions = append(consoleOptions, goarklog.WithConsoleLayout(layout))
	}
	appenders := []goarklog.Appender{goarklog.NewConsoleAppender(consoleOptions...)}
	refs := []string{"console"}

	fileName := strings.TrimSpace(properties.fileName)
	if fileName == "" && strings.TrimSpace(properties.filePath) != "" {
		fileName = filepath.Join(strings.TrimSpace(properties.filePath), "goark.log")
	}
	if fileName != "" {
		fileOptions := []goarklog.FileOption{goarklog.WithFileName("file")}
		if strings.TrimSpace(properties.filePattern) != "" {
			layout, err := goarklog.NewPatternLayout(properties.filePattern)
			if err != nil {
				return options, err
			}
			fileOptions = append(fileOptions, goarklog.WithFileLayout(layout))
		}
		appender, err := goarklog.NewFileAppender(fileName, fileOptions...)
		if err != nil {
			return options, err
		}
		appenders = append(appenders, appender)
		refs = append(refs, "file")
	}
	options.Appenders = appenders
	options.Root.AppenderRefs = refs
	options.Root.AppenderRefControls = nil
	return options, nil
}

func applyLoggerLevels(options *goarklog.Options, properties loggingProperties) {
	if properties.rootLevel != nil {
		options.Root.Level = *properties.rootLevel
	}
	indices := make(map[string]int, len(options.Loggers))
	for index := range options.Loggers {
		indices[options.Loggers[index].Name] = index
	}
	for _, name := range sortedLevelNames(properties.loggerLevels) {
		level := properties.loggerLevels[name]
		if index, found := indices[name]; found {
			options.Loggers[index].Level = levelPointer(level)
			continue
		}
		options.Loggers = append(options.Loggers, goarklog.LoggerRule{Name: name, Level: levelPointer(level)})
	}
}

func applyAppenderThreshold(options *goarklog.Options, wanted string, level *slog.Level) {
	if level == nil {
		return
	}
	name := findAppenderName(options.Appenders, wanted)
	if name == "" {
		return
	}
	applyReferenceThreshold(&options.Root.AppenderRefs, &options.Root.AppenderRefControls, name, *level)
	for index := range options.Loggers {
		applyReferenceThreshold(&options.Loggers[index].AppenderRefs, &options.Loggers[index].AppenderRefControls, name, *level)
	}
}

func applyReferenceThreshold(refs *[]string, controls *[]goarklog.AppenderRef, name string, level slog.Level) {
	hadSimpleRef := containsAppenderRef(*refs, name)
	*refs = removeAppenderRef(*refs, name)
	for index := range *controls {
		if strings.EqualFold((*controls)[index].Ref, name) {
			(*controls)[index].Level = levelPointer(level)
			return
		}
	}
	if hadSimpleRef {
		*controls = append(*controls, goarklog.NewAppenderRef(name, goarklog.WithAppenderRefLevel(level)))
	}
}

func findAppenderName(appenders []goarklog.Appender, wanted string) string {
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
