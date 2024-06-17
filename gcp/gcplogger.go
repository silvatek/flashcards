package gcp

import (
	"context"
	"fmt"

	"cloud.google.com/go/logging"
)

type GcpLogger struct {
	project string
	client  *logging.Client
	logs    *logging.Logger
}

func (logger *GcpLogger) init() {
	logger.project = "flashcards-425408"
	client, err := logging.NewClient(context.Background(), logger.project)
	if err == nil {
		logger.client = client
		logger.logs = client.Logger("flashcards")
	}
}

func (logger *GcpLogger) Debug(template string, args ...any) {
	logger.log(logging.Debug, template, args...)
}

func (logger *GcpLogger) Info(template string, args ...any) {
	logger.log(logging.Info, template, args...)
}

func (logger *GcpLogger) Error(template string, args ...any) {
	logger.log(logging.Error, template, args...)
}

func (logger *GcpLogger) log(severity logging.Severity, template string, args ...any) {
	labels := make(map[string]string)
	logger.logs.Log(logging.Entry{
		Payload:  fmt.Sprintf(template, args...),
		Severity: severity,
		Labels:   labels,
	})
}
