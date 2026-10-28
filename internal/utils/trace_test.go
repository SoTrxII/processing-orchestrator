package utils

import (
	"context"
	"testing"

	"google.golang.org/grpc/metadata"
)

func TestTraceparent_TheCallersIsKept(t *testing.T) {
	tp := "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01"
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("traceparent", tp))
	if got := Traceparent(ctx); got != tp {
		t.Fatalf("got %q, want the caller's %q", got, tp)
	}
	if got := Traceparent(context.Background()); got != "" {
		t.Fatalf("no caller trace, got %q", got)
	}
}

// A job stored without a trace must not send an empty one
func TestTraceContext_NothingWithoutATrace(t *testing.T) {
	if md, ok := metadata.FromOutgoingContext(TraceContext("")); ok && len(md.Get("traceparent")) > 0 {
		t.Fatal("an empty traceparent was sent")
	}
}
