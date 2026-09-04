# Configuration Reference

[简体中文](configuration-reference.zh-CN.md)

## Starter

| Property | Default | Effect |
| --- | --- | --- |
| `goark.log.enabled` | `true` | Creates the managed goark-log runtime |
| `goark.log.install-default` | `true` | Installs it as `slog.Default()` |
| `logging.config` | auto-discovery | Selects an exact native goark-log resource |
| `logging.register-shutdown-hook` | `true` | Lets Boot close and drain goark-log |

## Levels And Groups

`logging.level.root` sets the root level. `logging.level.<logger>` sets a named
logger. `logging.group.<name>` defines comma-separated members and
`logging.level.<group>` expands a level to them. Built-in groups are `web` and
`sql`; a user definition replaces the built-in group, and a direct logger level
wins over group expansion. Supported levels are `ALL`, `TRACE`, `DEBUG`,
`INFO`, `WARN`, `ERROR`, `FATAL`, and `OFF`.

## Text Output

| Property | Default |
| --- | --- |
| `logging.console.enabled` | `true` |
| `logging.charset.console` / `logging.charset.file` | `UTF-8` |
| `logging.pattern.console` / `logging.pattern.file` | Goark Spring Boot style |
| `logging.pattern.dateformat` | `yyyy-MM-dd HH:mm:ss.SSS` |
| `logging.pattern.level` | `%5p` |
| `logging.pattern.correlation` | empty |
| `logging.exception-conversion-word` | `%wEx` |
| `logging.include-application-name` / `group` | `true` |
| `logging.threshold.console` / `file` | `TRACE` |

Identity comes from `goark.application.name` and `goark.application.group`.
Unknown charsets fail startup. UTF-8 adds no conversion layer.

## File And Rolling

| Property | Default |
| --- | --- |
| `logging.file.name` | unset |
| `logging.file.path` | unset |
| `logging.file.clean-history-on-start` | `false` |
| `logging.file.max-history` | `7` |
| `logging.file.max-size` | `10M` |
| `logging.file.total-size-cap` | unlimited |
| `logging.pattern.rolling-file-name` | `${LOG_FILE}.%d{yyyy-MM-dd}.%i.gz` |

`logging.file.name` wins over `logging.file.path`; the latter creates
`goark.log` in that directory. Sizes support case-insensitive `K/M/G/T`,
`KB/MB/GB/TB`, and `KiB/MiB/GiB/TiB`. Negative values and overflow fail.

## Structured JSON

| Property | Values / effect |
| --- | --- |
| `logging.structured.format.console` / `file` | `ecs`, `gelf`, `logstash`, or unset |
| `logging.structured.json.include` | Comma-separated full member paths |
| `logging.structured.json.exclude` | Comma-separated full member paths |
| `logging.structured.json.rename.<path>` | Renames a member |
| `logging.structured.json.add.<path>` | Adds a constant string member |
| `logging.structured.json.context.include` | Includes slog context; default `true` |
| `logging.structured.json.context.prefix` | Prefixes context names |
| `logging.structured.json.stacktrace.root` | `first` or `last` |
| `logging.structured.json.stacktrace.max-length` | Positive maximum UTF-8 bytes; truncated output ends in `...` |
| `logging.structured.json.stacktrace.max-throwable-depth` | Positive maximum frame count for every throwable |
| `logging.structured.json.stacktrace.include-common-frames` | Includes common frames |
| `logging.structured.json.stacktrace.include-hashes` | Adds stable frame hashes |
| `logging.structured.json.stacktrace.printer` | `standard` or `logging-system` |
| `logging.structured.ecs.service.environment` | ECS environment |
| `logging.structured.ecs.service.name` | ECS name; falls back to application name |
| `logging.structured.ecs.service.node-name` | ECS node name |
| `logging.structured.ecs.service.version` | ECS version |
| `logging.structured.gelf.host` | GELF host; falls back to application name |
| `logging.structured.gelf.service.version` | GELF service version |

Include is applied first, exclusion wins, then rename. Additions and Go
customizers use the same filter. Internal throwable, marker, thread, and stack
control attributes are not duplicated as context. Java class customizers fail
with guidance to use `WithStructuredJSONCustomizers`.

ECS interprets dots in context and added member paths as nested objects; GELF
and Logstash keep context names flat. ECS and Logstash marker hierarchies are
sorted `tags` arrays. Built-in protocol members win on a duplicate output name.
The default stacktrace printer is `logging-system` unless another stacktrace
property is configured, matching Spring Boot 4.1.1.

## Precedence

1. Boot config-data builds the final `Environment`.
2. `logging.config` loads native engine configuration when present.
3. `logging.level.*` overrides native root and named levels.
4. Thresholds override references to matching console/file appenders.
5. Other default-output properties apply only without explicit native config.

Invalid booleans, integers, sizes, levels, formats, resource locations, and
Java class customizers fail startup and identify the property in the error.
