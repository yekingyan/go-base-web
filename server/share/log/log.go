package sharelog

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogLevel is the log level.
const LogLevel zapcore.Level = zap.DebugLevel

func formatEncodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(fmt.Sprintf("%d%02d%02d_%02d%02d%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second()))
}

// InitZapLog is logger creator.
func InitZapLog(name string) *zap.Logger {
	logDir := "/tmp"
	infoPath := filepath.Join(logDir, name+"_info.log")
	errPath := filepath.Join(logDir, name+"_err.log")
	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(LogLevel),
		Development: true,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "t",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "trace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     formatEncodeTime,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{infoPath, "stdout"},
		ErrorOutputPaths: []string{errPath, "stdout"},
		InitialFields: map[string]interface{}{
			"app": name,
		},
	}
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err = os.Mkdir(logDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	logger, err := cfg.Build()
	if err != nil {
		panic("log init fail:" + err.Error())
	}
	return logger
}
