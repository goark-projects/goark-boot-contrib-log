# Goark Boot Contrib Log

[English](README.md)

`goark.dev/gbc-log` 负责 `goark.dev/log` 与 Goark Boot 的集成：映射
`logging.*`，安装 primary `*slog.Logger`，提供运行期级别控制，并保证日志
运行期只有一个生命周期所有者。

稳定对齐基线是 Spring Boot 4.1.1；Spring Boot 4.2.0-M1 只用于前向审计，
里程碑版本中的可变行为不会成为 Goark 稳定契约。

## 职责边界

- `goark.dev/log` 负责 Handler、路由、Appender、Layout、Filter、异步投递、
  原生配置、重载和编码。
- `goark.dev/gbc-log` 负责 Boot 属性映射、Bean、默认 Logger 和生命周期顺序。
- 框架代码只依赖 `log/slog` 或 `goark.dev/log`；日志引擎不依赖 Boot 和观测。

## 快速使用

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
  group:
    web: goark.dev.arkhos,goark.dev.goark.web
  file:
    name: logs/admin.log
    max-size: 20M
    max-history: 14
    total-size-cap: 1G
```

默认日期格式为 `yyyy-MM-dd HH:mm:ss.SSS`，`goark.dev.arkhos.hertz` 默认缩写
为 `g.d.arkhos.hertz`。

## 结构化日志

控制台和文件可分别选择 `ecs`、`gelf`、`logstash`。`logging.structured.json`
支持 include、exclude、rename、add、context 和 stacktrace 控制；ECS 支持服务
元数据，GELF 支持 host 和服务版本。Java 类 Customizer 会被拒绝，Go 实现通过
`WithStructuredJSONCustomizers` 显式注册。

## 原生配置

`logging.config` 可指定 goark-log YAML、TOML、JSON、XML 或 properties 资源，
支持相对路径、绝对路径、`file:` 和显式注册文件系统中的 `classpath:`。资源
缺失、不可读、格式错误或存在未知字段时启动失败。

原生配置负责 Appender、Layout、Filter、additivity、异步和重载；
`logging.level.*` 及匹配的 console/file threshold 作为 Boot 覆盖项继续生效。
默认 pattern、file、charset、rolling、structured 属性不会替换原生 Appender。

## 运行期 API

- `goark.log.logger`：primary `*slog.Logger`。
- `goark.log.context`：生命周期中性的 `*goarklog.LoggerContext` 访问别名。
- `goark.log.system`：原子修改级别和查询快照的 `LoggingSystem`。
- `goark.log.lifecycle`：恢复默认 Logger、排空和关闭日志的唯一所有者。

`logging.register-shutdown-hook=true` 表示由 Boot 关闭并排空 goark-log，不会注册
第二套 OS 信号 Hook。只有外部组件明确负责关闭时才应设置为 `false`。

## 参考

- [配置参考](docs/configuration-reference.zh-CN.md)
- [Spring Boot 对齐说明](docs/spring-boot-logging-parity.zh-CN.md)

## 验证

```bash
go test ./...
go test -race ./...
go vet ./...
GOWORK=off go test ./...
```

使用 Apache License 2.0。
