package gbclog

import internalproperties "goark.dev/gbc-log/internal/properties"

const (
	// PropertyConfig 指定独立 goark-log 配置资源。
	PropertyConfig = internalproperties.Config
	// PropertyRootLevel 设置根 Logger 级别。
	PropertyRootLevel = internalproperties.RootLevel
	// PropertyLevelPrefix 设置命名 Logger 或日志组级别前缀。
	PropertyLevelPrefix = internalproperties.LevelPrefix
	// PropertyGroupPrefix 设置日志组成员前缀。
	PropertyGroupPrefix = internalproperties.GroupPrefix
	// PropertyConsolePattern 设置默认控制台输出格式。
	PropertyConsolePattern = internalproperties.ConsolePattern
	// PropertyFilePattern 设置默认文件输出格式。
	PropertyFilePattern = internalproperties.FilePattern
	// PropertyFileName 设置默认日志文件名。
	PropertyFileName = internalproperties.FileName
	// PropertyFilePath 设置默认日志文件目录。
	PropertyFilePath = internalproperties.FilePath
	// PropertyConsoleThreshold 设置默认控制台输出阈值。
	PropertyConsoleThreshold = internalproperties.ConsoleThreshold
	// PropertyFileThreshold 设置默认文件输出阈值。
	PropertyFileThreshold = internalproperties.FileThreshold
)
