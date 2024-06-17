package platform

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/logging"
)

const LOCAL_LOGS = 0
const GCLOUD_LOGS = 1

type Logger interface {
	Debug(template string, args ...any)
	Info(template string, args ...any)
	Error(template string, args ...any)
}

type ConsoleLogger struct {
	mode    int
	project string
	client  *logging.Client
	logs    *logging.Logger
}

func (logger *ConsoleLogger) Init() {
	if runningOnGCloud() {
		logger.mode = GCLOUD_LOGS
		logger.project = "flashcards-425408"
		client, err := logging.NewClient(context.Background(), logger.project)
		if err == nil {
			logger.client = client
			logger.logs = client.Logger("flashcards")
		}
	} else {
		logger.mode = LOCAL_LOGS
	}
}

func runningOnGCloud() bool {
	projectId := os.Getenv("GCLOUD_PROJECT")
	return len(projectId) > 0
}

func (logger *ConsoleLogger) Debug(template string, args ...any) {
	logger.debug1(context.TODO(), template, args...)
}

func (logger *ConsoleLogger) debug1(ctx context.Context, template string, args ...any) {
	switch logger.mode {
	case GCLOUD_LOGS:
		logger.gCloudLog(ctx, logging.Debug, template, args...)
	default:
		log.Printf("DEBUG "+template, args...)
	}
}

func (logger *ConsoleLogger) Info(template string, args ...any) {
	logger.info1(context.TODO(), template, args...)
}

func (logger *ConsoleLogger) info1(ctx context.Context, template string, args ...any) {
	switch logger.mode {
	case GCLOUD_LOGS:
		logger.gCloudLog(ctx, logging.Info, template, args...)
	default:
		log.Printf("INFO  "+template, args...)
	}
}

func (logger *ConsoleLogger) Error(template string, args ...any) {
	logger.error1(context.TODO(), template, args...)
}

func (logger *ConsoleLogger) error1(ctx context.Context, template string, args ...any) {
	switch logger.mode {
	case GCLOUD_LOGS:
		logger.gCloudLog(ctx, logging.Error, template, args...)
	default:
		log.Printf("ERROR "+template, args...)
	}
}

func (logger *ConsoleLogger) gCloudLog(ctx context.Context, severity logging.Severity, template string, args ...any) {
	labels := make(map[string]string)
	logger.logs.Log(logging.Entry{
		Payload:  fmt.Sprintf(template, args...),
		Severity: severity,
		Labels:   labels,
	})
}