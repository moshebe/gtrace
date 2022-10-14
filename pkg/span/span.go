package span

import (
	"errors"
	"fmt"
	"io"
	"sort"
	"text/template"
	"time"

	"github.com/moshebe/gtrace/pkg/filter"
	"google.golang.org/genproto/googleapis/devtools/cloudtrace/v1"
)

func Names(spans []*cloudtrace.TraceSpan) []string {
	names := make(map[string]struct{})

	for _, s := range spans {
		if _, found := names[s.Name]; found {
			continue
		}
		names[s.Name] = struct{}{}
	}

	results := make([]string, 0, len(names))
	for name := range names {
		results = append(results, name)
	}
	return results
}

func Filter(spans []*cloudtrace.TraceSpan, f *filter.Filter) []*cloudtrace.TraceSpan {
	result := make([]*cloudtrace.TraceSpan, 0, len(spans))
	for _, s := range spans {
		if !f.Pass(s.Name) {
			continue
		}
		result = append(result, s)
	}

	return result
}

func Sort(spans []*cloudtrace.TraceSpan) {
	sort.Slice(spans, func(i, j int) bool {
		return spans[i].StartTime.AsTime().Before(spans[j].StartTime.AsTime())
	})
}

func FilterLabels(labels map[string]string, f *filter.Filter) map[string]string {
	result := make(map[string]string, len(labels))
	for name, value := range labels {
		if !f.Pass(name) {
			continue
		}
		result[name] = value
	}

	return result
}

func Format(spans []*cloudtrace.TraceSpan, format string, writer io.Writer) error {
	type ExtSpan struct {
		*cloudtrace.TraceSpan
		Duration   time.Duration
		Start, End time.Time
	}

	t, err := template.New("").Parse(format)
	if err != nil {
		return fmt.Errorf("parse format: %w", err)
	}

	for _, s := range spans {
		err = t.Execute(writer, ExtSpan{
			TraceSpan: s,
			Start:     s.StartTime.AsTime(),
			End:       s.EndTime.AsTime(),
			Duration:  s.EndTime.AsTime().Sub(s.StartTime.AsTime()),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func ListRootSpans(traces []*cloudtrace.Trace) map[string][]string {
	results := make(map[string][]string, len(traces))
	for _, t := range traces {
		if len(t.Spans) == 0 {
			continue
		}
		name := t.Spans[0].Name
		results[name] = append(results[name], t.TraceId)
	}
	return results
}

func DurationSummary(span *cloudtrace.TraceSpan) string {
	return fmt.Sprintf("%s (%d) - took %s",
		span.GetName(),
		span.GetSpanId(),
		Duration(span))
}

func Duration(span *cloudtrace.TraceSpan) time.Duration {
	return span.GetEndTime().AsTime().Sub(span.GetStartTime().AsTime())
}

func FilterMinDuration(spans []*cloudtrace.TraceSpan, threshold time.Duration) []*cloudtrace.TraceSpan {
	results := make([]*cloudtrace.TraceSpan, 0, len(spans))
	for _, span := range spans {
		if !span.GetEndTime().IsValid() || !span.GetStartTime().IsValid() {
			continue
		}
		duration := span.GetEndTime().AsTime().Sub(span.GetStartTime().AsTime())
		if duration < threshold {
			continue
		}
		results = append(results, span)
	}
	return results
}

func SubTree(spans []*cloudtrace.TraceSpan, rootID uint64) ([]*cloudtrace.TraceSpan, error) {
	parents := map[uint64]struct{}{rootID: {}}
	results := make([]*cloudtrace.TraceSpan, 0, len(spans))
	for _, span := range spans {
		if span.GetSpanId() == rootID {
			results = append(results, span)
			continue
		}
		_, ok := parents[span.GetParentSpanId()]
		if !ok {
			continue
		}
		parents[span.GetSpanId()] = struct{}{}
		results = append(results, span)
	}
	if len(results) == 0 {
		return nil, errors.New("root span was not found")
	}
	return results, nil
}
