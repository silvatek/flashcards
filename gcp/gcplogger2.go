package gcp

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"flashcards/platform"
)

type GcpLogger2 struct {
	encoder *json.Encoder
	project string
}

type LogEntry struct {
	Severity    string            `json:"severity"`
	Timestamp   time.Time         `json:"timestamp"`
	Message     interface{}       `json:"message,omitempty"`
	TextPayload interface{}       `json:"textPayload,omitempty"`
	Labels      map[string]string `json:"logging.googleapis.com/labels,omitempty"`
	TraceID     string            `json:"logging.googleapis.com/trace,omitempty"`
	SpanID      string            `json:"logging.googleapis.com/spanId,omitempty"`
	HttpRequest HttpRequestLog    `json:"httpRequest,omitempty"`
}

type HttpRequestLog struct {
	RequestMethod string `json:"requestMethod,omitempty"`
	RequestUrl    string `json:"requestUrl,omitempty"`
}

func (logger *GcpLogger2) Debug(template string, args ...any) {
	logger.DebugCtx(context.Background(), template, args...)
}

func (logger *GcpLogger2) DebugCtx(ctx context.Context, template string, args ...any) {
	logger.logJson(ctx, "DEBUG", template, args...)
}

func (logger *GcpLogger2) Info(template string, args ...any) {
	logger.logJson(context.Background(), "INFO", template, args...)
}

func (logger *GcpLogger2) Error(template string, args ...any) {
	logger.logJson(context.Background(), "ERROR", template, args...)
}

func (logger *GcpLogger2) logJson(ctx context.Context, severity string, template string, args ...any) {
	if logger.encoder == nil {
		logger.encoder = json.NewEncoder(os.Stderr)
	}
	if logger.project == "" {
		logger.project = os.Getenv("GCLOUD_PROJECT")
	}

	labels := map[string]string{
		"appname": "flashcards",
	}

	entry := LogEntry{
		Severity:  severity,
		Timestamp: time.Now(),
		Message:   fmt.Sprintf(template, args...),
	}

	req := platform.HttpRequestFromContext(ctx)
	if req != nil {
		entry.HttpRequest = HttpRequestLog{RequestMethod: req.Method, RequestUrl: req.RequestURI}

		traceID, spanID, _ := platform.ParseCloudTraceHeader(req.Header["X-Cloud-Trace-Context"])
		if traceID != "" {
			entry.TraceID = fmt.Sprintf("projects/%s/traces/%s", logger.project, traceID)
			entry.SpanID = spanID
		}
	}

	entry.Labels = labels

	logger.encoder.Encode(entry)
}
