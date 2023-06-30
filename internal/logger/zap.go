package logger

import "go.uber.org/zap"

var log *zap.Logger = nil

func Initialize(level string) *zap.Logger {

	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		panic(err)
	}
	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	zl, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	log = zl
	return zl
}
func GetLogger() *zap.Logger {
	return log
}
