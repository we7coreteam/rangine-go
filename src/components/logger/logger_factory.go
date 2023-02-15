package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const DEFAULT_CHANNEL = "default"

type LoggerFactory struct {
	loggerMap map[string]*zap.Logger
}

func NewLoggerFactory() *LoggerFactory {
	return &LoggerFactory{
		loggerMap: make(map[string]*zap.Logger),
	}
}

func (loggerFactory *LoggerFactory) Channel(channel string) *zap.Logger {
	logger, exists := loggerFactory.loggerMap[channel]
	if !exists && channel == DEFAULT_CHANNEL {
		panic("logger channel " + channel + " not exists")
	}

	if !exists {
		return loggerFactory.Channel(DEFAULT_CHANNEL)
	}

	return logger
}

func (loggerFactory *LoggerFactory) MakeLogger(log Config) *zap.Logger {
	maxSize := log.MaxSize
	maxAge := log.MaxDays
	if maxSize <= 0 {
		maxSize = 2
	}
	if maxAge <= 0 {
		maxAge = 2
	}
	hook := lumberjack.Logger{
		Filename:   "./runtimes/logs/" + log.Path,
		MaxSize:    maxSize,
		MaxBackups: 2,
		MaxAge:     maxAge,
		Compress:   true,
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,    // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	switch log.Level {
	case "debug":
		atomicLevel.SetLevel(zap.DebugLevel)
	case "info":
		atomicLevel.SetLevel(zap.InfoLevel)
	case "warn":
		atomicLevel.SetLevel(zap.WarnLevel)
	case "error":
		atomicLevel.SetLevel(zap.ErrorLevel)
	case "fatal":
		atomicLevel.SetLevel(zap.FatalLevel)
	case "Panic":
		atomicLevel.SetLevel(zap.PanicLevel)
	default:
		atomicLevel.SetLevel(zap.DebugLevel)
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),               // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(&hook)), // 打印到控制台和文件
		atomicLevel, // 日志级别
	)

	return zap.New(core)
}

func (loggerFactory *LoggerFactory) Register(maps map[string]Config) {
	for key, value := range maps {
		loggerFactory.loggerMap[key] = loggerFactory.MakeLogger(value)
	}
}
