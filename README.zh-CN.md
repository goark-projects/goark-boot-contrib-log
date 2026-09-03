# Goark Boot Contrib Log

[English](README.md)

`goark.dev/gbc-log` 负责 Goark Boot 与 `goark.dev/log` 的集成。它在基础设施
Starter 之前创建日志运行期，暴露 primary `*slog.Logger`，默认安装为
`slog.Default`，并在观测 Exporter 停止后排空日志。

## 职责边界

- `goark.dev/log` 负责 Logger API、Handler、Appender、Layout、Filter、异步
  投递、独立配置文件及重载。
- `goark.dev/gbc-log` 读取 Boot `Environment`，解释 `logging.*`，注册 Bean，
  安装进程默认 Logger，并管理生命周期顺序。
- 其他库依赖 `log/slog` 或 `goark.dev/log`，不依赖本 Starter。

## 使用方式

```go
app, err := boot.Run(ctx, boot.WithAutoConfiguration(
	gbclog.AutoConfigure(),
))
```

## 配置属性

| 属性 | 默认值 | 说明 |
| --- | --- | --- |
| `goark.log.enabled` | `true` | 是否启用受管 goark-log 运行期 |
| `goark.log.install-default` | `true` | 是否安装为 `slog.Default` |
| `logging.config` | 自动发现 | 独立 goark-log 配置文件的精确路径 |
| `logging.level.root` | `INFO` | 根 Logger 级别 |
| `logging.level.<logger>` | 继承 | 命名 Logger 级别 |
| `logging.group.<name>` | 未设置 | 日志组包含的 Logger 名称，逗号分隔 |
| `logging.level.<group>` | 未设置 | 应用于组内全部 Logger 的级别 |
| `logging.pattern.console` | Spring Boot 风格 | 默认控制台 PatternLayout 格式 |
| `logging.pattern.file` | Spring Boot 风格 | 默认文件 PatternLayout 格式 |
| `logging.file.name` | 未设置 | 默认日志文件的精确路径 |
| `logging.file.path` | 未设置 | 默认 `goark.log` 文件所在目录 |
| `logging.threshold.console` | 未设置 | 控制台 Appender 阈值 |
| `logging.threshold.file` | 未设置 | 文件 Appender 阈值 |

`logging.file.name` 优先于 `logging.file.path`。日志级别支持 goark-log 注册表中的
全部名称，包括 `ALL`、`TRACE`、`DEBUG`、`INFO`、`WARN`、`ERROR`、`FATAL`
和 `OFF`。直接 Logger 配置优先于日志组展开结果。

## 独立日志配置

`logging.config` 可指定 goark-log 原生 YAML、TOML、JSON、XML 或 properties
配置文件。使用独立配置文件时：

- Appender、Layout、Filter、additivity、异步参数和重载策略由独立文件负责。
- `logging.level.*` 覆盖根 Logger 与命名 Logger 级别。
- `logging.threshold.console` 和 `logging.threshold.file` 覆盖同名 Appender 引用
  的阈值，同时保留 Filter 与调用位置设置。
- `logging.pattern.*` 和 `logging.file.*` 只定制 Starter 默认输出，不替换独立
  配置中的 Appender 与 Layout。

## Bean 与生命周期

- `goark.log.logger`：primary `*slog.Logger`。
- `goark.log.context`：启用时由容器管理的 `*goarklog.LoggerContext`。
- `goark.log.lifecycle`：负责恢复默认 Logger 和关闭日志运行期。

运行期顺序为 `-11000`。它在顺序为 `-10000` 的观测 Provider 之前启动、之后
停止，保证 Exporter 刷新和关闭产生的最终记录仍可写入，随后再排空 goark-log。

## 验证

```bash
go test ./...
go test -race ./...
go vet ./...
```

使用 Apache License 2.0。
