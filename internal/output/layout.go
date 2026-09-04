package output

import (
	"fmt"
	"strings"

	"goark.dev/gbc-log/internal/properties"
	goarklog "goark.dev/log"
)

func outputLayout(configuration properties.Properties, file bool, customizers []goarklog.StructuredJSONCustomizer) (goarklog.Layout, error) {
	format := configuration.Structured.ConsoleFormat
	if file {
		format = configuration.Structured.FileFormat
	}
	var layout goarklog.Layout
	var err error
	if strings.TrimSpace(format) == "" {
		layout, err = textLayout(configuration, file)
	} else {
		layout, err = structuredLayout(configuration, format, customizers)
	}
	if err != nil {
		return nil, err
	}
	charset := configuration.ConsoleCharset
	if file {
		charset = configuration.FileCharset
	}
	encoded, err := goarklog.NewCharsetLayout(layout, charset)
	if err != nil {
		return nil, fmt.Errorf("gbc-log: configure logging charset: %w", err)
	}
	return encoded, nil
}

func textLayout(configuration properties.Properties, file bool) (goarklog.Layout, error) {
	pattern := configuration.ConsolePattern
	if file {
		pattern = configuration.FilePattern
	}
	if strings.TrimSpace(pattern) == "" {
		pattern = defaultPattern(configuration, file)
	}
	layout, err := goarklog.NewPatternLayout(pattern)
	if err != nil {
		return nil, fmt.Errorf("gbc-log: compile logging pattern: %w", err)
	}
	return layout, nil
}

func structuredLayout(configuration properties.Properties, format string, customizers []goarklog.StructuredJSONCustomizer) (goarklog.Layout, error) {
	includeContext := true
	if configured := configuration.Structured.JSON.ContextInclude; configured != nil {
		includeContext = *configured
	}
	stacktrace := configuration.Structured.JSON.Stacktrace
	options := goarklog.StructuredJSONOptions{
		Format:         goarklog.StructuredFormat(strings.ToLower(strings.TrimSpace(format))),
		Include:        configuration.Structured.JSON.Include,
		Exclude:        configuration.Structured.JSON.Exclude,
		Rename:         configuration.Structured.JSON.Rename,
		Add:            configuration.Structured.JSON.Add,
		IncludeContext: includeContext,
		ContextPrefix:  configuration.Structured.JSON.ContextPrefix,
		Stacktrace: goarklog.StructuredStacktraceOptions{
			RootFirst:           strings.EqualFold(stacktrace.Root, "first"),
			MaxLength:           optionalIntValue(stacktrace.MaxLength),
			MaxThrowableDepth:   optionalIntValue(stacktrace.MaxThrowableDepth),
			IncludeCommonFrames: optionalBoolValue(stacktrace.IncludeCommonFrames),
			IncludeHashes:       optionalBoolValue(stacktrace.IncludeHashes),
		},
		ECS: goarklog.StructuredECSOptions{
			ServiceEnvironment: configuration.Structured.ECS.ServiceEnvironment,
			ServiceName:        firstText(configuration.Structured.ECS.ServiceName, configuration.ApplicationName),
			ServiceNodeName:    configuration.Structured.ECS.ServiceNodeName,
			ServiceVersion:     configuration.Structured.ECS.ServiceVersion,
		},
		GELF: goarklog.StructuredGELFOptions{
			Host:           firstText(configuration.Structured.GELF.Host, configuration.ApplicationName),
			ServiceName:    firstText(configuration.Structured.GELF.ServiceName, configuration.ApplicationName),
			ServiceVersion: configuration.Structured.GELF.ServiceVersion,
		},
		Customizers: customizers,
	}
	layout, err := goarklog.NewStructuredJSONLayout(options)
	if err != nil {
		return nil, fmt.Errorf("gbc-log: configure structured logging: %w", err)
	}
	return layout, nil
}

func optionalBoolValue(value *bool) bool {
	return value != nil && *value
}

func optionalIntValue(value *int) int {
	if value == nil {
		return 0
	}
	return *value
}

func firstText(values ...string) string {
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			return value
		}
	}
	return ""
}
