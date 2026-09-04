# Goark Boot Contrib Log

[简体中文](README.zh-CN.md)

`goark.dev/gbc-log` integrates `goark.dev/log` with Goark Boot. It maps
`logging.*`, installs the primary `*slog.Logger`, exposes runtime level control,
and gives the logging runtime one lifecycle owner.

The stable compatibility baseline is Spring Boot 4.1.1. Spring Boot 4.2.0-M1
is audited for forward compatibility only; milestone behavior is not a stable
Goark contract.

## Ownership

- `goark.dev/log` owns handlers, routing, appenders, layouts, filters, async
  delivery, native configuration, reload, and encoding.
- `goark.dev/gbc-log` owns Boot property mapping, beans, default logger
  installation, and lifecycle ordering.
- Framework code logs through `log/slog` or `goark.dev/log`; the engine does not
  depend on Boot or observability.

## Quick Start

```go
app, err := boot.Run(ctx, boot.WithAutoConfiguration(gbclog.AutoConfigure()))
```

```yaml
goark:
  application:
    name: admin
  log:
    enabled: true
    install-default: true

logging:
  level:
    root: info
    web: debug
    goark.dev.orm: warn
  group:
    web: goark.dev.arkhos,goark.dev.goark.web
  file:
    name: logs/admin.log
    max-size: 20M
    max-history: 14
    total-size-cap: 1G
  threshold:
    console: info
    file: debug
```

Default output uses `yyyy-MM-dd HH:mm:ss.SSS`. The logger
`goark.dev.arkhos.hertz` is abbreviated as `g.d.arkhos.hertz` by default.

## Structured Logging

Console and file outputs independently support `ecs`, `gelf`, and `logstash`:

```yaml
logging:
  structured:
    format:
      console: ecs
      file: logstash
    json:
      context:
        include: true
        prefix: ctx.
      exclude: process.thread.name
      rename:
        message: msg
      add:
        deployment: production
      stacktrace:
        root: first
        max-length: 8192
        max-throwable-depth: 8
        include-common-frames: false
        include-hashes: true
    ecs:
      service:
        environment: production
        version: 1.2.0
```

Java class names in `logging.structured.json.customizer` are rejected. Use the
typed `WithStructuredJSONCustomizers` option for Go implementations.

## Native Configuration

`logging.config` accepts a goark-log YAML, TOML, JSON, XML, or properties
resource. Relative paths, absolute paths, `file:`, and explicitly registered
`classpath:` resources are supported. Missing, unreadable, malformed, or
unknown resources fail startup.

Native configuration remains authoritative for appenders, layouts, filters,
additivity, async behavior, and reload. `logging.level.*` and matching
console/file thresholds are applied as Boot overrides. Default pattern, file,
charset, rolling, and structured properties do not replace native appenders.

## Runtime API

- `goark.log.logger`: primary `*slog.Logger`.
- `goark.log.context`: lifecycle-neutral `*goarklog.LoggerContext` alias.
- `goark.log.system`: `gbclog.LoggingSystem` for atomic level changes and
  snapshots.
- `goark.log.lifecycle`: sole owner of restoration and shutdown.

```go
system := goark.MustGet[gbclog.LoggingSystem](ctx, appContext, gbclog.BeanNameSystem)
debug := slog.LevelDebug
_ = system.SetLogLevel("admin.service", &debug)
configuration, found := system.LogLevel("admin.service")
_ = system.SetLogLevel("admin.service", nil)
```

`logging.register-shutdown-hook=true` means Boot closes and drains goark-log.
No second OS signal hook is registered. Set it to `false` only when another
owner closes the `LoggerContext`.

## Reference

- [Configuration reference](docs/configuration-reference.md)
- [Spring Boot parity](docs/spring-boot-logging-parity.md)

## Verification

```bash
go test ./...
go test -race ./...
go vet ./...
GOWORK=off go test ./...
```

Licensed under Apache License 2.0.
