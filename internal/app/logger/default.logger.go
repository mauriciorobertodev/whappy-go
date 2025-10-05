package logger

import (
	"fmt"
	"log"
)

type DefaultLogger struct {
	prefix string
	level  Level
}

func NewDefaultLogger(prefix string, level Level) *DefaultLogger {
	fmt.Println("Creating DefaultLogger with level:", level)
	return &DefaultLogger{prefix: prefix, level: level}
}

func (l *DefaultLogger) Info(msg string, args ...any) {
	if l.level == LevelNone {
		fmt.Println("Log level is NONE, skipping log.")
		return
	}
	log.Printf("[INFO] %s: %s", l.prefix, fmt.Sprintf(msg, args...))
}

func (l *DefaultLogger) Error(msg string, args ...any) {
	if l.level == LevelNone {
		fmt.Println("Log level is NONE, skipping log.")
		return
	}
	log.Printf("[ERROR] %s: %s", l.prefix, fmt.Sprintf(msg, args...))
}

func (l *DefaultLogger) Debug(msg string, args ...any) {
	if l.level == LevelNone {
		fmt.Println("Log level is NONE, skipping log.")
		return
	}
	log.Printf("[DEBUG] %s: %s", l.prefix, fmt.Sprintf(msg, args...))
}

func (l *DefaultLogger) Warn(msg string, args ...any) {
	if l.level == LevelNone {
		fmt.Println("Log level is NONE, skipping log.")
		return
	}
	log.Printf("[WARN] %s: %s", l.prefix, fmt.Sprintf(msg, args...))
}
