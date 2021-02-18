package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// CustomTimeEncoder function of own formulating time for output to the log
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// GetLoggerConfigFromFile - create a logger for settings from a file
func GetLoggerConfigFromFile(pathToConfig string) (*zap.Logger, error) {
	var err error
	var configRaw []byte

	if configRaw, err = ioutil.ReadFile(pathToConfig); err != nil {
		return nil, err
	}
	zapConfig := zap.Config{}
	if err = json.Unmarshal(configRaw, &zapConfig); err != nil {
		return nil, err
	}
	zapConfig.EncoderConfig.EncodeTime = CustomTimeEncoder

	var logger *zap.Logger
	if logger, err = zapConfig.Build(); err != nil {
		return nil, err
	}
	return logger, nil
}

// GetLoggerConfigFromFileWithRotate - GetLoggerConfigFromFile + log rotate
func GetLoggerConfigFromFileWithRotate(pathToConfig string, maxsize, MaxBackups, MaxAge int, compress bool) (*zap.Logger, error) {
	var err error
	var configRaw []byte

	if configRaw, err = ioutil.ReadFile(pathToConfig); err != nil {
		return nil, err
	}
	zapConfig := zap.Config{}
	if err = json.Unmarshal(configRaw, &zapConfig); err != nil {
		return nil, err
	}
	zapConfig.EncoderConfig.EncodeTime = CustomTimeEncoder

	if len(zapConfig.OutputPaths) == 0 {
		return nil, fmt.Errorf("zap config has no output file")
	}
	writeSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   zapConfig.OutputPaths[0],
		MaxSize:    maxsize, // megabytes
		MaxBackups: MaxBackups,
		MaxAge:     MaxAge, // days
		Compress:   compress,
	})
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zapConfig.EncoderConfig),
		writeSyncer,
		zapConfig.Level,
	)

	var logger *zap.Logger
	logger = zap.New(core, zap.AddCaller())
	//if logger, err = zapConfig.Build(); err != nil {
	//	return nil, err
	//}

	return logger, nil
}

// GetDefaultZapConfig - return default zap Config
func GetDefaultZapConfig() zap.Config {
	var err error
	configRaw := []byte(`{
	"level":"info",
	"encoding":"console",
	"outputPaths": ["stdout"],
	"errorOutputPaths": ["stderr"],
	"encoderConfig": {
		"messageKey": "message",
		"levelKey": "level",
		"nameKey": "logger",
		"timeKey": "time",
		"stacktraceKey": "stacktrace",
		"callstackKey": "callstack",
		"errorKey": "error",
		"timeEncoder": "rfc3339",
		"fileKey": "file",
		"levelEncoder": "capitalColor",
		"durationEncoder": "second",
		"sampling": {
			"initial": "3",
			"thereafter": "10"
		}
	}
}`)
	zapConfig := zap.Config{}
	if err = json.Unmarshal(configRaw, &zapConfig); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(2)
	}
	zapConfig.EncoderConfig.EncodeTime = CustomTimeEncoder
	return zapConfig
}

func GetStandartLogger(loggerType string) (*zap.Logger, error) {
	var err error
	var logger *zap.Logger
	var config zap.Config = GetDefaultZapConfig()
	switch loggerType {
	case "production":
		config.Level.SetLevel(zapcore.InfoLevel)
	case "development":
		config.Level.SetLevel(zapcore.DebugLevel)
	default:
		config.Level.SetLevel(zapcore.InfoLevel)
	}
	if logger, err = config.Build(); err != nil {
		return nil, err
	}
	return logger, nil
}
