# Goark Boot Contrib Log

[English](README.md)

`goark.dev/gbc-log` 是 `goark.dev/log` 的 Goark Boot Starter。它创建应用日志运行期，暴露 primary `*slog.Logger`，可选安装为 `slog.Default`，并在观测 exporter 关闭前排空日志。

## 使用方式

```go
app, err := boot.Run(ctx, boot.WithAutoConfiguration(
    gbclog.AutoConfigure(),
))
```

Starter 直接复用 goark-log 的配置发现机制。可以通过 `goark.log.config`、`goark.logging.config`、`logging.config`、`GOARK_LOG_CONFIG` 或标准 `conf/goark-log.*` 路径指定配置。

## Bean

- `goark.log.logger`：primary `*slog.Logger`。
- `goark.log.context`：启用时由容器管理的 `*goarklog.LoggerContext`。
- `goark.log.lifecycle`：负责恢复默认 Logger 和关闭日志运行期。

## 配置属性

- `goark.log.enabled`：默认 `true`。
- `goark.log.install-default`：默认 `true`。
- `goark.log.config`：可选的 goark-log 配置文件路径。

## 生命周期

运行期顺序是 `-9000`。它在普通应用组件之后、顺序为 `-10000` 的观测 Provider 之前停止，保证输出到 goark-log 的 telemetry 在 exporter 关闭前排空。

## 验证

```bash
go test ./...
go test -race ./...
go vet ./...
```

使用 Apache License 2.0。
