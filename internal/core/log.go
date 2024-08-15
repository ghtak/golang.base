package core

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

const (
	LogClassConsole string = "console"
	LogClassFile           = "file"
)

const (
	LogEncoderConsole string = "console"
	LogEncoderJson           = "json"
)

func NewLogger(env Env) *zap.Logger {
	handlers := env.viper.GetStringSlice("log.handlers")
	var cores []zapcore.Core
	for _, handler := range handlers {
		core := NewCore(handler, env)
		if core != nil {
			cores = append(cores, core)
		}
	}
	return zap.New(zapcore.NewTee(cores...))
}

func NewCore(handler string, env Env) zapcore.Core {
	subEnv := env.Sub(fmt.Sprintf("log.%s", handler))
	class := subEnv.GetString("class", "")
	switch class {
	case LogClassConsole:
		return NewConsoleCore(NewLogConsoleEnv(subEnv))
	case LogClassFile:
		return NewFileCore(NewLogFileEnv(subEnv))
	}
	return nil
}

func NewConsoleCore(env LogConsoleEnv) zapcore.Core {
	level, err := zapcore.ParseLevel(env.Level)
	if err != nil {
		return nil
	}
	return zapcore.NewCore(NewEncoder(env.Encoder), zapcore.AddSync(os.Stdout), level)
}

func NewFileCore(env LogFileEnv) zapcore.Core {
	level, err := zapcore.ParseLevel(env.Level)
	if err != nil {
		return nil
	}
	return zapcore.NewCore(
		NewEncoder(env.Encoder),
		zapcore.AddSync(&lumberjack.Logger{
			Filename:   env.Filename,
			MaxSize:    env.MaxSize,
			MaxBackups: env.MaxBackups,
			MaxAge:     env.MaxAge,
		}),
		level)
}

func NewZapConfig() zapcore.EncoderConfig {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return encoderConfig
}

func NewEncoder(encoder string) zapcore.Encoder {
	if encoder == LogEncoderJson {
		return zapcore.NewJSONEncoder(NewZapConfig())
	}
	return zapcore.NewConsoleEncoder(NewZapConfig())
}

type LogEnv struct {
	Class   string
	Level   string
	Encoder string
}

type LogConsoleEnv struct {
	LogEnv
}

func NewLogConsoleEnv(env SubEnv) LogConsoleEnv {
	return LogConsoleEnv{
		LogEnv{
			Class:   env.GetString("class", ""),
			Level:   env.GetString("level", ""),
			Encoder: env.GetString("encoder", ""),
		},
	}
}

type LogFileEnv struct {
	LogEnv
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
}

func NewLogFileEnv(env SubEnv) LogFileEnv {
	return LogFileEnv{
		LogEnv: LogEnv{
			Class:   env.GetString("class", ""),
			Level:   env.GetString("level", ""),
			Encoder: env.GetString("encoder", ""),
		},
		Filename:   env.GetString("filename", "./logs/base.log"),
		MaxSize:    env.GetInt("max_size", 500), // megabytes
		MaxBackups: env.GetInt("max_backups", 3),
		MaxAge:     env.GetInt("max_age", 28), // days
	}
}
