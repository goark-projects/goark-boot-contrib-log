package properties

import (
	"fmt"
	"log/slog"
	"strings"

	coreenv "goark.dev/goark/core/env"
	"goark.dev/log"
)

// Read 从 Environment 编译日志配置。
func Read(environment coreenv.Environment) (Properties, error) {
	properties := defaults()
	if environment == nil {
		return properties, nil
	}
	properties.Config = optionalString(environment, Config)
	properties.ConsoleCharset = optionalString(environment, ConsoleCharset)
	properties.FileCharset = optionalString(environment, FileCharset)
	properties.ConsolePattern = optionalString(environment, ConsolePattern)
	properties.FilePattern = optionalString(environment, FilePattern)
	properties.DateFormatPattern = stringWithDefault(
		environment,
		DateFormatPattern,
		DefaultDateFormat,
	)
	properties.LevelPattern = stringWithDefault(environment, LevelPattern, DefaultLevelPattern)
	properties.CorrelationPattern = optionalString(environment, CorrelationPattern)
	properties.ExceptionConversionWord = stringWithDefault(
		environment,
		ExceptionConversionWord,
		DefaultExceptionConversionWord,
	)
	properties.FileName = firstNonEmpty(environment, FileName, LegacyFileName)
	properties.FilePath = firstNonEmpty(environment, FilePath, LegacyFilePath)
	properties.RollingFileNamePattern = optionalString(environment, RollingFileNamePattern)
	properties.ApplicationName = optionalString(environment, ApplicationName)
	properties.ApplicationGroup = optionalString(environment, ApplicationGroup)
	properties.Structured.ConsoleFormat = optionalString(environment, StructuredConsoleFormat)
	properties.Structured.FileFormat = optionalString(environment, StructuredFileFormat)
	properties.Structured.JSON.Include = splitNonEmpty(
		optionalString(environment, StructuredJSONInclude),
	)
	properties.Structured.JSON.Exclude = splitNonEmpty(
		optionalString(environment, StructuredJSONExclude),
	)
	properties.Structured.JSON.ContextPrefix = optionalString(environment, StructuredContextPrefix)
	properties.Structured.JSON.Customizers = splitNonEmpty(
		optionalString(environment, StructuredJSONCustomizer),
	)
	properties.Structured.JSON.Stacktrace.Printer = optionalString(environment, StacktracePrinter)
	properties.Structured.JSON.Stacktrace.Root = optionalString(environment, StacktraceRoot)
	properties.Structured.ECS.ServiceEnvironment = optionalString(
		environment,
		ECSServiceEnvironment,
	)
	properties.Structured.ECS.ServiceName = optionalString(environment, ECSServiceName)
	properties.Structured.ECS.ServiceNodeName = optionalString(environment, ECSServiceNodeName)
	properties.Structured.ECS.ServiceVersion = optionalString(environment, ECSServiceVersion)
	properties.Structured.GELF.Host = optionalString(environment, GELFHost)
	properties.Structured.GELF.ServiceVersion = optionalString(environment, GELFServiceVersion)

	var err error
	properties.ConsoleEnabled, err = boolWithDefault(environment, ConsoleEnabled, true)
	if err != nil {
		return Properties{}, err
	}
	properties.IncludeApplicationName, err = boolWithDefault(
		environment, IncludeApplicationName, true,
	)
	if err != nil {
		return Properties{}, err
	}
	properties.IncludeApplicationGroup, err = boolWithDefault(
		environment, IncludeApplicationGroup, true,
	)
	if err != nil {
		return Properties{}, err
	}
	properties.RegisterShutdownHook, err = boolWithDefault(
		environment, RegisterShutdownHook, true,
	)
	if err != nil {
		return Properties{}, err
	}
	properties.CleanHistoryOnStart, err = boolWithDefault(
		environment, FileCleanHistoryOnStart, false,
	)
	if err != nil {
		return Properties{}, err
	}
	properties.MaxHistory, err = intWithDefault(
		environment, FileMaxHistory, DefaultMaxHistory,
	)
	if err != nil {
		return Properties{}, err
	}
	properties.MaxFileSize, err = dataSizeWithDefault(
		environment, FileMaxSize, DefaultMaxFileSize,
	)
	if err != nil {
		return Properties{}, err
	}
	properties.TotalSizeCap, err = dataSizeWithDefault(environment, FileTotalSizeCap, 0)
	if err != nil {
		return Properties{}, err
	}
	if properties.RootLevel, err = optionalLevel(environment, RootLevel); err != nil {
		return Properties{}, err
	}
	properties.ConsoleThreshold, err = levelWithDefault(
		environment, ConsoleThreshold, log.LevelTrace,
	)
	if err != nil {
		return Properties{}, err
	}
	properties.FileThreshold, err = levelWithDefault(
		environment, FileThreshold, log.LevelTrace,
	)
	if err != nil {
		return Properties{}, err
	}
	properties.Structured.JSON.ContextInclude, err = optionalBool(
		environment, StructuredContextInclude,
	)
	if err != nil {
		return Properties{}, err
	}
	properties.Structured.JSON.Stacktrace.IncludeCommonFrames, err = optionalBool(
		environment, StacktraceCommonFrames,
	)
	if err != nil {
		return Properties{}, err
	}
	properties.Structured.JSON.Stacktrace.IncludeHashes, err = optionalBool(
		environment, StacktraceHashes,
	)
	if err != nil {
		return Properties{}, err
	}
	properties.Structured.JSON.Stacktrace.MaxLength, err = optionalInt(
		environment, StacktraceMaxLength,
	)
	if err != nil {
		return Properties{}, err
	}
	properties.Structured.JSON.Stacktrace.MaxThrowableDepth, err = optionalInt(
		environment, StacktraceMaxDepth,
	)
	if err != nil {
		return Properties{}, err
	}
	if properties.Groups, err = readGroups(environment); err != nil {
		return Properties{}, err
	}
	if properties.LoggerLevels, err = readLevels(environment, properties.Groups); err != nil {
		return Properties{}, err
	}
	properties.Structured.JSON.Rename, err = readStringMap(environment, StructuredJSONRename)
	if err != nil {
		return Properties{}, err
	}
	properties.Structured.JSON.Add, err = readStringMap(environment, StructuredJSONAdd)
	if err != nil {
		return Properties{}, err
	}
	if err := validate(properties); err != nil {
		return Properties{}, err
	}
	return properties, nil
}

func defaults() Properties {
	return Properties{
		ConsoleEnabled:          true,
		DateFormatPattern:       DefaultDateFormat,
		LevelPattern:            DefaultLevelPattern,
		ExceptionConversionWord: DefaultExceptionConversionWord,
		MaxHistory:              DefaultMaxHistory,
		MaxFileSize:             DefaultMaxFileSize,
		IncludeApplicationName:  true,
		IncludeApplicationGroup: true,
		RegisterShutdownHook:    true,
		ConsoleThreshold:        levelPointer(log.LevelTrace),
		FileThreshold:           levelPointer(log.LevelTrace),
		LoggerLevels:            make(map[string]slog.Level),
		Groups:                  cloneGroups(builtInGroups),
		Structured: StructuredProperties{JSON: JSONProperties{
			Rename: make(map[string]string),
			Add:    make(map[string]string),
		}},
	}
}

func firstNonEmpty(environment coreenv.Environment, keys ...string) string {
	for _, key := range keys {
		if value := optionalString(environment, key); value != "" {
			return value
		}
	}
	return ""
}

func validate(properties Properties) error {
	if len(properties.Structured.JSON.Customizers) > 0 {
		return fmt.Errorf(
			"gbc-log: %q contains Java class names; "+
				"register Go customizers with WithStructuredJSONCustomizers",
			StructuredJSONCustomizer,
		)
	}
	for key, format := range map[string]string{
		StructuredConsoleFormat: properties.Structured.ConsoleFormat,
		StructuredFileFormat:    properties.Structured.FileFormat,
	} {
		switch strings.ToLower(format) {
		case "", "ecs", "gelf", "logstash":
		default:
			return fmt.Errorf(
				"gbc-log: unsupported structured logging format %q for %q",
				format,
				key,
			)
		}
	}
	root := strings.ToLower(properties.Structured.JSON.Stacktrace.Root)
	if root != "" && root != "first" && root != "last" {
		return fmt.Errorf("gbc-log: %q must be first or last", StacktraceRoot)
	}
	printer := strings.ToLower(properties.Structured.JSON.Stacktrace.Printer)
	if printer != "" && printer != "standard" && printer != "logging-system" {
		return fmt.Errorf("gbc-log: %q must be standard or logging-system", StacktracePrinter)
	}
	return nil
}
