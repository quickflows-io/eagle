package log

import (
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-eagle/eagle/pkg/utils"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// WriterConsole console output
	WriterConsole = "console"
	// WriterFile file output
	WriterFile = "file"
)

const (
	// RotateTimeDaily cut by day
	RotateTimeDaily = "daily"
	// RotateTimeHourly cut by the hour
	RotateTimeHourly = "hourly"
)

const defaultSkip = 1 // zapLogger wraps a layer of zap.Logger, which skips one layer by default

// For mapping config logger to app logger levels
var loggerLevelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

// Prevent data race from occurring during zap.AddStacktrace
var zapStacktraceMutex sync.Mutex

func getLoggerLevel(cfg *Config) zapcore.Level {
	level, exist := loggerLevelMap[cfg.Level]
	if !exist {
		return zapcore.DebugLevel
	}

	return level
}

// zapLogger logger struct
type zapLogger struct {
	sugarLogger *zap.SugaredLogger
}

// newZapLogger new zap logger
func newZapLogger(cfg *Config) (*zap.Logger, error) {
	return buildLogger(cfg, defaultSkip), nil
}

// newLoggerWithCallerSkip new logger with caller skip
func newLoggerWithCallerSkip(cfg *Config, skip int) (Logger, error) {
	return &zapLogger{sugarLogger: buildLogger(cfg, defaultSkip+skip).Sugar()}, nil
}

// newLogger new logger
func newLogger(cfg *Config) (Logger, error) {
	return newLoggerWithCallerSkip(cfg, 0)
}

func buildLogger(cfg *Config, skip int) *zap.Logger {
	var encoderCfg zapcore.EncoderConfig
	if cfg.Development {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	var encoder zapcore.Encoder
	if cfg.Encoding == WriterConsole {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	var cores []zapcore.Core
	var options []zap.Option
	// init option
	hostname, _ := os.Hostname()
	option := zap.Fields(
		zap.String("ip", utils.GetLocalIP()),
		zap.String("app_id", cfg.Name),
		zap.String("instance_id", hostname),
	)
	options = append(options, option)

	writers := strings.Split(cfg.Writers, ",")
	for _, w := range writers {
		switch w {
		case WriterConsole:
			cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), getLoggerLevel(cfg)))
		case WriterFile:
			// info
			cores = append(cores, getInfoCore(encoder, cfg))

			// warning
			core, option := getWarnCore(encoder, cfg)
			cores = append(cores, core)
			if option != nil {
				options = append(options, option)
			}

			// error
			core, option = getErrorCore(encoder, cfg)
			cores = append(cores, core)
			if option != nil {
				options = append(options, option)
			}
		default:
			// console
			cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), getLoggerLevel(cfg)))
			// file
			cores = append(cores, getAllCore(encoder, cfg))
		}
	}

	combinedCore := zapcore.NewTee(cores...)

	// Open development mode, stack trace
	if !cfg.DisableCaller {
		caller := zap.AddCaller()
		options = append(options, caller)
	}

	// Skip file call levels
	addCallerSkip := zap.AddCallerSkip(skip)
	options = append(options, addCallerSkip)

	// Construction log
	return zap.New(combinedCore, options...)
}

func getAllCore(encoder zapcore.Encoder, cfg *Config) zapcore.Core {
	allWriter := getLogWriterWithTime(cfg, cfg.LoggerFile)
	allLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl <= zapcore.FatalLevel
	})
	return zapcore.NewCore(encoder, zapcore.AddSync(allWriter), allLevel)
}

func getInfoCore(encoder zapcore.Encoder, cfg *Config) zapcore.Core {
	infoWrite := getLogWriterWithTime(cfg, cfg.LoggerFile)
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl <= zapcore.InfoLevel
	})
	return zapcore.NewCore(encoder, zapcore.AddSync(infoWrite), infoLevel)
}

func getWarnCore(encoder zapcore.Encoder, cfg *Config) (zapcore.Core, zap.Option) {
	warnWrite := getLogWriterWithTime(cfg, cfg.LoggerWarnFile)
	var stacktrace zap.Option
	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		if !cfg.DisableCaller {
			zapStacktraceMutex.Lock()
			stacktrace = zap.AddStacktrace(zapcore.WarnLevel)
			zapStacktraceMutex.Unlock()
		}
		return lvl == zapcore.WarnLevel
	})
	return zapcore.NewCore(encoder, zapcore.AddSync(warnWrite), warnLevel), stacktrace
}

func getErrorCore(encoder zapcore.Encoder, cfg *Config) (zapcore.Core, zap.Option) {
	errorFilename := cfg.LoggerErrorFile
	errorWrite := getLogWriterWithTime(cfg, errorFilename)
	var stacktrace zap.Option
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		if !cfg.DisableCaller {
			zapStacktraceMutex.Lock()
			stacktrace = zap.AddStacktrace(zapcore.ErrorLevel)
			zapStacktraceMutex.Unlock()
		}
		return lvl >= zapcore.ErrorLevel
	})
	return zapcore.NewCore(encoder, zapcore.AddSync(errorWrite), errorLevel), stacktrace
}

// getLogWriterWithTime cuts by time (hours)
func getLogWriterWithTime(cfg *Config, filename string) io.Writer {
	logFullPath := filename
	rotationPolicy := cfg.LogRollingPolicy
	backupCount := cfg.LogBackupCount
	// 默认
	var rotateDuration time.Duration
	if rotationPolicy == RotateTimeHourly {
		rotateDuration = time.Hour
	} else if rotationPolicy == RotateTimeDaily {
		rotateDuration = time.Hour * 24
	}
	hook, err := rotatelogs.New(
		logFullPath+".%Y%m%d%H",                     // The time format uses the shell's date time format
		rotatelogs.WithLinkName(logFullPath),        // Generate a soft link pointing to the latest log file
		rotatelogs.WithRotationCount(backupCount),   // Maximum number of files to save
		rotatelogs.WithRotationTime(rotateDuration), // log cutting interval
	)

	if err != nil {
		panic(err)
	}
	return hook
}

// Debug logger
func (l *zapLogger) Debug(args ...interface{}) {
	l.sugarLogger.Debug(args...)
}

// Info logger
func (l *zapLogger) Info(args ...interface{}) {
	l.sugarLogger.Info(args...)
}

// Warn logger
func (l *zapLogger) Warn(args ...interface{}) {
	l.sugarLogger.Warn(args...)
}

// Error logger
func (l *zapLogger) Error(args ...interface{}) {
	l.sugarLogger.Error(args...)
}

func (l *zapLogger) Fatal(args ...interface{}) {
	l.sugarLogger.Fatal(args...)
}

func (l *zapLogger) Debugf(format string, args ...interface{}) {
	l.sugarLogger.Debugf(format, args...)
}

func (l *zapLogger) Infof(format string, args ...interface{}) {
	l.sugarLogger.Infof(format, args...)
}

func (l *zapLogger) Warnf(format string, args ...interface{}) {
	l.sugarLogger.Warnf(format, args...)
}

func (l *zapLogger) Errorf(format string, args ...interface{}) {
	l.sugarLogger.Errorf(format, args...)
}

func (l *zapLogger) Fatalf(format string, args ...interface{}) {
	l.sugarLogger.Fatalf(format, args...)
}

func (l *zapLogger) Panicf(format string, args ...interface{}) {
	l.sugarLogger.Panicf(format, args...)
}

func (l *zapLogger) WithFields(fields Fields) Logger {
	var f = make([]interface{}, 0)
	for k, v := range fields {
		f = append(f, k)
		f = append(f, v)
	}
	newLogger := l.sugarLogger.With(f...)
	return &zapLogger{newLogger}
}
