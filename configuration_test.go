package gbclog_test

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"

	"goark.dev/boot"
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
	if err := app.Close(t.Context()); err != nil {
		t.Fatalf("Close: %v", err)
	}
}
