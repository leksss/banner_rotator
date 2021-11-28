package interfaces

import "go.uber.org/zap"

type Log interface {
	GetLogger() *zap.Logger
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
}
