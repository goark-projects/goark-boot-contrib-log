package gbclog

import (
	"context"
	"log/slog"
	"sync"

	goarklog "goark.dev/log"
)

// Runtime 保存应用日志器及其受管生命周期。
type Runtime struct {
	context         *goarklog.LoggerContext
	logger          *slog.Logger
	previousDefault *slog.Logger
	installed       bool
	closeOnce       sync.Once
	closeErr        error
}

// Stop 在观测 Provider 停止前排空并关闭日志运行期。
func (r *Runtime) Stop(context.Context) error { return r.Close() }

// Logger 返回应用级结构化日志器。
func (r *Runtime) Logger() *slog.Logger {
	if r == nil || r.logger == nil {
		return slog.Default()
	}
	return r.logger
}

// Context 返回底层 goark-log 运行期；禁用模式返回 nil。
func (r *Runtime) Context() *goarklog.LoggerContext {
	if r == nil {
		return nil
	}
	return r.context
}

// Close 恢复默认 Logger，并幂等排空和关闭 goark-log。
func (r *Runtime) Close() error {
	if r == nil {
		return nil
	}
	r.closeOnce.Do(func() {
		if r.installed && slog.Default() == r.logger && r.previousDefault != nil {
			slog.SetDefault(r.previousDefault)
		}
		if r.context != nil {
			r.closeErr = r.context.Close()
		}
	})
	return r.closeErr
}

// Order 让日志在普通组件之后、观测 Provider 之前关闭。
func (*Runtime) Order() int { return -9000 }
