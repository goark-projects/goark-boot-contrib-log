# 变更日志

[English](CHANGELOG.md) | 中文

这里记录 Goark Boot Contrib Log 的重要变更。

## [未发布]

暂无未发布变更。

## [0.0.1] - 2026-09-07

### 新增

- 为 `goark.dev/log` 提供 Goark Boot 自动配置。
- 与 Spring Boot 4.1 对齐的 `logging.*` 和 `goark.log.*` 属性映射。
- Console、滚动文件、ECS、GELF 和 Logstash 输出配置。
- Logger Group、运行时级别控制、字符集选择、资源解析和日志关闭的单一生命周期所有者。
- 基于 Go 1.26 的跨平台测试、vet 和 race 门禁。

### 变更

- 将所有实际使用的 `golang.org/x` 模块对齐到最新稳定版本。

### 修复

- 日志保持可用，直到遥测和应用关闭完成。
- 默认控制台输出使用 stdout，并应用隐式 Appender Threshold。
- 结构化 GELF 和 JSON 设置遵循文档化的 Boot 语义。

[未发布]: https://github.com/goark-projects/goark-boot-contrib-log/compare/v0.0.1...HEAD
[0.0.1]: https://github.com/goark-projects/goark-boot-contrib-log/releases/tag/v0.0.1
