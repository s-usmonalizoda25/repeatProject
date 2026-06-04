package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)
type Logger struct {
	*zap.Logger
	Audit *zap.Logger
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
	mainLogger, err:=cfg.Build()
	if err!=nil{
		return nil, fmt.Errorf("mainLogger.Build: %w", err)
	}
	auditCfg:=zap.NewProductionConfig()
	auditCfg.EncoderConfig.TimeKey="timestamp"
	auditCfg.EncoderConfig.EncodeTime=zapcore.ISO8601TimeEncoder
	auditCfg.OutputPaths=[]string{"data/auditLog.txt"}
	auditLogger, err:=auditCfg.Build()
	if err!=nil{
		return nil, fmt.Errorf("auditLogger.Build:%w", err)
	}
	return &Logger{
		Logger: mainLogger,
		Audit: auditLogger,
	}, nil
}