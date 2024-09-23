package logging

import (
	"context"

	"go.uber.org/zap"
)

const (
	FieldService  = "service"
	FieldEndpoint = "endpoint"
)

type LoggetCtxKey struct{}

func GetLoggerFromContext(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(LoggetCtxKey{}).(*zap.Logger); ok {
		return logger
	}
	return zap.NewNop()
}
