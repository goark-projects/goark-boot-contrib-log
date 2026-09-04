package gbclog

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	coreenv "goark.dev/goark/core/env"
	coreresource "goark.dev/goark/core/resource"
	goarklog "goark.dev/log"
)

func TestDefaultLoggerContextFactory_whenClasspathConfigUsed_shouldLoadAndResolveBootProperties(t *testing.T) {
	filesystem := fstest.MapFS{
		"logging/goark-log.yml": &fstest.MapFile{Data: []byte(`
appenders:
  console:
    type: console
root:
  level: ${logging.level.root}
  appenderRefs: [console]
`)},
	}
	loader, err := coreresource.NewLoader(coreresource.WithFS("classpath", filesystem))
	if err != nil {
		t.Fatalf("NewLoader() error = %v", err)
	}
	environment := newLoggingEnvironment(t, map[string]any{
		PropertyConfig:    "classpath:/logging/goark-log.yml",
		PropertyRootLevel: "error",
	})
	loggerContext, err := defaultLoggerContextFactory(context.Background(), environment, loader)
	if err != nil {
		t.Fatalf("defaultLoggerContextFactory() error = %v", err)
	}
	defer loggerContext.Close()
	result := loggerContext.ConfigResult()
	if result == nil || result.Path != "classpath:/logging/goark-log.yml" {
		t.Fatalf("ConfigResult() = %#v", result)
	}
	configurations := loggerContext.LoggerConfigurations()
	if len(configurations) == 0 || configurations[0].Name != "ROOT" || configurations[0].EffectiveLevel != slog.LevelError {
		t.Fatalf("LoggerConfigurations() = %#v", configurations)
	}
}

func TestConfigResourceOption_whenFileLocationUsed_shouldPassResolvedPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "goark-log.yml")
	if err := os.WriteFile(path, []byte(`
appenders:
  console:
    type: console
root:
  level: info
  appenderRefs: [console]
`), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	loader, err := coreresource.NewLoader()
	if err != nil {
		t.Fatalf("NewLoader() error = %v", err)
	}
	environment := newLoggingEnvironment(t, map[string]any{PropertyConfig: "file:" + filepath.ToSlash(path)})
	option, err := configResourceOption(context.Background(), environment, loader)
	if err != nil {
		t.Fatalf("configResourceOption() error = %v", err)
	}
	options, result, err := goarklog.LoadOptions(context.Background(), option)
	if err != nil {
		t.Fatalf("LoadOptions() error = %v", err)
	}
	defer closeAppenders(t, options.Appenders)
	if result.Path != path {
		t.Fatalf("config path = %q, want %q", result.Path, path)
	}
}

func TestConfigResourceOption_whenResourceMissing_shouldFail(t *testing.T) {
	loader, err := coreresource.NewLoader(coreresource.WithFS("classpath", fstest.MapFS{}))
	if err != nil {
		t.Fatalf("NewLoader() error = %v", err)
	}
	environment, err := coreenv.NewStandardEnvironment()
	if err != nil {
		t.Fatalf("NewStandardEnvironment() error = %v", err)
	}
	source, err := coreenv.NewMapPropertySource("test", map[string]any{PropertyConfig: "classpath:/missing.yml"})
	if err != nil {
		t.Fatalf("NewMapPropertySource() error = %v", err)
	}
	if err := environment.PropertySources().AddFirst(source); err != nil {
		t.Fatalf("AddFirst() error = %v", err)
	}
	if _, err := configResourceOption(context.Background(), environment, loader); err == nil {
		t.Fatal("missing logging.config should fail")
	}
}
