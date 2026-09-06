# Changelog

English | [中文](CHANGELOG.zh-CN.md)

All notable changes to Goark Boot Contrib Log are recorded here.

## [Unreleased]

No unreleased changes.

## [0.0.1] - 2026-09-06

### Added

- Goark Boot auto-configuration for `goark.dev/log`.
- `logging.*` and `goark.log.*` property mapping aligned with Spring Boot 4.1.
- Console, rolling file, ECS, GELF, and Logstash output configuration.
- Logger groups, runtime level control, charset selection, resource resolution,
  and one lifecycle owner for logging shutdown.
- Cross-platform CI with Go 1.26 tests, vet, and race gates.

### Changed

- Aligned all used `golang.org/x` modules with their latest stable releases.

### Fixed

- Logging remains available until telemetry and application shutdown complete.
- Default console output uses stdout and applies implicit appender thresholds.
- Structured GELF and JSON settings follow the documented Boot semantics.

[Unreleased]: https://github.com/goark-projects/goark-boot-contrib-log/compare/v0.0.1...HEAD
[0.0.1]: https://github.com/goark-projects/goark-boot-contrib-log/releases/tag/v0.0.1
