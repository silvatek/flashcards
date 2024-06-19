package gcp

import (
	"context"
	"net/http"
	"testing"

	"flashcards/internal/platform"
)

func TestAddRequestDetails(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add(TRACE_HEADER_NAME, "10000/200;o=1")
	ctx := context.WithValue(context.Background(), platform.HttpRequestKey, req)

	logs := GcpLogger{project: "TEST"}

	entry := LogEntry{}

	logs.addRequestDetails(&entry, ctx)

	if entry.HttpRequest.RequestMethod != http.MethodGet {
		t.Errorf("Unexpected LogEntry http request method: %s", entry.HttpRequest.RequestMethod)
	}
	if entry.TraceID != "projects/TEST/traces/10000" {
		t.Errorf("Unexpected LogEntry trace ID: %s", entry.TraceID)
	}
	if entry.SpanID != "200" {
		t.Errorf("Unexpected LogEntry span ID: %s", entry.SpanID)
	}
}
