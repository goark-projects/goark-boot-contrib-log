package properties

import (
	"log/slog"
	"testing"

	coreenv "goark.dev/goark/core/env"
)

func TestRead_whenEnvironmentIsNil_shouldReturnGoarkDefaults(t *testing.T) {
	properties, err := Read(nil)
	if err != nil {
		t.Fatalf("Read(nil) error = %v", err)
	}
	if !properties.ConsoleEnabled || !properties.IncludeApplicationName || !properties.IncludeApplicationGroup || !properties.RegisterShutdownHook {
		t.Fatalf("boolean defaults = %#v", properties)
	}
	if properties.DateFormatPattern != DefaultDateFormat || properties.LevelPattern != DefaultLevelPattern {
		t.Fatalf("pattern defaults = %q, %q", properties.DateFormatPattern, properties.LevelPattern)
	}
	if properties.MaxFileSize != 10_000_000 || properties.MaxHistory != 7 {
		t.Fatalf("rolling defaults = %d, %d", properties.MaxFileSize, properties.MaxHistory)
	}
	if properties.ConsoleThreshold == nil || *properties.ConsoleThreshold != slog.Level(-8) ||
		properties.FileThreshold == nil || *properties.FileThreshold != slog.Level(-8) {
		t.Fatalf("threshold defaults = %v, %v", properties.ConsoleThreshold, properties.FileThreshold)
	}
}

func TestRead_whenCurrentAndLegacyFilePropertiesExist_shouldPreferCurrentNames(t *testing.T) {
	environment := newEnvironment(t, map[string]any{
		FileName:       "current.log",
		LegacyFileName: "legacy.log",
		FilePath:       "current",
		LegacyFilePath: "legacy",
	})
	properties, err := Read(environment)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	if properties.FileName != "current.log" || properties.FilePath != "current" {
		t.Fatalf("file properties = %q, %q", properties.FileName, properties.FilePath)
	}
}

func TestRead_whenDataSizesUseSupportedUnits_shouldParseDecimalAndBinaryValues(t *testing.T) {
	environment := newEnvironment(t, map[string]any{
		FileMaxSize:      "1M",
		FileTotalSizeCap: "2GiB",
	})
	properties, err := Read(environment)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	if properties.MaxFileSize != 1_000_000 || properties.TotalSizeCap != 2*1024*1024*1024 {
		t.Fatalf("data sizes = %d, %d", properties.MaxFileSize, properties.TotalSizeCap)
	}
}

func TestRead_whenBuiltInGroupIsConfigured_shouldExpandGoarkLoggers(t *testing.T) {
	environment := newEnvironment(t, map[string]any{LevelPrefix + "web": "debug"})
	properties, err := Read(environment)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	if properties.LoggerLevels["goark.dev.arkhos"] != slog.LevelDebug {
		t.Fatalf("web group levels = %#v", properties.LoggerLevels)
	}
}

func TestRead_whenUserOverridesBuiltInGroup_shouldUseUserMembers(t *testing.T) {
	environment := newEnvironment(t, map[string]any{
		GroupPrefix + "web": "example.http",
		LevelPrefix + "web": "warn",
	})
	properties, err := Read(environment)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	if properties.LoggerLevels["example.http"] != slog.LevelWarn {
		t.Fatalf("web group levels = %#v", properties.LoggerLevels)
	}
	if _, exists := properties.LoggerLevels["goark.dev.arkhos"]; exists {
		t.Fatal("overridden built-in group member should not remain")
	}
}

func TestRead_whenDirectLevelTargetsGroupMember_shouldOverrideExpandedLevel(t *testing.T) {
	environment := newEnvironment(t, map[string]any{
		LevelPrefix + "web":              "debug",
		LevelPrefix + "goark.dev.arkhos": "error",
	})
	properties, err := Read(environment)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	if properties.LoggerLevels["goark.dev.arkhos"] != slog.LevelError {
		t.Fatalf("direct logger level = %s", properties.LoggerLevels["goark.dev.arkhos"])
	}
}

func TestRead_whenValuesAreInvalid_shouldFail(t *testing.T) {
	tests := map[string]map[string]any{
		"boolean":            {ConsoleEnabled: "sometimes"},
		"integer":            {FileMaxHistory: "-1"},
		"data size":          {FileMaxSize: "large"},
		"structured format":  {StructuredConsoleFormat: "unknown"},
		"stacktrace root":    {StacktraceRoot: "middle"},
		"stacktrace printer": {StacktracePrinter: "custom.Type"},
		"Java customizer":    {StructuredJSONCustomizer: "com.example.LoggingCustomizer"},
	}
	for name, values := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := Read(newEnvironment(t, values)); err == nil {
				t.Fatal("Read() error = nil")
			}
		})
	}
}

func newEnvironment(t *testing.T, values map[string]any) coreenv.Environment {
	t.Helper()
	environment, err := coreenv.NewStandardEnvironment()
	if err != nil {
		t.Fatalf("NewStandardEnvironment() error = %v", err)
	}
	source, err := coreenv.NewMapPropertySource("test", values)
	if err != nil {
		t.Fatalf("NewMapPropertySource() error = %v", err)
	}
	if err := environment.PropertySources().AddFirst(source); err != nil {
		t.Fatalf("AddFirst() error = %v", err)
	}
	return environment
}
