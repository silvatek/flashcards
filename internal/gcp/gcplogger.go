package gcp

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"flashcards/internal/platform"
)

type GcpLogger struct {
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

const TRACE_HEADER_NAME = "X-Cloud-Trace-Context"

func (logger *GcpLogger) Debug(ctx context.Context, template string, args ...any) {
	logger.logJson(ctx, "DEBUG", template, args...)
}

func (logger *GcpLogger) Info(ctx context.Context, template string, args ...any) {
	logger.logJson(ctx, "INFO", template, args...)
}

func (logger *GcpLogger) Error(ctx context.Context, template string, args ...any) {
	logger.logJson(ctx, "ERROR", template, args...)
}

func (logger *GcpLogger) logJson(ctx context.Context, severity string, template string, args ...any) {
	if logger.encoder == nil {
		logger.encoder = json.NewEncoder(os.Stderr)
	}
	if logger.project == "" {
		logger.project = os.Getenv("GCLOUD_PROJECT")
	}

	entry := LogEntry{
		Severity:  severity,
		Timestamp: time.Now(),
		Message:   fmt.Sprintf(template, args...),
	}

	logger.addRequestDetails(&entry, ctx)
	logger.addStartupDetails(&entry, ctx)

	entry.Labels = map[string]string{
		"appname": "flashcards",
	}

	logger.encoder.Encode(entry)
}

func (logger *GcpLogger) addRequestDetails(entry *LogEntry, ctx context.Context) {
	req := platform.HttpRequestFromContext(ctx)
	if req != nil {
		entry.HttpRequest = HttpRequestLog{RequestMethod: req.Method, RequestUrl: req.RequestURI}

		traceID, spanID, _ := platform.ParseCloudTraceHeader(req.Header[TRACE_HEADER_NAME])
		if traceID != "" {
			entry.TraceID = fmt.Sprintf("projects/%s/traces/%s", logger.project, traceID)
			entry.SpanID = spanID
		}
	}
}

func (logger *GcpLogger) addStartupDetails(entry *LogEntry, ctx context.Context) {
	value := ctx.Value(platform.StartupContextKey)
	if value == nil {
		return
	}
	data, ok := value.(platform.StartupContextData)
	if ok {
		entry.TraceID = data.TraceID
		entry.SpanID = data.SpanID
	}
}
