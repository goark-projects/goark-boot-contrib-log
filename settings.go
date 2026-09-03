package gbclog

import (
	"context"
	"fmt"

	coreenv "goark.dev/goark/core/env"
	goarklog "goark.dev/log"
)

func newSettings(environment coreenv.Environment, options []Option) (settings, error) {
	resolved := settings{factory: defaultLoggerContextFactory}
	if environment != nil {
		enabled, err := coreenv.ResolveValueAs[bool](environment, "${"+PropertyEnabled+":true}")
		if err != nil {
			return settings{}, err
		}
		installDefault, err := coreenv.ResolveValueAs[bool](environment, "${"+PropertyInstallDefault+":true}")
		if err != nil {
			return settings{}, err
		}
		resolved.enabled = &enabled
		resolved.installDefault = &installDefault
	}
	for _, option := range options {
		if option != nil {
			option(&resolved)
		}
	}
	if resolved.enabled == nil {
		value := DefaultEnabled
		resolved.enabled = &value
	}
	if resolved.installDefault == nil {
		value := DefaultInstallDefault
		resolved.installDefault = &value
	}
	if resolved.factory == nil {
		return settings{}, fmt.Errorf("gbc-log: logger context factory is nil")
	}
	return resolved, nil
}

func defaultLoggerContextFactory(ctx context.Context, environment coreenv.Environment) (*goarklog.LoggerContext, error) {
	returnContext, _, err := goarklog.NewConfiguredLoggerContext(ctx, goarklog.WithBootPropertyResolver(environment))
	return returnContext, err
}
