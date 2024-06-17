package platform

import (
	"log"
)

type Logger interface {
	Debug(template string, args ...any)
	Info(template string, args ...any)
	Error(template string, args ...any)
}

type ConsoleLogger struct {
}

func (logger *ConsoleLogger) Debug(template string, args ...any) {
	log.Printf("DEBUG "+template, args...)
}

func (logger *ConsoleLogger) Info(template string, args ...any) {
	log.Printf("INFO  "+template, args...)
}

func (logger *ConsoleLogger) Error(template string, args ...any) {
	log.Printf("ERROR "+template, args...)
}
