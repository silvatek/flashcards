package test

import (
	"context"
	"fmt"
	"strings"
)

type LogRecorder struct {
	Entries []string
}

func (r *LogRecorder) Debug(ctx context.Context, template string, args ...any) {
	r.capture("DEBUG", template, args)
}

func (r *LogRecorder) Info(ctx context.Context, template string, args ...any) {
	r.capture("INFO", template, args)
}

func (r *LogRecorder) Error(ctx context.Context, template string, args ...any) {
	r.capture("ERROR", template, args)
}

func (r *LogRecorder) capture(level string, template string, args ...any) {
	if r.Entries == nil {
		r.Entries = make([]string, 0)
	}
	r.Entries = append(r.Entries, fmt.Sprintf(level+" "+template, args))
}

func (r *LogRecorder) HasEntryWithPrefix(prefix string) bool {
	if r.Entries == nil {
		return false
	}
	for _, entry := range r.Entries {
		if strings.HasPrefix(entry, prefix) {
			return true
		}
	}
	return false
}
