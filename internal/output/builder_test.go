package output

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
	"time"

	"goark.dev/gbc-log/internal/properties"
	goarklog "goark.dev/log"
)

func TestApplyDefault_whenConsoleDisabledWithoutFile_shouldInstallDiscardAppender(t *testing.T) {
	configuration, err := properties.Read(nil)
	if err != nil {
		t.Fatalf("Read(nil) error = %v", err)
	}
	configuration.ConsoleEnabled = false
	options, err := ApplyDefault(goarklog.DefaultOptions(), configuration)
	if err != nil {
		t.Fatalf("ApplyDefault() error = %v", err)
	}
	defer closeAppenders(t, options.Appenders)
	if len(options.Appenders) != 1 || options.Appenders[0].Name() != "discard" {
		t.Fatalf("appenders = %#v", options.Appenders)
	}
}

func TestApplyDefault_whenFileConfigured_shouldUseRollingFileAppender(t *testing.T) {
	configuration, err := properties.Read(nil)
	if err != nil {
		t.Fatalf("Read(nil) error = %v", err)
	}
	configuration.ConsoleEnabled = false
	configuration.FileName = t.TempDir() + "/admin.log"
	options, err := ApplyDefault(goarklog.DefaultOptions(), configuration)
	if err != nil {
		t.Fatalf("ApplyDefault() error = %v", err)
	}
	defer closeAppenders(t, options.Appenders)
	if len(options.Appenders) != 1 {
		t.Fatalf("appenders = %#v", options.Appenders)
	}
	if _, ok := options.Appenders[0].(*goarklog.RollingFileAppender); !ok {
		t.Fatalf("file appender type = %T", options.Appenders[0])
	}
	if options.Appenders[0].Name() != "file" {
		t.Fatalf("file appender name = %q", options.Appenders[0].Name())
	}
}

func TestDefaultPattern_shouldRenderConfiguredIdentityAndGoarkDateFormat(t *testing.T) {
	configuration, err := properties.Read(nil)
	if err != nil {
		t.Fatalf("Read(nil) error = %v", err)
	}
	configuration.ApplicationName = "admin%service"
	configuration.ApplicationGroup = "backend"
	layout, err := goarklog.NewPatternLayout(defaultPattern(configuration, false))
	if err != nil {
		t.Fatalf("NewPatternLayout() error = %v", err)
	}
	var output bytes.Buffer
	err = layout.Format(&output, goarklog.Event{
		Time:    time.Date(2026, 9, 3, 20, 16, 52, 951_000_000, time.Local),
		Level:   slog.LevelInfo,
		Logger:  "goark.dev.arkhos.hertz",
		Message: "started",
	})
	if err != nil {
		t.Fatalf("Format() error = %v", err)
	}
	text := output.String()
	for _, wanted := range []string{
		"2026-09-03 20:16:52.951",
		"[admin%service] [backend]",
		"g.d.arkhos.hertz",
		"started",
	} {
		if !strings.Contains(text, wanted) {
			t.Fatalf("output %q does not contain %q", text, wanted)
		}
	}
}

func TestTextLayout_whenFileCharsetConfigured_shouldEncodeOutput(t *testing.T) {
	configuration, err := properties.Read(nil)
	if err != nil {
		t.Fatalf("Read(nil) error = %v", err)
	}
	configuration.FilePattern = "%msg"
	configuration.FileCharset = "ISO-8859-1"
	layout, err := outputLayout(configuration, true, nil)
	if err != nil {
		t.Fatalf("textLayout() error = %v", err)
	}
	var output bytes.Buffer
	if err := layout.Format(&output, goarklog.Event{Message: "caf\u00e9"}); err != nil {
		t.Fatalf("Format() error = %v", err)
	}
	if !bytes.Equal(output.Bytes(), []byte{'c', 'a', 'f', 0xe9}) {
		t.Fatalf("encoded output = %v", output.Bytes())
	}
}

func TestTextLayout_whenCharsetUnknown_shouldFail(t *testing.T) {
	configuration, err := properties.Read(nil)
	if err != nil {
		t.Fatalf("Read(nil) error = %v", err)
	}
	configuration.ConsoleCharset = "unknown-charset"
	if _, err := outputLayout(configuration, false, nil); err == nil {
		t.Fatal("unknown charset should fail")
	}
}

func TestOutputLayout_whenStructuredFormatConfigured_shouldMapProperties(t *testing.T) {
	configuration, err := properties.Read(nil)
	if err != nil {
		t.Fatalf("Read(nil) error = %v", err)
	}
	configuration.Structured.ConsoleFormat = "ecs"
	configuration.ApplicationName = "admin"
	configuration.Structured.ECS.ServiceVersion = "1.0.0"
	configuration.Structured.JSON.ContextInclude = boolPointer(true)
	configuration.Structured.JSON.ContextPrefix = "ctx."
	configuration.Structured.JSON.Add["build"] = "42"
	configuration.Structured.JSON.Rename["message"] = "msg"
	layout, err := outputLayout(configuration, false, []goarklog.StructuredJSONCustomizer{
		goarklog.StructuredJSONCustomizerFunc(func(_ goarklog.Event, fields goarklog.StructuredJSONFieldAppender) {
			fields.Add("custom", slog.StringValue("ok"))
		}),
	})
	if err != nil {
		t.Fatalf("outputLayout() error = %v", err)
	}
	var output bytes.Buffer
	if err := layout.Format(&output, goarklog.Event{
		Time: time.Date(2026, 9, 4, 10, 0, 0, 0, time.UTC), Logger: "admin", Message: "started",
		Attrs: []slog.Attr{slog.String("trace", "trace-1")},
	}); err != nil {
		t.Fatalf("Format() error = %v", err)
	}
	text := output.String()
	for _, wanted := range []string{`"service":{"name":"admin","version":"1.0.0"}`, `"msg":"started"`, `"ctx":{"trace":"trace-1"}`, `"build":"42"`, `"custom":"ok"`} {
		if !strings.Contains(text, wanted) {
			t.Fatalf("structured output %q does not contain %q", text, wanted)
		}
	}
}

func TestStructuredStacktracePrinterMatchesBootDefaults(t *testing.T) {
	if actual := structuredStacktracePrinter(properties.StacktraceProperties{}); actual != goarklog.StructuredStacktracePrinterLoggingSystem {
		t.Fatalf("default printer = %q", actual)
	}
	includeCommonFrames := false
	configured := properties.StacktraceProperties{IncludeCommonFrames: &includeCommonFrames}
	if actual := structuredStacktracePrinter(configured); actual != goarklog.StructuredStacktracePrinterStandard {
		t.Fatalf("configured printer = %q", actual)
	}
	configured.Printer = "logging-system"
	if actual := structuredStacktracePrinter(configured); actual != goarklog.StructuredStacktracePrinterLoggingSystem {
		t.Fatalf("explicit logging-system printer = %q", actual)
	}
}

func boolPointer(value bool) *bool { return &value }

func closeAppenders(t *testing.T, appenders []goarklog.Appender) {
	t.Helper()
	for _, appender := range appenders {
		if err := appender.Close(); err != nil {
			t.Errorf("Close(%q) error = %v", appender.Name(), err)
		}
	}
}
