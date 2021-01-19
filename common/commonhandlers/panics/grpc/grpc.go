package grpc

import (
	"context"
	"github.com/ahmadrezamusthafa/logwatcher/common/commonhandlers/panics/internal"
	"github.com/ahmadrezamusthafa/logwatcher/common/logger"
	"google.golang.org/grpc"
	"os"
	"runtime/debug"
)

func PanicInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	defer internal.HandlePanic(func(err error) {
		logger.Panic("%+v", err)
		stack := internal.TrimStackTrace(debug.Stack())
		os.Stderr.Write(stack)
		internal.PublishError(err, stack, nil)
	})
	return handler(ctx, req)
}