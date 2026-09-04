# Spring Boot Logging Parity

[简体中文](spring-boot-logging-parity.zh-CN.md)

## Baseline

Spring Boot 4.1.1 is the stable behavioral baseline. Spring Boot 4.2.0-M1 is
reviewed only for upcoming changes; unstable milestone behavior is not a Goark
contract.

## Capability Matrix

| Capability | Goark behavior |
| --- | --- |
| Root, named, inherited levels | Supported, including atomic runtime overrides |
| Logging groups | Built-in `web` and `sql`, user replacement, direct-level precedence |
| Console/file thresholds | Supported for default and matching native appenders |
| Console disablement | Supported; discard output keeps an engine valid when no file exists |
| Pattern and logger abbreviation | Supported; requested Goark defaults are applied |
| File rolling and retention | Startup cleanup, history, file size, and total-size cap |
| Output charset | Supported; unknown encodings fail startup |
| `logging.config` | YAML, TOML, JSON, XML, properties, file and registered classpath resources |
| ECS, GELF, Logstash | Dedicated layouts, not aliases of generic JSON |
| JSON transformations | Include, exclude, rename, add, context, and stacktrace controls |
| Runtime control | `LoggingSystem` set, restore, lookup, and snapshot API |
| Shutdown hook | Mapped to Boot lifecycle ownership; no duplicate OS hook |

## Intentional Go Differences

- Goark defaults to the requested `yyyy-MM-dd HH:mm:ss.SSS`; Spring Boot 4.1.1
  itself defaults to an ISO offset timestamp.
- Goark uses goark-log, not Logback or Log4j2. Backend-specific
  `logging.logback.*` and `logging.log4j2.*` keys are not Goark contracts;
  equivalent rolling behavior uses backend-neutral properties.
- Java fully-qualified formatter, stack printer, and JSON customizer classes
  cannot be loaded. Goark uses explicit typed options and interfaces.
- Boot currently emits no framework logs before `gbc-log` exists. Installing a
  global buffer during `AutoConfigure()` would leak process state if config-data
  loading failed. A future pre-context logger requires a failure-safe Boot
  bootstrap lifecycle first.
- `logging.register-shutdown-hook` controls the existing Boot owner and never
  installs another signal handler.

## Performance Contract

Property parsing, transforms, format selection, and resource loading run at
startup. Structured output writes JSON in order without a per-event
`map[string]any` plus generic marshal. Optional customizers own their allocation
cost. Dynamic level changes atomically publish immutable routing snapshots, so
log calls do not acquire the configuration mutex.

Correctness gates are unit, race, vet, isolated-module, and `admin-minimal`
runtime tests. Performance claims must remain same-host and workload-specific.
