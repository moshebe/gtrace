package span

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"google.golang.org/genproto/googleapis/devtools/cloudtrace/v1"
	"gtrace/filter"
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

func Summary(spans []*cloudtrace.TraceSpan, writer io.Writer) error {
	var buf bytes.Buffer
	for _, s := range spans {
		duration := s.EndTime.AsTime().Sub(s.StartTime.AsTime())
		buf.WriteString(fmt.Sprintf("%s (%v - %s)\n", s.Name, s.StartTime.AsTime(), duration))
		if len(s.Labels) > 0 {
			labels, err := json.Marshal(s.Labels)
			if err != nil {
				return err
			}
			buf.WriteString(fmt.Sprintf("\t%s\n", labels))
		}

		if _, err := writer.Write(buf.Bytes()); err != nil {
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
