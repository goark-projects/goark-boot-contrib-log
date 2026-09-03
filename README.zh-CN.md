# Goark Boot Contrib Log

[English](README.md)

`goark.dev/gbc-log` 是 `goark.dev/log` 的 Goark Boot Starter。它创建应用日志运行期，暴露 primary `*slog.Logger`，可选安装为 `slog.Default`，并保持日志可用直至观测 exporter 关闭。

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

运行期顺序是 `-11000`。它在顺序为 `-10000` 的观测 Provider 之前启动、之后停止，保证 Provider 刷新和关闭期间产生的最终记录仍可写入，随后再排空 goark-log。

## 验证

```bash
go test ./...
go test -race ./...
go vet ./...
```

使用 Apache License 2.0。
