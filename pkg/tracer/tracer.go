package tracer

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

// Get retrieve tracer from a specific project by the tracer id.
func (t *Tracer) Get(ctx context.Context, projectID, traceID string) (*cloudtrace.Trace, error) {
	return t.client.GetTrace(ctx, &cloudtrace.GetTraceRequest{
		ProjectId: projectID,
		TraceId:   traceID,
	})
}

// MultiGet retrieve tracer from multiple projects by the tracer id and aggregate the spans.
func (t *Tracer) MultiGet(ctx context.Context, projects []string, traceIDs []string) (*cloudtrace.Trace, error) {
	result := &cloudtrace.Trace{
		TraceId:   strings.Join(traceIDs, "+"),
		ProjectId: strings.Join(projects, "+"),
	}
	for _, project := range projects {
		for _, traceID := range traceIDs {
			res, err := t.Get(ctx, project, traceID)
			if err == nil {
				result.Spans = append(result.Spans, res.Spans...)
			}
		}
	}
	if len(result.Spans) == 0 {
		return nil, fmt.Errorf("no spans found for trace %q", traceIDs)
	}
	return result, nil
}

// List returns of a list of traces that match the specified options conditions.
func (t *Tracer) List(ctx context.Context, projectID string, limit int32, opts ...ListOption) ([]*cloudtrace.Trace, error) {
	req := &cloudtrace.ListTracesRequest{
		ProjectId: projectID,
		View:      cloudtrace.ListTracesRequest_COMPLETE,
	}
	for _, opt := range opts {
		opt(req)
	}

	var traces []*cloudtrace.Trace
	it := t.client.ListTraces(ctx, req)
	var count int32 = 0
	for {
		trace, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		count++
		traces = append(traces, trace)
		if count >= limit {
			break
		}
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
