package output

import (
	"fmt"
	"strings"

	"goark.dev/gbc-log/internal/properties"
	goarklog "goark.dev/log"
)

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
