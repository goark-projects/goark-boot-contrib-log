package output

import (
	"io"
	"path/filepath"
	"strings"

	"goark.dev/gbc-log/internal/properties"
	goarklog "goark.dev/log"
)

// ApplyDefault 使用 Boot 属性替换 goark-log 的默认输出端。
func ApplyDefault(options goarklog.Options, configuration properties.Properties, customizers ...goarklog.StructuredJSONCustomizer) (goarklog.Options, error) {
	appenders := make([]goarklog.Appender, 0, 2)
	refs := make([]string, 0, 2)
	if configuration.ConsoleEnabled {
		layout, err := outputLayout(configuration, false, customizers)
		if err != nil {
			return options, err
		}
		appenders = append(appenders, goarklog.NewConsoleAppender(goarklog.WithConsoleLayout(layout)))
		refs = append(refs, "console")
	}
	fileName := resolvedFileName(configuration)
	if fileName != "" {
		appender, err := newRollingFileAppender(fileName, configuration, customizers)
		if err != nil {
			return options, err
		}
		appenders = append(appenders, appender)
		refs = append(refs, "file")
	}
	if len(appenders) == 0 {
		appenders = append(appenders, goarklog.NewConsoleAppender(
			goarklog.WithConsoleName("discard"),
			goarklog.WithConsoleWriter(io.Discard),
		))
		refs = append(refs, "discard")
	}
	options.Appenders = appenders
	options.Root.AppenderRefs = refs
	options.Root.AppenderRefControls = nil
	return options, nil
}

func resolvedFileName(configuration properties.Properties) string {
	if name := strings.TrimSpace(configuration.FileName); name != "" {
		return name
	}
	if path := strings.TrimSpace(configuration.FilePath); path != "" {
		return filepath.Join(path, "goark.log")
	}
	return ""
}

func newRollingFileAppender(fileName string, configuration properties.Properties, customizers []goarklog.StructuredJSONCustomizer) (*goarklog.RollingFileAppender, error) {
	layout, err := outputLayout(configuration, true, customizers)
	if err != nil {
		return nil, err
	}
	filePattern := strings.TrimSpace(configuration.RollingFileNamePattern)
	if filePattern == "" {
		filePattern = fileName + ".%d{yyyy-MM-dd}.%i.gz"
	}
	filePattern = strings.ReplaceAll(filePattern, "${LOG_FILE}", fileName)
	return goarklog.NewRollingFileAppender(fileName,
		goarklog.WithRollingFileName("file"),
		goarklog.WithRollingFileLayout(layout),
		goarklog.WithRollingFilePattern(filePattern),
		goarklog.WithRollingMaxSize(configuration.MaxFileSize),
		goarklog.WithRollingMaxBackups(configuration.MaxHistory),
		goarklog.WithRollingTotalSizeCap(configuration.TotalSizeCap),
		goarklog.WithRollingCleanHistoryOnStart(configuration.CleanHistoryOnStart),
	)
}
