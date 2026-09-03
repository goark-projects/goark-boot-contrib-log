package gbclog

import (
	"context"

	coreenv "goark.dev/goark/core/env"
	goarklog "goark.dev/log"
)

// LoggerContextFactory 创建由 Starter 管理的 goark-log 运行期。
type LoggerContextFactory func(ctx context.Context, environment coreenv.Environment) (*goarklog.LoggerContext, error)

type settings struct {
	enabled        *bool
	installDefault *bool
	factory        LoggerContextFactory
}

// Option 调整日志自动配置。
type Option func(*settings)

// WithEnabled 控制是否启用 goark-log。
func WithEnabled(enabled bool) Option {
	return func(settings *settings) { settings.enabled = &enabled }
}

// WithInstallDefault 控制是否设置进程级 slog 默认 Logger。
func WithInstallDefault(enabled bool) Option {
	return func(settings *settings) { settings.installDefault = &enabled }
}

// WithLoggerContextFactory 替换日志运行期工厂，主要用于自定义 Appender 或测试。
func WithLoggerContextFactory(factory LoggerContextFactory) Option {
	return func(settings *settings) { settings.factory = factory }
}
