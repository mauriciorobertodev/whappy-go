package logger

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/mauriciorobertodev/whappy-go/internal/app/logger"
)

type CharmLogger struct {
	*log.Logger
	level logger.Level
}

func NewCuteLogger(prefix string, level logger.Level) logger.Logger {
	return &CharmLogger{
		level: level,
		Logger: log.NewWithOptions(os.Stderr, log.Options{
			Level:           appLevelToCharm(level),
			Prefix:          prefix,
			ReportTimestamp: true,
			TimeFormat:      "01/02 03:04PM",
		}),
	}
}

func (l *CharmLogger) Info(msg string, args ...any) {
	if l.level != logger.LevelNone {
		l.Logger.Info(msg, args...)
	}
}

func (l *CharmLogger) Error(msg string, args ...any) {
	if l.level != logger.LevelNone {
		l.Logger.Error(msg, args...)
	}
}

func (l *CharmLogger) Debug(msg string, args ...any) {
	if l.level != logger.LevelNone {
		l.Logger.Debug(msg, args...)
	}
}

func (l *CharmLogger) Warn(msg string, args ...any) {
	if l.level != logger.LevelNone {
		l.Logger.Warn(msg, args...)
	}
}

func appLevelToCharm(level logger.Level) log.Level {
	switch level {
	case logger.LevelDebug:
		return log.DebugLevel
	case logger.LevelInfo:
		return log.InfoLevel
	case logger.LevelWarn:
		return log.WarnLevel
	case logger.LevelError:
		return log.ErrorLevel
	default:
		return log.InfoLevel
	}
}
