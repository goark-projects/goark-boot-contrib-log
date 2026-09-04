# 配置参考

[English](configuration-reference.md)

## Starter

| 属性 | 默认值 | 作用 |
| --- | --- | --- |
| `goark.log.enabled` | `true` | 创建受管 goark-log 运行期 |
| `goark.log.install-default` | `true` | 安装为 `slog.Default()` |
| `logging.config` | 自动发现 | 指定原生 goark-log 配置资源 |
| `logging.register-shutdown-hook` | `true` | 由 Boot 关闭并排空日志 |

## 级别和日志组

`logging.level.root` 配置根级别，`logging.level.<logger>` 配置命名 Logger，
`logging.group.<name>` 定义逗号分隔的成员，`logging.level.<group>` 展开组级别。
内置组为 `web` 和 `sql`；用户同名组替换内置组，直接 Logger 级别覆盖组展开。
支持 `ALL`、`TRACE`、`DEBUG`、`INFO`、`WARN`、`ERROR`、`FATAL`、`OFF`。

## 文本输出

| 属性 | 默认值 |
| --- | --- |
| `logging.console.enabled` | `true` |
| `logging.charset.console` / `file` | `UTF-8` |
| `logging.pattern.console` / `file` | Goark Spring Boot 风格 |
| `logging.pattern.dateformat` | `yyyy-MM-dd HH:mm:ss.SSS` |
| `logging.pattern.level` | `%5p` |
| `logging.pattern.correlation` | 空 |
| `logging.exception-conversion-word` | `%wEx` |
| `logging.include-application-name` / `group` | `true` |
| `logging.threshold.console` / `file` | `TRACE` |

应用名称和组来自 `goark.application.name`、`goark.application.group`。未知字符集
会导致启动失败；UTF-8 不增加转换层。

## 文件和滚动

`logging.file.name` 优先于 `logging.file.path`，后者在指定目录创建 `goark.log`。
支持 `clean-history-on-start`、`max-history`、`max-size`、`total-size-cap` 和
`logging.pattern.rolling-file-name`。容量支持大小写不敏感的 `K/M/G/T`、
`KB/MB/GB/TB`、`KiB/MiB/GiB/TiB`；负数和溢出会导致启动失败。默认值分别为
`false`、`7`、`10M`、不限制和 `${LOG_FILE}.%d{yyyy-MM-dd}.%i.gz`。

## 结构化 JSON

`logging.structured.format.console/file` 支持 `ecs`、`gelf`、`logstash`。
`logging.structured.json` 支持 `include`、`exclude`、`rename`、`add`、
`context.include`、`context.prefix`，以及异常的 `root`、`max-length`、
`max-throwable-depth`、`include-common-frames`、`include-hashes`、`printer`。
ECS 支持 service environment/name/node-name/version；GELF 支持 host 和 service
version。

执行顺序为 include、exclude（优先）、rename。常量和 Go Customizer 走相同
过滤链。Throwable、marker、线程和上下文栈控制属性不会重复输出。Java 类
Customizer 会启动失败，必须使用 `WithStructuredJSONCustomizers`。

## 优先级

1. Boot config-data 合成最终 `Environment`。
2. `logging.config` 存在时加载原生引擎配置。
3. `logging.level.*` 覆盖原生根级别和命名 Logger 级别。
4. threshold 覆盖匹配的 console/file Appender 引用。
5. 其他默认输出属性只在未选择原生配置时生效。

非法 bool、整数、容量、级别、格式、资源地址和 Java 类 Customizer 都会中止
启动，错误包含具体属性名。
