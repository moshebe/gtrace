package filter

import (
	"fmt"
	"regexp"
	"strings"
)

type Filter struct {
	pattern *regexp.Regexp
	include bool
}

func New(patterns []string, include bool) (*Filter, error) {
	rePatterns := make([]string, len(patterns))
	for i := range patterns {
		rePatterns[i] = fmt.Sprintf("(?:%s)", patterns[i])
	}
	pattern, err := regexp.Compile(strings.Join(rePatterns, "|"))
	if err != nil {
		return nil, fmt.Errorf("compile regex: %w", err)
	}

	return &Filter{
		pattern: pattern,
		include: include,
	}, nil
}

func (f *Filter) Pass(value string) bool {
	if value == "" {
		return !f.include
	}
	match := f.pattern.MatchString(value)
	if f.include {
		return match
	} else {
		return !match
	}
}
