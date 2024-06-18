package platform

import (
	"context"
	"log"
	"net/http"
	"strings"
)

type Logger interface {
	Debug(ctx context.Context, template string, args ...any)
	Info(ctx context.Context, template string, args ...any)
	Error(ctx context.Context, template string, args ...any)
}

type HttpRequestContext string

const HttpRequestKey HttpRequestContext = "HttpRequestContext"

type ConsoleLogger struct {
}

func HttpRequestFromContext(ctx context.Context) *http.Request {
	req := ctx.Value(HttpRequestKey)
	if req == nil {
		return nil
	}
	req2, ok := req.(*http.Request)
	if ok {
		return req2
	}
	return nil
}

func ParseCloudTraceHeader(header []string) (string, string, string) {
	if len(header) > 0 {
		return ParseCloudTrace(header[0])
	} else {
		return "", "", ""
	}
}

func ParseCloudTrace(trace string) (string, string, string) {
	if strings.Contains(trace, "/") {
		parts := strings.Split(trace, "/")

		if len(parts) >= 2 {
			if strings.Contains(parts[1], ";") {
				spanParts := strings.Split(parts[1], ";")
				return parts[0], spanParts[0], spanParts[1]
			} else {
				return parts[0], parts[1], ""
			}
		}
	}
	return "", "", ""
}

func (logger *ConsoleLogger) Debug(ctx context.Context, template string, args ...any) {
	log.Printf("DEBUG "+template, args...)
}

func (logger *ConsoleLogger) Info(ctx context.Context, template string, args ...any) {
	log.Printf("INFO  "+template, args...)
}

func (logger *ConsoleLogger) Error(ctx context.Context, template string, args ...any) {
	log.Printf("ERROR "+template, args...)
}
