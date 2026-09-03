# Goark Boot Contrib Log

[简体中文](README.zh-CN.md)

`goark.dev/gbc-log` is the Goark Boot starter for `goark.dev/log`. It creates the application logging runtime, exposes the primary `*slog.Logger`, optionally installs it as `slog.Default`, and drains logging before observability exporters shut down.

## Usage

```go
app, err := boot.Run(ctx, boot.WithAutoConfiguration(
    gbclog.AutoConfigure(),
))
```

The starter reuses goark-log configuration discovery. Configure a file through `goark.log.config`, `goark.logging.config`, `logging.config`, `GOARK_LOG_CONFIG`, or the standard `conf/goark-log.*` paths.

## Beans

- `goark.log.logger`: primary `*slog.Logger`.
- `goark.log.context`: managed `*goarklog.LoggerContext` when enabled.
- `goark.log.lifecycle`: runtime responsible for default-logger restoration and shutdown.

## Properties

- `goark.log.enabled`: defaults to `true`.
- `goark.log.install-default`: defaults to `true`.
- `goark.log.config`: optional goark-log configuration path.

## Lifecycle

The runtime order is `-9000`. It stops after ordinary application components and before the observability provider at `-10000`, so telemetry exported to goark-log is drained before exporter shutdown.

## Validation

```bash
go test ./...
go test -race ./...
go vet ./...
```

Licensed under Apache License 2.0.
