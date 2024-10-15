package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v3"
)

// Config represents the setting for zap logger.
type Config struct {
	ZapConfig zap.Config        `json:"zap_config" yaml:"zap_config"`
	LogRotate lumberjack.Logger `json:"log_rotate" yaml:"log_rotate"`
}

// Logger is an alternative implementation of *gorm.Logger
type Logger struct {
	Zap *zap.SugaredLogger
}

var log *Logger

func NewLogger(folder string) *Logger {
	configYaml, err := os.ReadFile(filepath.Join(folder, "server.zaplogger.yml"))
	if err != nil {
		fmt.Printf("Failed to read logger configuration: %s", err)
		os.Exit(2)
	}
	var myConfig *Config
	if err = yaml.Unmarshal(configYaml, &myConfig); err != nil {
		fmt.Printf("Failed to read zap logger configuration: %s", err)
		os.Exit(2)
	}
	var zap *zap.Logger
	zap, err = build(myConfig)
	if err != nil {
		fmt.Printf("Failed to compose zap logger : %s", err)
		os.Exit(2)
	}
	sugar := zap.Sugar()
	// set package varriable logger.
	logger := &Logger{Zap: sugar}
	logger.Zap.Infof("Success to read zap logger configuration: server.zaplogger.yml")
	_ = zap.Sync()
	log = logger
	return logger
}

// GetZapLogger returns zapSugaredLogger
func (l *Logger) GetZapLogger() *zap.SugaredLogger {
	return l.Zap
}

func Log() *zap.SugaredLogger {
	if log == nil {
		logger, _ := zap.NewProduction()
		log = &Logger{Zap: logger.Sugar()}
	}
	return log.GetZapLogger()
}

func GetLogger() *Logger {
	Log()
	return log
}
