package platform

import (
	"context"
	"log"
	"net/http"
)

type Logger interface {
	Debug(template string, args ...any)
	DebugCtx(ctx context.Context, template string, args ...any)
	Info(template string, args ...any)
	Error(template string, args ...any)
}

type HttpRequestContext string

const HttpRequestKey HttpRequestContext = "HttpRequestContext"

type ConsoleLogger struct {
}

func HttpRequestFromContext(ctx context.Context) *http.Request {
	req := ctx.Value(HttpRequestKey)
	if req != nil {
		return nil
	}
	req2, ok := req.(*http.Request)
	if ok {
		return req2
	}
	return nil
}

func (logger *ConsoleLogger) Debug(template string, args ...any) {
	log.Printf("DEBUG "+template, args...)
}

func (logger *ConsoleLogger) DebugCtx(ctx context.Context, template string, args ...any) {

	log.Printf("DEBUG "+template, args...)
}

func (logger *ConsoleLogger) Info(template string, args ...any) {
	log.Printf("INFO  "+template, args...)
}

func (logger *ConsoleLogger) Error(template string, args ...any) {
	log.Printf("ERROR "+template, args...)
}
