package gbclog

import (
	"context"
	"fmt"
	"strings"

	coreenv "goark.dev/goark/core/env"
	coreresource "goark.dev/goark/core/resource"
	goarklog "goark.dev/log"
)

func newSettings(environment coreenv.Environment, options []Option) (settings, error) {
	resolved := settings{}
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
		loader, err := defaultResourceLoader(resolved)
		if err != nil {
			return settings{}, err
		}
		resolved.factory = func(ctx context.Context, environment coreenv.Environment) (*goarklog.LoggerContext, error) {
			return defaultLoggerContextFactory(ctx, environment, loader)
		}
	}
	return resolved, nil
}

func defaultLoggerContextFactory(ctx context.Context, environment coreenv.Environment, loader coreresource.Loader) (*goarklog.LoggerContext, error) {
	loadOptions := []goarklog.ConfigLoadOption{
		goarklog.WithBootPropertyResolver(environment),
		goarklog.WithOptionsCustomizer(loggingOptionsCustomizer(environment)),
	}
	resourceOption, err := configResourceOption(ctx, environment, loader)
	if err != nil {
		return nil, err
	}
	if resourceOption != nil {
		loadOptions = append(loadOptions, resourceOption)
	}
	returnContext, _, err := goarklog.NewConfiguredLoggerContext(
		ctx,
		loadOptions...,
	)
	return returnContext, err
}

func defaultResourceLoader(configuration settings) (coreresource.Loader, error) {
	if configuration.resourceLoader != nil {
		return configuration.resourceLoader, nil
	}
	options := make([]coreresource.LoaderOption, 0, 1)
	if configuration.classpathFS != nil {
		options = append(options, coreresource.WithFS("classpath", configuration.classpathFS))
	}
	loader, err := coreresource.NewLoader(options...)
	if err != nil {
		return nil, fmt.Errorf("gbc-log: create resource loader: %w", err)
	}
	return loader, nil
}

func configResourceOption(ctx context.Context, environment coreenv.Environment, loader coreresource.Loader) (goarklog.ConfigLoadOption, error) {
	if environment == nil {
		return nil, nil
	}
	location, found := environment.GetProperty(PropertyConfig)
	location = strings.TrimSpace(location)
	if !found || location == "" {
		return nil, nil
	}
	loadLocation := location
	if strings.HasPrefix(location, "classpath:") {
		path := strings.TrimLeft(strings.TrimPrefix(location, "classpath:"), "/\\")
		if path == "" {
			return nil, fmt.Errorf("gbc-log: logging.config classpath is empty")
		}
		loadLocation = "fs:classpath/" + strings.ReplaceAll(path, "\\", "/")
	}
	resource, err := loader.Load(loadLocation)
	if err != nil {
		return nil, fmt.Errorf("gbc-log: load logging.config %q: %w", location, err)
	}
	exists, err := resource.Exists(ctx)
	if err != nil {
		return nil, fmt.Errorf("gbc-log: inspect logging.config %q: %w", location, err)
	}
	if !exists {
		return nil, fmt.Errorf("gbc-log: logging.config %q does not exist", location)
	}
	if file, ok := resource.(*coreresource.FileResource); ok {
		return goarklog.WithConfigPath(file.Path()), nil
	}
	data, err := resource.ReadAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("gbc-log: read logging.config %q: %w", location, err)
	}
	return goarklog.WithConfigData(location, data), nil
}
