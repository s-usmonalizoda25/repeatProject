package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)
type Logger struct {
	*zap.Logger
}

func New(devMode bool) (*Logger, error) {
	var cfg zap.Config
	if devMode {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.OutputPaths = []string{"stdout"}
	l, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	
	return &Logger{
		l,
	}, nil
}