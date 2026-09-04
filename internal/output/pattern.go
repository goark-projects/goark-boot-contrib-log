package output

import (
	"strings"

	"goark.dev/gbc-log/internal/properties"
)

func defaultPattern(configuration properties.Properties, file bool) string {
	thread := "%15.15thread"
	if file {
		thread = "%thread"
	}
	var pattern strings.Builder
	pattern.Grow(160)
	pattern.WriteString("%d{")
	pattern.WriteString(configuration.DateFormatPattern)
	pattern.WriteString("}  ")
	pattern.WriteString(configuration.LevelPattern)
	pattern.WriteString(" %pid - ")
	appendApplicationIdentity(&pattern, configuration)
	pattern.WriteByte('[')
	pattern.WriteString(thread)
	pattern.WriteString("] ")
	pattern.WriteString("%-40.40logger{1.2*} : ")
	pattern.WriteString(configuration.CorrelationPattern)
	pattern.WriteString("%msg%attrs%n")
	pattern.WriteString(exceptionPattern(configuration.ExceptionConversionWord))
	return pattern.String()
}

func appendApplicationIdentity(pattern *strings.Builder, configuration properties.Properties) {
	if configuration.IncludeApplicationName && configuration.ApplicationName != "" {
		pattern.WriteByte('[')
		pattern.WriteString(escapePatternLiteral(configuration.ApplicationName))
		pattern.WriteString("] ")
	}
	if configuration.IncludeApplicationGroup && configuration.ApplicationGroup != "" {
		pattern.WriteByte('[')
		pattern.WriteString(escapePatternLiteral(configuration.ApplicationGroup))
		pattern.WriteString("] ")
	}
}

func exceptionPattern(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "%wex", "%ex", "%throwable", "%exception":
		return "%ex"
	case "%nopex":
		return ""
	default:
		return value
	}
}

func escapePatternLiteral(value string) string {
	return strings.ReplaceAll(value, "%", "%%")
}
