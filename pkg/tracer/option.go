package tracer

import (
	"strings"
	"time"

	"google.golang.org/genproto/googleapis/devtools/cloudtrace/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ListOption func(request *cloudtrace.ListTracesRequest)

func WithStartTime(start time.Time) ListOption {
	return func(r *cloudtrace.ListTracesRequest) {
		r.StartTime = timestamppb.New(start)
	}
}
func WithEndTime(end time.Time) ListOption {
	return func(r *cloudtrace.ListTracesRequest) {
		r.EndTime = timestamppb.New(end)
	}
}
func WithSince(since time.Duration) ListOption {
	return func(r *cloudtrace.ListTracesRequest) {
		now := time.Now()
		r.StartTime = timestamppb.New(now.Add(-since))
		r.EndTime = timestamppb.New(now)
	}
}

func WithFilter(filters ...string) ListOption {
	return func(r *cloudtrace.ListTracesRequest) {
		if len(filters) <= 0 {
			return
		}
		if r.Filter != "" {
			r.Filter += "&"
		}

		r.Filter += strings.Join(filters, "&")
	}
}

func WithLimit(limit int32) ListOption {
	return func(r *cloudtrace.ListTracesRequest) {
		r.PageSize = limit
	}
}

func WithOrderBy(field string, desc bool) ListOption {
	return func(r *cloudtrace.ListTracesRequest) {
		r.OrderBy = field
		if desc {
			r.OrderBy += " desc"
		}
	}
}
func WithOrderByDuration(desc bool) ListOption {
	return WithOrderBy("duration", desc)
}
func WithOrderByStart(desc bool) ListOption {
	return WithOrderBy("start", desc)
}
func WithOrderByName(desc bool) ListOption {
	return WithOrderBy("name", desc)
}

func WithOnlyRootSpanView() ListOption {
	return func(r *cloudtrace.ListTracesRequest) {
		r.View = cloudtrace.ListTracesRequest_ROOTSPAN
	}
}
