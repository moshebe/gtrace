package trace

import (
	"context"
	"fmt"
	"strings"

	traceapi "cloud.google.com/go/trace/apiv1"
	"google.golang.org/api/iterator"
	"google.golang.org/genproto/googleapis/devtools/cloudtrace/v1"
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
func (t *Tracer) Get(ctx context.Context, projectID, traceID string) (*cloudtrace.Trace, error) {
	return t.client.GetTrace(ctx, &cloudtrace.GetTraceRequest{
		ProjectId: projectID,
		TraceId:   traceID,
	})
}

// MultiGet retrieve trace from multiple projects by the trace id and aggregate the spans.
func (t *Tracer) MultiGet(ctx context.Context, projects []string, traceID string) (*cloudtrace.Trace, error) {
	result := &cloudtrace.Trace{
		TraceId:   traceID,
		ProjectId: strings.Join(projects, "+"),
	}
	for _, project := range projects {
		res, err := t.Get(ctx, project, traceID)
		if err == nil {
			result.Spans = append(result.Spans, res.Spans...)
		}
	}
	if len(result.Spans) == 0 {
		return nil, fmt.Errorf("trace not found")
	}
	return result, nil
}

// List returns of a list of traces that match the specified options conditions.
func (t *Tracer) List(ctx context.Context, projectID string, opts ...ListOption) ([]*cloudtrace.Trace, error) {
	req := &cloudtrace.ListTracesRequest{
		ProjectId: projectID,
		View:      cloudtrace.ListTracesRequest_COMPLETE,
	}
	for _, opt := range opts {
		opt(req)
	}

	var traces []*cloudtrace.Trace
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

// Close closes the inner client connection to the API service.
func (t *Tracer) Close() error {
	if t.client == nil {
		return nil
	}
	return t.client.Close()
}
