package gbclog

import (
	"context"
	"log/slog"
	"path/filepath"
	"testing"

	coreenv "goark.dev/goark/core/env"
	goarklog "goark.dev/log"
)

func TestReadLoggingProperties_whenLevelsAndGroupsExist_shouldResolveDirectLevelLast(t *testing.T) {
	environment := newLoggingEnvironment(t, map[string]any{
		PropertyRootLevel:                       "warn",
		PropertyGroupPrefix + "web":             "goark.dev.arkhos, goark.dev.goark",
		PropertyLevelPrefix + "web":             "debug",
		PropertyLevelPrefix + "goark.dev.goark": "error",
		PropertyConsoleThreshold:                "info",
	})

	properties, err := readLoggingProperties(environment)
	if err != nil {
		t.Fatalf("read logging properties failed: %v", err)
	}
	if properties.rootLevel == nil || *properties.rootLevel != slog.LevelWarn {
		t.Fatalf("root level = %v, want WARN", properties.rootLevel)
	}
	if got := properties.loggerLevels["goark.dev.arkhos"]; got != slog.LevelDebug {
		t.Fatalf("group logger level = %s, want DEBUG", got)
	}
	if got := properties.loggerLevels["goark.dev.goark"]; got != slog.LevelError {
		t.Fatalf("direct logger level = %s, want ERROR", got)
	}
}

func TestReadLoggingProperties_whenLevelIsInvalid_shouldReturnError(t *testing.T) {
	environment := newLoggingEnvironment(t, map[string]any{PropertyRootLevel: "invalid"})
	if _, err := readLoggingProperties(environment); err == nil {
		t.Fatal("invalid logger level must fail startup")
	}
}

func TestLoggingOptionsCustomizer_whenDefaultSource_shouldBuildDefaultOutputs(t *testing.T) {
	logFile := filepath.Join(t.TempDir(), "admin.log")
	environment := newLoggingEnvironment(t, map[string]any{
		PropertyRootLevel:        "debug",
		PropertyConsolePattern:   "%level %msg%n",
		PropertyFilePattern:      "%msg%n",
		PropertyFileName:         logFile,
		PropertyConsoleThreshold: "warn",
		PropertyFileThreshold:    "error",
	})
	customize := loggingOptionsCustomizer(environment)
	options, err := customize(context.Background(), goarklog.DefaultOptions(), &goarklog.ConfigResult{Source: goarklog.ConfigSourceDefault})
	if err != nil {
		t.Fatalf("customize options failed: %v", err)
	}
	defer closeAppenders(t, options.Appenders)
	if options.Root.Level != slog.LevelDebug {
		t.Fatalf("root level = %s, want DEBUG", options.Root.Level)
	}
	if len(options.Appenders) != 2 || options.Appenders[0].Name() != "console" || options.Appenders[1].Name() != "file" {
		t.Fatalf("appenders = %#v, want console and file", options.Appenders)
	}
	if len(options.Root.AppenderRefControls) != 2 {
		t.Fatalf("root controls = %#v, want two thresholds", options.Root.AppenderRefControls)
	}
}

func TestLoggingOptionsCustomizer_whenFileConfigExists_shouldPreserveAppenderAndRouteFields(t *testing.T) {
	console := goarklog.NewConsoleAppender(goarklog.WithConsoleName("CONSOLE"))
	location := true
	original := goarklog.Options{
		Appenders: []goarklog.Appender{console},
		Root: goarklog.RootLogger{
			Level: slog.LevelInfo,
			AppenderRefControls: []goarklog.AppenderRef{{
				Ref:             "CONSOLE",
				IncludeLocation: &location,
			}},
		},
		Loggers: []goarklog.LoggerRule{{Name: "example", Additivity: false, AdditivitySet: true, AppenderRefs: []string{"CONSOLE"}}},
	}
	environment := newLoggingEnvironment(t, map[string]any{
		PropertyConsolePattern:          "%msg%n",
		PropertyConsoleThreshold:        "error",
		PropertyLevelPrefix + "example": "debug",
	})
	customize := loggingOptionsCustomizer(environment)
	options, err := customize(context.Background(), original, &goarklog.ConfigResult{Source: goarklog.ConfigSourceExplicit})
	if err != nil {
		t.Fatalf("customize options failed: %v", err)
	}
	defer closeAppenders(t, options.Appenders)
	if len(options.Appenders) != 1 || options.Appenders[0] != console {
		t.Fatal("explicit appender must be preserved")
	}
	control := options.Root.AppenderRefControls[0]
	if control.IncludeLocation == nil || !*control.IncludeLocation || control.Level == nil || *control.Level != slog.LevelError {
		t.Fatalf("root control was not preserved with threshold override: %#v", control)
	}
	rule := options.Loggers[0]
	if rule.Level == nil || *rule.Level != slog.LevelDebug || !rule.AdditivitySet || rule.Additivity {
		t.Fatalf("logger route fields were not preserved: %#v", rule)
	}
	if len(rule.AppenderRefs) != 0 || len(rule.AppenderRefControls) != 1 || rule.AppenderRefControls[0].Level == nil || *rule.AppenderRefControls[0].Level != slog.LevelError {
		t.Fatalf("logger appender threshold was not applied: %#v", rule)
	}
}

func newLoggingEnvironment(t *testing.T, values map[string]any) coreenv.Environment {
	t.Helper()
	environment, err := coreenv.NewStandardEnvironment()
	if err != nil {
		t.Fatalf("new environment failed: %v", err)
	}
	source, err := coreenv.NewMapPropertySource("test", values)
	if err != nil {
		t.Fatalf("new property source failed: %v", err)
	}
	if err := environment.PropertySources().AddFirst(source); err != nil {
		t.Fatalf("add property source failed: %v", err)
	}
	return environment
}

func closeAppenders(t *testing.T, appenders []goarklog.Appender) {
	t.Helper()
	for _, appender := range appenders {
		if err := appender.Close(); err != nil {
			t.Errorf("close appender %q failed: %v", appender.Name(), err)
		}
	}
}
