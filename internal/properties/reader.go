package properties

import (
	"fmt"
	"log/slog"
	"strings"

	coreenv "goark.dev/goark/core/env"
	goarklog "goark.dev/log"
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
	properties.DateFormatPattern = stringWithDefault(environment, DateFormatPattern, DefaultDateFormat)
	properties.LevelPattern = stringWithDefault(environment, LevelPattern, DefaultLevelPattern)
	properties.CorrelationPattern = optionalString(environment, CorrelationPattern)
	properties.ExceptionConversionWord = stringWithDefault(environment, ExceptionConversionWord, DefaultExceptionConversionWord)
	properties.FileName = firstNonEmpty(environment, FileName, LegacyFileName)
	properties.FilePath = firstNonEmpty(environment, FilePath, LegacyFilePath)
	properties.RollingFileNamePattern = optionalString(environment, RollingFileNamePattern)
	properties.ApplicationName = optionalString(environment, ApplicationName)
	properties.ApplicationGroup = optionalString(environment, ApplicationGroup)
	properties.Structured.ConsoleFormat = optionalString(environment, StructuredConsoleFormat)
	properties.Structured.FileFormat = optionalString(environment, StructuredFileFormat)
	properties.Structured.JSON.Include = splitNonEmpty(optionalString(environment, StructuredJSONInclude))
	properties.Structured.JSON.Exclude = splitNonEmpty(optionalString(environment, StructuredJSONExclude))
	properties.Structured.JSON.ContextPrefix = optionalString(environment, StructuredContextPrefix)
	properties.Structured.JSON.Customizers = splitNonEmpty(optionalString(environment, StructuredJSONCustomizer))
	properties.Structured.JSON.Stacktrace.Printer = optionalString(environment, StacktracePrinter)
	properties.Structured.JSON.Stacktrace.Root = optionalString(environment, StacktraceRoot)
	properties.Structured.ECS.ServiceEnvironment = optionalString(environment, ECSServiceEnvironment)
	properties.Structured.ECS.ServiceName = optionalString(environment, ECSServiceName)
	properties.Structured.ECS.ServiceNodeName = optionalString(environment, ECSServiceNodeName)
	properties.Structured.ECS.ServiceVersion = optionalString(environment, ECSServiceVersion)
	properties.Structured.GELF.Host = optionalString(environment, GELFHost)
	properties.Structured.GELF.ServiceName = optionalString(environment, GELFServiceName)
	properties.Structured.GELF.ServiceVersion = optionalString(environment, GELFServiceVersion)

	var err error
	if properties.ConsoleEnabled, err = boolWithDefault(environment, ConsoleEnabled, true); err != nil {
		return Properties{}, err
	}
	if properties.IncludeApplicationName, err = boolWithDefault(environment, IncludeApplicationName, true); err != nil {
		return Properties{}, err
	}
	if properties.IncludeApplicationGroup, err = boolWithDefault(environment, IncludeApplicationGroup, true); err != nil {
		return Properties{}, err
	}
	if properties.RegisterShutdownHook, err = boolWithDefault(environment, RegisterShutdownHook, true); err != nil {
		return Properties{}, err
	}
	if properties.CleanHistoryOnStart, err = boolWithDefault(environment, FileCleanHistoryOnStart, false); err != nil {
		return Properties{}, err
	}
	if properties.MaxHistory, err = intWithDefault(environment, FileMaxHistory, DefaultMaxHistory); err != nil {
		return Properties{}, err
	}
	if properties.MaxFileSize, err = dataSizeWithDefault(environment, FileMaxSize, DefaultMaxFileSize); err != nil {
		return Properties{}, err
	}
	if properties.TotalSizeCap, err = dataSizeWithDefault(environment, FileTotalSizeCap, 0); err != nil {
		return Properties{}, err
	}
	if properties.RootLevel, err = optionalLevel(environment, RootLevel); err != nil {
		return Properties{}, err
	}
	if properties.ConsoleThreshold, err = levelWithDefault(environment, ConsoleThreshold, goarklog.LevelTrace); err != nil {
		return Properties{}, err
	}
	if properties.FileThreshold, err = levelWithDefault(environment, FileThreshold, goarklog.LevelTrace); err != nil {
		return Properties{}, err
	}
	if properties.Structured.JSON.ContextInclude, err = optionalBool(environment, StructuredContextInclude); err != nil {
		return Properties{}, err
	}
	if properties.Structured.JSON.Stacktrace.IncludeCommonFrames, err = optionalBool(environment, StacktraceCommonFrames); err != nil {
		return Properties{}, err
	}
	if properties.Structured.JSON.Stacktrace.IncludeHashes, err = optionalBool(environment, StacktraceHashes); err != nil {
		return Properties{}, err
	}
	if properties.Structured.JSON.Stacktrace.MaxLength, err = optionalInt(environment, StacktraceMaxLength); err != nil {
		return Properties{}, err
	}
	if properties.Structured.JSON.Stacktrace.MaxThrowableDepth, err = optionalInt(environment, StacktraceMaxDepth); err != nil {
		return Properties{}, err
	}
	if properties.Groups, err = readGroups(environment); err != nil {
		return Properties{}, err
	}
	if properties.LoggerLevels, err = readLevels(environment, properties.Groups); err != nil {
		return Properties{}, err
	}
	if properties.Structured.JSON.Rename, err = readStringMap(environment, StructuredJSONRename); err != nil {
		return Properties{}, err
	}
	if properties.Structured.JSON.Add, err = readStringMap(environment, StructuredJSONAdd); err != nil {
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
		ConsoleThreshold:        levelPointer(goarklog.LevelTrace),
		FileThreshold:           levelPointer(goarklog.LevelTrace),
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
		return fmt.Errorf("gbc-log: %q contains Java class names; register Go customizers with WithStructuredJSONCustomizers", StructuredJSONCustomizer)
	}
	for key, format := range map[string]string{
		StructuredConsoleFormat: properties.Structured.ConsoleFormat,
		StructuredFileFormat:    properties.Structured.FileFormat,
	} {
		switch strings.ToLower(format) {
		case "", "ecs", "gelf", "logstash":
		default:
			return fmt.Errorf("gbc-log: unsupported structured logging format %q for %q", format, key)
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
