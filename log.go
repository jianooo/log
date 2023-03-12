package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

//zap记录器的实例
var logger *zap.Logger

type mode int

const (
	production mode = iota
	development
)

func NewProduction() {
	newLogger(production)
}
func NewDevelopment() {
	newLogger(development)
}

func newLogger(model mode) {

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "file",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalColorLevelEncoder, // 大小写编码器
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder, // 时间到秒
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.InfoLevel)

	if model == production {
		encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
		atomicLevel.SetLevel(zap.ErrorLevel)
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig), // 编码器配置
			zapcore.NewMultiWriteSyncer(),         // 打印到文件
			atomicLevel,                           // 日志级别
		)
		logger = zap.New(core)
	} else {
		// 开启开发模式，堆栈跟踪
		//caller := zap.AddCaller()
		// 开启文件及行号
		//development := zap.Development()
		core := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),                // 编码器配置
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), // 打印到控制台和文件
			atomicLevel, // 日志级别
		)
		logger = zap.New(core)
	}

}

func Debug(msg string, args ...zap.Field) {
	logger.Debug(msg, args...)
}

func Info(msg string, args ...zap.Field) {
	logger.Info(msg, args...)
}

func Warn(msg string, args ...zap.Field) {
	logger.Warn(msg, args...)
}

func Error(msg string, args ...zap.Field) {
	logger.Error(msg, args...)
}

func DPanic(msg string, args ...zap.Field) {
	logger.DPanic(msg, args...)
}

func Panic(msg string, fields ...zap.Field) {
	logger.Panic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}

func Sugar() *zap.SugaredLogger {
	return logger.Sugar()
}

func Named(s string) *zap.Logger {
	return logger.Named(s)
}

func Core() zapcore.Core {
	return logger.Core()
}

func Check(lvl zapcore.Level, msg string) *zapcore.CheckedEntry {
	return logger.Check(lvl, msg)
}

func Sync() error {
	return logger.Sync()
}
