package logger

import (
	"log"
	"path"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logger *zap.Logger
}

type LoggConf struct {
	Level string
	File  string
}

func New(zapConfig zap.Config, config LoggConf, projectRoot string) *Logger {
	var zapLevel zapcore.Level
	zapLevel.Set(config.Level)
	zapConfig.Level.SetLevel(zapLevel)
	zapConfig.OutputPaths = []string{
		"stdout",
		path.Join(projectRoot, config.File),
	}
	logger, err := zapConfig.Build(zap.AddCallerSkip(1))
	if err != nil {
		log.Fatalf("can't initialize logger logger: %v", err)
	}
	defer logger.Sync() // flushes buffer, if any

	return &Logger{
		logger: logger,
	}
}

func (l Logger) GetLogger() *zap.Logger {
	return l.logger
}

func (l Logger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l Logger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

func (l Logger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

func (l Logger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}
