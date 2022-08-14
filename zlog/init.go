package zlog

import (
	"github.com/derekAHua/goLib/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

type (
	LogConfig struct {
		Level  zapcore.Level `yaml:"level"`
		Stdout bool          `yaml:"stdout"`
	}

	loggerConfig struct {
		ZapLevel zapcore.Level

		Stdout bool
		Path   string
	}
)

// Only change by Init function.
var (
	logConfig = loggerConfig{
		ZapLevel: zapcore.InfoLevel,
		Stdout:   false,
		Path:     "./log",
	}

	onceLogInit sync.Once
)

func Init(conf LogConfig) {
	onceLogInit.Do(func() {
		logConfig.ZapLevel = conf.Level
		logConfig.Stdout = conf.Stdout
		logConfig.Path = env.GetLogDirPath()

		zapLogs := []string{LogNameServer, LogNameAccess, LogNameMysql, LogNameRedis, LogNameLua, LogNameRMQ, LogNameRpc, LogNameES}

		for _, v := range zapLogs {
			if _, ok := mapZapLogger[v]; !ok {
				mapZapLogger[v] = newLogger(v).WithOptions(zap.AddCallerSkip(1))
			}
		}

		sugaredLogger = mapZapLogger[LogNameServer].Sugar()
	})
}
