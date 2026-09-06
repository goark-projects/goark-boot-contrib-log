package gbclog

import (
	"fmt"
	"log/slog"
	"strings"

	"goark.dev/log"
)

// LoggingSystem 提供运行期 Logger 级别控制和配置查询。
type LoggingSystem interface {
	SetLogLevel(logger string, level *slog.Level) error
	LogLevel(logger string) (log.LoggerConfiguration, bool)
	Loggers() []log.LoggerConfiguration
}

// SetLogLevel 原子修改 Logger 级别；nil 恢复静态配置或继承关系。
func (r *Runtime) SetLogLevel(logger string, level *slog.Level) error {
	if r == nil || r.context == nil {
		return fmt.Errorf("gbc-log: logging system is disabled")
	}
	return r.context.SetLevel(logger, level)
}

// LogLevel 返回指定 Logger 的配置快照。
func (r *Runtime) LogLevel(logger string) (log.LoggerConfiguration, bool) {
	logger = strings.TrimSpace(logger)
	if strings.EqualFold(logger, "root") {
		logger = "ROOT"
	}
	for _, configuration := range r.Loggers() {
		if configuration.Name == logger {
			return configuration, true
		}
	}
	return log.LoggerConfiguration{}, false
}

// Loggers 返回 Root 和所有已配置 Logger 的稳定快照。
func (r *Runtime) Loggers() []log.LoggerConfiguration {
	if r == nil || r.context == nil {
		return nil
	}
	return r.context.LoggerConfigurations()
}

var _ LoggingSystem = (*Runtime)(nil)
