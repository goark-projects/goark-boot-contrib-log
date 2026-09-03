package gbclog

import (
	"context"
	"log/slog"

	"goark.dev/boot"
	goarkcontainer "goark.dev/goark/container"
	appcontext "goark.dev/goark/context"
	goarklog "goark.dev/log"
)

// AutoConfigure 创建 goark-log 自动配置。
func AutoConfigure(options ...Option) boot.AutoConfiguration {
	copied := append([]Option(nil), options...)
	return boot.NewAutoConfiguration(StarterID, func(_ context.Context, app *appcontext.ApplicationContext) error {
		return app.RegisterConfiguration(configuration{options: copied})
	}, boot.WithAutoConfigurationOrder(-10000))
}

type configuration struct{ options []Option }

func (configuration) Name() string { return StarterID + ".configuration" }
func (configuration) Order() int   { return -20000 }
func (c configuration) Register(ctx context.Context, registry *goarkcontainer.Registry) error {
	return c.RegisterWithContext(ctx, appcontext.NewConfigurationContext(nil, registry))
}
func (c configuration) RegisterWithContext(_ context.Context, config appcontext.ConfigurationContext) error {
	resolved, err := newSettings(config.Environment(), c.options)
	if err != nil {
		return err
	}
	if err := goarkcontainer.Register[*Runtime](config.Registry(), BeanNameLifecycle, func(ctx context.Context, _ goarkcontainer.Resolver) (*Runtime, error) {
		previous := slog.Default()
		if !*resolved.enabled {
			return &Runtime{logger: previous}, nil
		}
		loggerContext, err := resolved.factory(ctx, config.Environment())
		if err != nil {
			return nil, err
		}
		logger := loggerContext.Logger("goark")
		runtime := &Runtime{context: loggerContext, logger: logger, previousDefault: previous}
		if *resolved.installDefault {
			slog.SetDefault(logger)
			runtime.installed = true
		}
		return runtime, nil
	}); err != nil {
		return err
	}
	if err := goarkcontainer.Register[*slog.Logger](config.Registry(), BeanNameLogger, func(ctx context.Context, resolver goarkcontainer.Resolver) (*slog.Logger, error) {
		runtime, err := goarkcontainer.Get[*Runtime](ctx, resolver, BeanNameLifecycle)
		if err != nil {
			return nil, err
		}
		return runtime.Logger(), nil
	}, goarkcontainer.WithDependsOn(BeanNameLifecycle), goarkcontainer.WithPrimary()); err != nil {
		return err
	}
	if !*resolved.enabled {
		return nil
	}
	return goarkcontainer.Register[*goarklog.LoggerContext](config.Registry(), BeanNameContext, func(ctx context.Context, resolver goarkcontainer.Resolver) (*goarklog.LoggerContext, error) {
		runtime, err := goarkcontainer.Get[*Runtime](ctx, resolver, BeanNameLifecycle)
		if err != nil {
			return nil, err
		}
		return runtime.Context(), nil
	}, goarkcontainer.WithDependsOn(BeanNameLifecycle))
}
