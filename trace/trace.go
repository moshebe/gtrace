package trace

import (
	"context"
	"fmt"

	traceapi "cloud.google.com/go/trace/apiv1"
	"google.golang.org/api/iterator"
	tracev1 "google.golang.org/genproto/googleapis/devtools/cloudtrace/v1"
)

type Tracer struct {
	client *traceapi.Client
}

func NewTracer(ctx context.Context) (*Tracer, error) {
	client, err := traceapi.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("new client: %w", err)
	}

	return &Tracer{client: client}, nil
}

// Get retrieve trace from a specific project by the trace id.
func (t *Tracer) Get(ctx context.Context, projectID, traceID string) (*tracev1.Trace, error) {
	return t.client.GetTrace(ctx, &tracev1.GetTraceRequest{
		ProjectId: projectID,
		TraceId:   traceID,
	})
}

// List returns of a list of traces that match the specified options conditions.
func (t *Tracer) List(ctx context.Context, projectID string, opts ...ListOption) ([]*tracev1.Trace, error) {
	req := &tracev1.ListTracesRequest{
		ProjectId: projectID,
		View:      tracev1.ListTracesRequest_COMPLETE,
	}
	for _, opt := range opts {
		opt(req)
	}

	var traces []*tracev1.Trace
	it := t.client.ListTraces(ctx, req)
	for {
		trace, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		traces = append(traces, trace)
	}

	return traces, nil
}
