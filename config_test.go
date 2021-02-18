package config

import (
	"fmt"
	"go.uber.org/zap"
	"testing"
)

const testConfig = "test_config/zapConfig.json"

func TestGetCrontab(t *testing.T) {
	var err error
	var logger *zap.Logger
	if logger, err = GetLoggerConfigFromFileWithRotate(testConfig, 200, 30, 30, true); err != nil {
		_, _ = fmt.Printf("error : %s", err)
	}
	logger.Sugar().Debugf("test")
}
