package execctx

import (
	"context"
	"github.com/krls256/dsd2024additional/pkg/repositories"
	"go.uber.org/zap"
)

type ContextFactory struct {
	executorFactory *ExecutorFactory
}

func NewContextFactory(executorFactory *ExecutorFactory) *ContextFactory {
	return &ContextFactory{
		executorFactory: executorFactory,
	}
}

func (f *ContextFactory) Wrap(ctx context.Context) ExecutorContext {
	return WrapContext(ctx, f.executorFactory.Executor())
}

func WrapContext(ctx context.Context, executor *Executor) ExecutorContext {
	return ExecutorContext{
		Context: ctx,

		executor: executor,
	}
}

type ExecutorContext struct {
	context.Context

	executor *Executor
}

func (ctx *ExecutorContext) Add(fn func(tx repositories.TrWrapper) error) {
	ctx.executor.Add(fn)
}

func (ctx *ExecutorContext) Exec() error {
	ok, err := ctx.executor.Exec(ctx)
	if err != nil {
		return err
	}

	if !ok {
		zap.S().Error("critical error: executed twice")
	}

	return nil
}
