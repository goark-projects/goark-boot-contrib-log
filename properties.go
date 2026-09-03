package gbclog

const (
	// PropertyConfig 指定独立 goark-log 配置文件。
	PropertyConfig = "logging.config"
	// PropertyRootLevel 设置根 Logger 级别。
	PropertyRootLevel = "logging.level.root"
	// PropertyLevelPrefix 设置命名 Logger 或日志组级别前缀。
	PropertyLevelPrefix = "logging.level."
	// PropertyGroupPrefix 设置日志组成员前缀。
	PropertyGroupPrefix = "logging.group."
	// PropertyConsolePattern 设置默认控制台输出格式。
	PropertyConsolePattern = "logging.pattern.console"
	// PropertyFilePattern 设置默认文件输出格式。
	PropertyFilePattern = "logging.pattern.file"
	// PropertyFileName 设置默认日志文件名。
	PropertyFileName = "logging.file.name"
	// PropertyFilePath 设置默认日志文件目录。
	PropertyFilePath = "logging.file.path"
	// PropertyConsoleThreshold 设置默认控制台输出阈值。
	PropertyConsoleThreshold = "logging.threshold.console"
	// PropertyFileThreshold 设置默认文件输出阈值。
	PropertyFileThreshold = "logging.threshold.file"
)
