package logger

import "go.uber.org/zap"

// 这是一个没有意义的解耦
// 因为使用的时候同样要初始化这个Logger
// 所以依然耦合
var _ *zap.Logger

type Logger interface {
	Debug(msg string, args ...Field)
	Info(msg string, args ...Field)
	Warn(msg string, args ...Field)
	Error(msg string, args ...Field)
}

type Field struct {
	Key   string
	Value any
}

type ZapLogger struct {
	logger *zap.Logger
}

func NewLogger(logger *zap.Logger) Logger {
	return &ZapLogger{
		logger: logger,
	}
}

func (z *ZapLogger) Debug(msg string, args ...Field) {
	z.logger.Debug(msg, z.toArgs(args)...)
}

func (z *ZapLogger) Info(msg string, args ...Field) {
	z.logger.Info(msg, z.toArgs(args)...)
}

func (z *ZapLogger) Warn(msg string, args ...Field) {
	z.logger.Warn(msg, z.toArgs(args)...)
}

func (z *ZapLogger) Error(msg string, args ...Field) {
	z.logger.Error(msg, z.toArgs(args)...)
}

func (z *ZapLogger) toArgs(args []Field) (res []zap.Field) {
	res = make([]zap.Field, 0, len(args))
	for i := range args {
		res = append(res, zap.Any(args[i].Key, args[i].Value))
	}
	return res
}
