package logger

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const ctxLoggerKey = "zapLogger"

type Logger struct {
	*zap.Logger
}

func NewLogger(conf *viper.Viper) *Logger {
	firstString := func(defaultValue string, keys ...string) string {
		for _, key := range keys {
			if val := conf.GetString(key); val != "" {
				return val
			}
		}
		return defaultValue
	}
	firstInt := func(defaultValue int, keys ...string) int {
		for _, key := range keys {
			if conf.IsSet(key) {
				if v := conf.GetInt(key); v != 0 {
					return v
				}
			}
		}
		return defaultValue
	}
	firstBool := func(defaultValue bool, keys ...string) bool {
		for _, key := range keys {
			if conf.IsSet(key) {
				return conf.GetBool(key)
			}
		}
		return defaultValue
	}

	levelValue := strings.ToLower(firstString("info", "logger.level", "log.log_level"))
	var level zapcore.Level
	switch levelValue {
	case "debug":
		level = zap.DebugLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}

	fileEnabled := firstBool(false, "logger.file.enabled", "log.enable_file_logger")
	logPath := firstString("./storage/logs/server.log", "logger.file.path", "log.log_file_name")
	encoderStyle := firstString("json", "logger.encoding", "log.encoding")

	hook := lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    firstInt(100, "logger.file.max_size", "log.max_size"),
		MaxBackups: firstInt(3, "logger.file.max_backups", "log.max_backups"),
		MaxAge:     firstInt(28, "logger.file.max_age", "log.max_age"),
		Compress:   firstBool(false, "logger.file.compress", "log.compress"),
	}

	var encoder zapcore.Encoder
	if encoderStyle == "console" {
		encoder = zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "Logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseColorLevelEncoder,
			EncodeTime:     timeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.FullCallerEncoder,
		})
	} else {
		encoder = zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.EpochTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		})
	}
	writers := []zapcore.WriteSyncer{zapcore.AddSync(os.Stdout)}
	if fileEnabled {
		writers = append(writers, zapcore.AddSync(&hook))
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(writers...),
		level,
	)

	env := strings.ToLower(firstString("", "app.env", "env"))
	if env != "prod" && env != "production" {
		return &Logger{zap.New(core, zap.Development(), zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))}
	}
	return &Logger{zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))}
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	//enc.AppendString(t.Format("2006-01-02 15:04:05"))
	enc.AppendString(t.Format("2006-01-02 15:04:05.000000000"))
}

// WithValue Adds a field to the specified context
func (l *Logger) WithValue(ctx context.Context, fields ...zapcore.Field) context.Context {
	if c, ok := ctx.(*gin.Context); ok {
		ctx = c.Request.Context()
		c.Request = c.Request.WithContext(context.WithValue(ctx, ctxLoggerKey, l.WithContext(ctx).With(fields...)))
		return c
	}
	return context.WithValue(ctx, ctxLoggerKey, l.WithContext(ctx).With(fields...))
}

// WithContext Returns a zap instance from the specified context
func (l *Logger) WithContext(ctx context.Context) *Logger {
	if c, ok := ctx.(*gin.Context); ok {
		ctx = c.Request.Context()
	}
	zl := ctx.Value(ctxLoggerKey)
	ctxLogger, ok := zl.(*zap.Logger)
	if ok {
		return &Logger{ctxLogger}
	}
	return l
}
