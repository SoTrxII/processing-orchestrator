package utils

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// A job outlives every call that drives it : it waits in the state store
// between steps, for hours when it is replayed after a restart, and the
// workers answer on topics of their own. Dapr starts a new trace for each call
// that carries none, so a session was a dozen unrelated traces in Tempo. The
// job keeps the trace it was started in, and every call made for it carries
// that again : the sidecar reads `traceparent` from gRPC metadata, and puts it
// in the CloudEvents it publishes, so the workers' spans join it too.

// Traceparent is the W3C trace context of an incoming call. Dapr sets it on
// every call it hands the app (SpanContextToGRPCMetadata, dapr 1.18), and the
// `tracing` Configuration samples everything. Empty from any other caller
func Traceparent(ctx context.Context) string {
	md, _ := metadata.FromIncomingContext(ctx)
	if v := md.Get("traceparent"); len(v) > 0 {
		return v[0]
	}
	return ""
}

// TraceContext is the context of a call made for a job started in this
// trace. Without one (jobs stored before, or not started through Dapr) the
// sidecar starts a new trace, as before
func TraceContext(traceparent string) context.Context {
	if traceparent == "" {
		return context.Background()
	}
	return metadata.AppendToOutgoingContext(context.Background(), "traceparent", traceparent)
}
