package core

//
//import (
//	"go.uber.org/zap"
//	"go.uber.org/zap/zapcore"
//)
//
//func NewEncoder(env Env) zapcore.Encoder {
//	encoderConfig := zap.NewProductionEncoderConfig()
//	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
//	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
//	switch env.LogEncoder {
//	case "console":
//		return zapcore.NewConsoleEncoder(encoderConfig)
//	case "json":
//		fallthrough
//	default:
//		return zapcore.NewJSONEncoder(encoderConfig)
//	}
//}
//
//func NewWriter(env Env) zapcore.WriteSyncer {
//
//}
//
//func NewLogger(
//	envs *EnvRepository,
//) *zap.Logger {
//	myEnv := envs.Envs[moduleName].(Env)
//	level, _ := zapcore.ParseLevel(myEnv.LogLevel)
//	core := zapcore.NewCore(NewEncoder(myEnv),
//		writeSyncer,
//		level,
//	)
//	return zap.New(core)
//}
