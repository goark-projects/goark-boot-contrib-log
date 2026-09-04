package gbclog

import (
	"context"
	"io/fs"

	coreenv "goark.dev/goark/core/env"
	coreresource "goark.dev/goark/core/resource"
	goarklog "goark.dev/log"
)

// LoggerContextFactory 创建由 Starter 管理的 goark-log 运行期。
type LoggerContextFactory func(ctx context.Context, environment coreenv.Environment) (*goarklog.LoggerContext, error)

type settings struct {
	enabled         *bool
	installDefault  *bool
	factory         LoggerContextFactory
	resourceLoader  coreresource.Loader
	classpathFS     fs.FS
	customizers     []goarklog.StructuredJSONCustomizer
	manageLifecycle bool
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

// WithResourceLoader 设置 logging.config 使用的资源加载器。
func WithResourceLoader(loader coreresource.Loader) Option {
	return func(settings *settings) { settings.resourceLoader = loader }
}

// WithClasspathFS 注册 classpath: 配置资源使用的文件系统。
func WithClasspathFS(filesystem fs.FS) Option {
	return func(settings *settings) { settings.classpathFS = filesystem }
}

// WithStructuredJSONCustomizers 注册显式、类型安全的结构化 JSON 定制器。
func WithStructuredJSONCustomizers(customizers ...goarklog.StructuredJSONCustomizer) Option {
	copied := append([]goarklog.StructuredJSONCustomizer(nil), customizers...)
	return func(settings *settings) { settings.customizers = copied }
}
