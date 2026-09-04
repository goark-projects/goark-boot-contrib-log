package gbclog

const (
	// ModulePath 是该官方 contrib 模块的 Go 导入路径。
	ModulePath = "goark.dev/gbc-log"
	// Repository 是该模块对应的官方 Git 仓库名。
	Repository = "goark-boot-contrib-log"
	// StarterID 是后续自动配置注册时使用的稳定启动器标识。
	StarterID = "goark.boot.log"
	// BeanNameContext 是 goark-log 运行期 Bean 名称。
	BeanNameContext = "goark.log.context"
	// BeanNameLogger 是应用级 slog Logger Bean 名称。
	BeanNameLogger = "goark.log.logger"
	// BeanNameLifecycle 是日志运行期生命周期 Bean 名称。
	BeanNameLifecycle = "goark.log.lifecycle"
	// BeanNameSystem 是日志运行时控制面 Bean 名称。
	BeanNameSystem = "goark.log.system"
)

const (
	// PropertyEnabled 控制日志 Starter 是否创建 goark-log 运行期。
	PropertyEnabled = "goark.log.enabled"
	// PropertyInstallDefault 控制是否安装进程级 slog 默认 Logger。
	PropertyInstallDefault = "goark.log.install-default"
)

const (
	// DefaultEnabled 表示默认启用 goark-log 运行期。
	DefaultEnabled = true
	// DefaultInstallDefault 表示默认安装进程级 slog Logger。
	DefaultInstallDefault = true
)
