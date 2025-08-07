package logger

import (
	"os"
	"path/filepath"
	"time"

	"github.com/usual2970/retell/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger 封装zap.Logger
type Logger struct {
	log *zap.Logger
}

var defaultLogger *Logger

// Init 初始化日志
func Init() error {
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zapcore.InfoLevel)

	// 配置编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    "func",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	conf := config.GetConfig()

	// 确保日志目录存在
	logPath := conf.LogPath
	if err := os.MkdirAll(filepath.Dir(logPath), 0744); err != nil {
		return err
	}
	// 生成日志文件路径，按日期命名
	logFile := filepath.Join(conf.LogPath, time.Now().Format("2006-01-02")+".log")
	// 配置输出
	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)

	// 配置日志文件切割
	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFile, // 日志文件路径
		MaxSize:    500,     // 每个日志文件最大尺寸，单位MB
		MaxBackups: 3,       // 保留旧文件的最大个数
		MaxAge:     28,      // 保留旧文件的最大天数
		Compress:   true,    // 是否压缩/归档旧文件
	})

	cores := []zapcore.Core{}

	if config.IsDevelopment() {
		// 开发环境输出到控制台
		cores = append(cores, zapcore.NewCore(jsonEncoder, zapcore.AddSync(os.Stdout), atomicLevel))
	} else {
		// 生产环境输出到文件
		cores = append(cores, zapcore.NewCore(jsonEncoder, writer, atomicLevel))
	}

	// 多输出核心
	core := zapcore.NewTee(
		cores...,
	)

	// 创建logger
	defaultLogger = &Logger{
		log: zap.New(core,
			zap.AddCaller(),
			zap.AddCallerSkip(2),
		),
	}

	return nil
}

// With 创建携带固定字段的新Logger
func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{
		log: l.log.With(fields...),
	}
}

func (l *Logger) WithField(key string, value interface{}) *Logger {
	return l.With(zap.Any(key, value))
}

// WithField 使用单个字段创建新的Logger
func WithField(key string, value interface{}) *Logger {
	return defaultLogger.With(zap.Any(key, value))
}

// WithFields 使用多个字段创建新的Logger
func WithFields(fields map[string]interface{}) *Logger {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return defaultLogger.With(zapFields...)
}

// Debug 输出debug级别日志
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.log.Debug(msg, fields...)
}

// Info 输出info级别日志
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.log.Info(msg, fields...)
}

// Warn 输出warn级别日志
func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.log.Warn(msg, fields...)
}

// Error 输出error级别日志
func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.log.Error(msg, fields...)
}

// Fatal 输出fatal级别日志
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.log.Fatal(msg, fields...)
}

// 全局方法
func Debug(msg string, fields ...zap.Field) {
	defaultLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	defaultLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	defaultLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	defaultLogger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	defaultLogger.Fatal(msg, fields...)
}
