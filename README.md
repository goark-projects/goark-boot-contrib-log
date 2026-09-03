# Goark Boot Contrib Log

[简体中文](README.zh-CN.md)

`goark.dev/gbc-log` integrates Goark Boot with `goark.dev/log`. It creates the
logging runtime before infrastructure starters, exposes the primary
`*slog.Logger`, installs it as `slog.Default` by default, and drains it after
observability exporters stop.

## Responsibilities

- `goark.dev/log` owns logger APIs, handlers, appenders, layouts, filters,
  asynchronous delivery, configuration files, and reload.
- `goark.dev/gbc-log` reads Boot's `Environment`, maps `logging.*`, registers
  beans, installs the process default logger, and manages lifecycle ordering.
- Other libraries depend on `log/slog` or `goark.dev/log`; they do not depend on
  this starter.

## Usage

```go
app, err := boot.Run(ctx, boot.WithAutoConfiguration(
	gbclog.AutoConfigure(),
))
```

## Properties

| Property | Default | Description |
| --- | --- | --- |
| `goark.log.enabled` | `true` | Enables the managed goark-log runtime |
| `goark.log.install-default` | `true` | Installs the logger as `slog.Default` |
| `logging.config` | auto-discovery | Exact goark-log config file path |
| `logging.level.root` | `INFO` | Root logger level |
| `logging.level.<logger>` | inherited | Named logger level |
| `logging.group.<name>` | unset | Comma-separated logger names in a group |
| `logging.level.<group>` | unset | Level applied to all group members |
| `logging.pattern.console` | Spring Boot style | Default console PatternLayout pattern |
| `logging.pattern.file` | Spring Boot style | Default file PatternLayout pattern |
| `logging.file.name` | unset | Exact default log file path |
| `logging.file.path` | unset | Directory for the default `goark.log` file |
| `logging.threshold.console` | unset | Console appender threshold |
| `logging.threshold.file` | unset | File appender threshold |

`logging.file.name` takes precedence over `logging.file.path`. Levels accept all
names registered in goark-log, including `ALL`, `TRACE`, `DEBUG`, `INFO`,
`WARN`, `ERROR`, `FATAL`, and `OFF`.

```yaml
goark:
  application:
    name: admin

logging:
  level:
    root: info
    web: debug
    goark.dev.orm: warn
  group:
    web: goark.dev.arkhos,goark.dev.goark
  pattern:
    console: "%d{yyyy-MM-dd HH:mm:ss.SSS} %5level %logger : %msg%n"
    file: "%d{ISO8601} %level %logger %msg%n"
  file:
    name: logs/admin.log
  threshold:
    console: info
    file: debug
```

Direct logger levels override levels expanded from groups.

## Independent Configuration

`logging.config` selects a native goark-log YAML, TOML, JSON, XML, or properties
configuration file. When one is selected:

- Its appenders, layouts, filters, additivity, async settings, and reload policy
  remain authoritative.
- `logging.level.*` overrides root and named logger levels.
- `logging.threshold.console` and `logging.threshold.file` override matching
  appender-reference thresholds while preserving filters and location settings.
- `logging.pattern.*` and `logging.file.*` only customize the starter's default
  outputs and do not replace explicit appenders or layouts.

## Beans

- `goark.log.logger`: primary `*slog.Logger`.
- `goark.log.context`: managed `*goarklog.LoggerContext` when enabled.
- `goark.log.lifecycle`: default-logger restoration and shutdown runtime.

The runtime order is `-11000`. It starts before the observability provider at
`-10000` and stops after it, allowing exporter flush and shutdown records to be
written before goark-log drains.

## Validation

```bash
go test ./...
go test -race ./...
go vet ./...
```

Licensed under Apache License 2.0.
