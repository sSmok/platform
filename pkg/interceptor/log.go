package interceptor

import (
	"context"
	"time"

	"github.com/sSmok/platform_common/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Log - логирует запросы
func Log(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	now := time.Now()

	resp, err := handler(ctx, req)
	if err != nil {
		logger.Error(err.Error(), zap.String("method", info.FullMethod), zap.Any("req", req))
	}
	logger.Info(
		"request",
		zap.String("method", info.FullMethod),
		zap.Any("req", req),
		zap.Any("resp", resp),
		zap.Duration("duration", time.Since(now)),
	)

	return resp, err
}
