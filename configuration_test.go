package gbclog_test

import (
	"bytes"
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"goark.dev/boot"
	"goark.dev/boot/configdata"
	gbclog "goark.dev/gbc-log"
	"goark.dev/goark"
	coreenv "goark.dev/goark/core/env"
	goarklog "goark.dev/log"
)

func TestAutoConfigureInstallsLoggerAndClosesContext(t *testing.T) {
	var output bytes.Buffer
	previous := slog.Default()
	app, err := boot.Run(t.Context(), boot.WithAutoConfiguration(gbclog.AutoConfigure(
		gbclog.WithLoggerContextFactory(func(context.Context, coreenv.Environment) (*goarklog.LoggerContext, error) {
			return goarklog.NewLoggerContext(goarklog.Options{
				Appenders: []goarklog.Appender{goarklog.NewConsoleAppender(goarklog.WithConsoleWriter(&output))},
				Root:      goarklog.RootLogger{Level: slog.LevelInfo, AppenderRefs: []string{"console"}},
			})
		}),
	)))
	if err != nil {
		t.Fatalf("boot.Run: %v", err)
	}
	appContext, _ := app.Context()
	logger := goark.MustGet[*slog.Logger](t.Context(), appContext, gbclog.BeanNameLogger)
	logger.InfoContext(t.Context(), "starter works", "key", "value")
	if err := app.Close(t.Context()); err != nil {
		t.Fatalf("Close: %v", err)
	}
	if slog.Default() != previous {
		t.Fatal("default logger was not restored")
	}
	if !strings.Contains(output.String(), "starter works") {
		t.Fatalf("output = %q", output.String())
	}
}

func TestAutoConfigureDisabledKeepsExistingLogger(t *testing.T) {
	previous := slog.Default()
	app, err := boot.Run(t.Context(), boot.WithAutoConfiguration(gbclog.AutoConfigure(gbclog.WithEnabled(false))))
	if err != nil {
		t.Fatalf("boot.Run: %v", err)
	}
	appContext, _ := app.Context()
	logger := goark.MustGet[*slog.Logger](t.Context(), appContext, gbclog.BeanNameLogger)
	if logger != previous {
		t.Fatal("disabled starter replaced existing logger")
	}
	system := goark.MustGet[gbclog.LoggingSystem](t.Context(), appContext, gbclog.BeanNameSystem)
	level := slog.LevelDebug
	if err := system.SetLogLevel("admin", &level); err == nil {
		t.Fatal("disabled LoggingSystem accepted level change")
	}
	if err := app.Close(t.Context()); err != nil {
		t.Fatalf("Close: %v", err)
	}
}

func TestLoggingSystemChangesAndRestoresLoggerLevel(t *testing.T) {
	app, err := boot.Run(t.Context(), boot.WithAutoConfiguration(gbclog.AutoConfigure()))
	if err != nil {
		t.Fatalf("boot.Run: %v", err)
	}
	defer func() { _ = app.Close(t.Context()) }()
	appContext, _ := app.Context()
	system := goark.MustGet[gbclog.LoggingSystem](t.Context(), appContext, gbclog.BeanNameSystem)
	level := slog.LevelDebug
	if err := system.SetLogLevel("admin.service", &level); err != nil {
		t.Fatalf("SetLogLevel(debug) error = %v", err)
	}
	configuration, found := system.LogLevel("admin.service")
	if !found || configuration.ConfiguredLevel == nil || *configuration.ConfiguredLevel != slog.LevelDebug || configuration.EffectiveLevel != slog.LevelDebug {
		t.Fatalf("logger configuration = %#v, found=%v", configuration, found)
	}
	if err := system.SetLogLevel("admin.service", nil); err != nil {
		t.Fatalf("SetLogLevel(nil) error = %v", err)
	}
	if _, found := system.LogLevel("admin.service"); found {
		t.Fatal("restored dynamic logger should no longer be explicitly configured")
	}
	if root, found := system.LogLevel("root"); !found || root.Name != "ROOT" {
		t.Fatalf("root configuration = %#v, found=%v", root, found)
	}
}

func TestRuntimeOrderKeepsLoggerAliveThroughProviderShutdown(t *testing.T) {
	if order := (&gbclog.Runtime{}).Order(); order != -11000 {
		t.Fatalf("runtime order = %d, want -11000", order)
	}
}

func TestAutoConfigure_whenLoggingPropertiesExist_shouldWriteConfiguredNamedLogger(t *testing.T) {
	root := t.TempDir()
	logFile := filepath.Join(root, "admin.log")
	config := "logging:\n" +
		"  level:\n" +
		"    root: info\n" +
		"    admin.service: debug\n" +
		"  pattern:\n" +
		"    file: '%level|%logger|%msg%n'\n" +
		"  file:\n" +
		"    name: '" + filepath.ToSlash(logFile) + "'\n"
	if err := os.WriteFile(filepath.Join(root, "app.yml"), []byte(config), 0o644); err != nil {
		t.Fatalf("write app config failed: %v", err)
	}

	app, err := boot.Run(
		t.Context(),
		boot.WithConfigDataOptions(configdata.WithLocations(root)),
		boot.WithAutoConfiguration(gbclog.AutoConfigure()),
	)
	if err != nil {
		t.Fatalf("boot.Run: %v", err)
	}
	appContext, _ := app.Context()
	loggerContext := goark.MustGet[*goarklog.LoggerContext](t.Context(), appContext, gbclog.BeanNameContext)
	loggerContext.Logger("admin.service").DebugContext(t.Context(), "named debug")
	loggerContext.Logger("other.service").DebugContext(t.Context(), "root debug")
	if err := app.Close(t.Context()); err != nil {
		t.Fatalf("Close: %v", err)
	}

	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("read log file failed: %v", err)
	}
	output := string(content)
	if !strings.Contains(output, "DEBUG|admin.service|named debug") {
		t.Fatalf("configured named logger output missing: %q", output)
	}
	if strings.Contains(output, "root debug") {
		t.Fatalf("root level did not suppress debug event: %q", output)
	}
}
