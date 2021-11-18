package logger

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Log interface {
	GetLogger() *zap.Logger
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
}

type Logger struct {
	logger *zap.Logger
}

type LoggConf struct {
	Level string
	File  string
}

func New(config LoggConf, projectRoot string, isDebug bool) *Logger {
	var zapLevel zapcore.Level
	zapLevel.Set(config.Level)

	var cfg zap.Config
	if isDebug {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}

	cfg.Level.SetLevel(zapLevel)

	cfg.OutputPaths = []string{
		"stdout",
		projectRoot + "/" + config.File,
	}

	logger, err := cfg.Build(zap.AddCallerSkip(1))
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
