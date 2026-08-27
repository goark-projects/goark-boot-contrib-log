# Goark Boot Contrib Log

Official Goark Boot starter module for application logging.

## Module

- Module path: `goark.dev/gbc-log`
- Repository: `github.com/goark-projects/goark-boot-contrib-log`
- License: Apache-2.0
- Default branch: `main`
- Development branch: `dev`

## Scope

`goark.dev/gbc-log` is the Goark-managed logging starter. It is intended to provide the Spring Boot logging starter equivalent for the Goark ecosystem:

- Goark Boot logging auto-configuration.
- Defaults backed by `goark.dev/log`.
- Structured logging conventions for application startup, lifecycle, and request flows.
- Future integration points for observability, tracing, and runtime diagnostics.

The initial repository bootstrap exposes stable module metadata only. Runtime auto-configuration APIs will be added in later implementation slices.

## Development

```bash
go test ./...
```

## Chinese

# Goark Boot Contrib Log（中文）

Goark 官方维护的应用日志启动器模块。

## 模块信息

- 模块路径：`goark.dev/gbc-log`
- 仓库地址：`github.com/goark-projects/goark-boot-contrib-log`
- 开源协议：Apache-2.0
- 默认分支：`main`
- 开发分支：`dev`

## 职责边界

`goark.dev/gbc-log` 是 Goark 生态中对标 Spring Boot 日志启动器的官方模块：

- 提供 Goark Boot 日志自动配置。
- 默认接入 `goark.dev/log`。
- 统一应用启动、生命周期和请求链路的结构化日志约定。
- 后续扩展可观测性、链路追踪和运行期诊断集成点。

当前初始化版本只暴露稳定的模块元数据。运行期自动配置 API 会在后续实现切片中补齐。

## 开发

```bash
go test ./...
```
