package middleware

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

// ServerInstrumentInterceptor instruments gRPC requests for errors and latency.
func ServerInstrumentInterceptor(duration *prometheus.HistogramVec) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		begin := time.Now()
		resp, err := handler(ctx, req)
		status := "success"
		if err != nil {
			status = "error"
		}
		duration.WithLabelValues(gRPC, info.FullMethod, status, "false").Observe(time.Since(begin).Seconds())
		return resp, err
	}
}
