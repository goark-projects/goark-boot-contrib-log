# Spring Boot 日志对齐说明

[English](spring-boot-logging-parity.md)

## 基线

稳定行为基线是 Spring Boot 4.1.1。Spring Boot 4.2.0-M1 只用于识别未来变化，
Goark 不会把里程碑版本的不稳定行为固化为契约。

## 已完成能力

- 根 Logger、命名 Logger、父级继承、日志组和原子动态级别。
- console/file threshold、关闭控制台和无输出时的 discard Appender。
- Pattern、Logger 缩写、字符集、文件滚动和完整保留策略。
- `logging.config` 的 YAML、TOML、JSON、XML、properties、文件和显式
  classpath 资源。
- 独立 ECS、GELF、Logstash Layout，以及 include/exclude/rename/add、上下文、
  异常顺序、深度、长度、公共帧和哈希。
- `LoggingSystem` Bean 的设置、恢复、单项查询和全量快照。
- `logging.register-shutdown-hook` 对 Boot 生命周期所有权的映射。

## Go 化差异

- 按用户要求，Goark 默认日期为 `yyyy-MM-dd HH:mm:ss.SSS`；Spring Boot
  4.1.1 自身默认值是带偏移的 ISO 时间。
- Goark 后端是 goark-log，不存在 Logback/Log4j2，因此不把
  `logging.logback.*`、`logging.log4j2.*` 固化为 Goark 契约；等价能力由
  后端无关属性提供。
- Go 不能加载 Java 全限定类；Formatter、StackTracePrinter、JSON Customizer
  使用显式、类型安全的 Go Option/Interface。
- 当前 Boot 在 `gbc-log` 创建前没有框架日志。在 `AutoConfigure()` 阶段安装
  全局缓存会在配置加载失败时泄漏进程状态；未来若增加前置日志，Boot 必须先提供
  失败安全的 bootstrap 生命周期。
- shutdown 属性复用 Boot 生命周期，不注册第二套信号处理。

## 性能契约

属性解析、变换编译、格式选择和资源加载都在启动阶段完成。结构化热路径按顺序
直接写 JSON，不走每事件 `map[string]any + json.Marshal`。动态级别原子发布
不可变路由快照，日志调用不获取配置互斥锁。自定义 Customizer 的分配成本由
实现方负责。

正确性门禁包括单元测试、race、vet、独立模块和 `admin-minimal` 运行测试。
性能结论必须来自同机、同负载基准。
