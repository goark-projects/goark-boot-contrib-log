package properties

import "log/slog"

// Properties 是启动阶段编译完成的日志配置快照。
type Properties struct {
	Config                  string
	ConsoleEnabled          bool
	ConsoleCharset          string
	FileCharset             string
	RootLevel               *slog.Level
	LoggerLevels            map[string]slog.Level
	Groups                  map[string][]string
	ConsolePattern          string
	FilePattern             string
	DateFormatPattern       string
	LevelPattern            string
	CorrelationPattern      string
	ExceptionConversionWord string
	FileName                string
	FilePath                string
	RollingFileNamePattern  string
	CleanHistoryOnStart     bool
	MaxHistory              int
	MaxFileSize             int64
	TotalSizeCap            int64
	ConsoleThreshold        *slog.Level
	FileThreshold           *slog.Level
	IncludeApplicationName  bool
	IncludeApplicationGroup bool
	RegisterShutdownHook    bool
	ApplicationName         string
	ApplicationGroup        string
	Structured              StructuredProperties
}

// StructuredProperties 描述结构化日志格式及成员变换。
type StructuredProperties struct {
	ConsoleFormat string
	FileFormat    string
	JSON          JSONProperties
	ECS           ECSProperties
	GELF          GELFProperties
}

// JSONProperties 描述结构化 JSON 成员和堆栈配置。
type JSONProperties struct {
	Include        []string
	Exclude        []string
	Rename         map[string]string
	Add            map[string]string
	ContextInclude *bool
	ContextPrefix  string
	Stacktrace     StacktraceProperties
}

// StacktraceProperties 描述结构化异常输出限制。
type StacktraceProperties struct {
	IncludeCommonFrames *bool
	IncludeHashes       *bool
	MaxLength           *int
	MaxThrowableDepth   *int
	Printer             string
	Root                string
}

// ECSProperties 描述 ECS 服务字段。
type ECSProperties struct {
	ServiceEnvironment string
	ServiceName        string
	ServiceNodeName    string
	ServiceVersion     string
}

// GELFProperties 描述 GELF 服务字段。
type GELFProperties struct {
	Host           string
	ServiceVersion string
}
