package gcp

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type GcpLogger2 struct {
	encoder *json.Encoder
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
	logger.logJson("DEBUG", template, args...)
}

func (logger *GcpLogger2) Info(template string, args ...any) {
	logger.logJson("INFO", template, args...)
}

func (logger *GcpLogger2) Error(template string, args ...any) {
	logger.logJson("ERROR", template, args...)
}

func (logger *GcpLogger2) logJson(severity string, template string, args ...any) {
	if logger.encoder == nil {
		logger.encoder = json.NewEncoder(os.Stderr)
	}

	labels := map[string]string{
		"appname": "flashcards",
	}

	entry := LogEntry{
		Severity:  severity,
		Timestamp: time.Now(),
		Message:   fmt.Sprintf(template, args...),
		Labels:    labels,
	}

	logger.encoder.Encode(entry)
}
