package logger

import (
	"fmt"

	"go.uber.org/zap"
)

type Logger struct {
	*zap.Logger
	audit *zap.Logger
}

func New(development bool) (*Logger, error) {
	var cfg zap.Config
	if development {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}

	loggy, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build main logger: %w", err)
	}

	auditCfg := zap.NewProductionConfig()
	auditCfg.OutputPaths = []string{"data/auditLog.txt"}

	auditLoggy, err := auditCfg.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build audit logger: %w", err)
	}

	return &Logger{
		Logger: loggy,
		audit:  auditLoggy,
	}, nil
}

func (l *Logger) AuditLogger() *zap.Logger {
	return l.audit
}
