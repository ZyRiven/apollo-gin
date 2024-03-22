package setting

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Option custom setup config
type Option func(*option)

type option struct {
	level          zapcore.Level
	fields         map[string]string
	file           io.Writer
	timeLayout     string
	disableConsole bool
}

var ZAPS *zap.SugaredLogger

const (
	// DefaultLevel the default log level
	DefaultLevel = zapcore.InfoLevel

	// DefaultTimeLayout the default time layout;
	DefaultTimeLayout = time.DateTime
)

// InitLogger 日志初始化
func InitLogger() {
	var zapLogger *zap.Logger
	var err error

	if LogToFile == true {
		exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		absPath := exeCurDir + LogFile
		zapLogger, err = NewJSONLogger(WithTimeLayout(DefaultTimeLayout), WithFileRotationP(absPath))
	} else {
		zapLogger, err = NewJSONLogger(WithTimeLayout(DefaultTimeLayout))
	}
	if err != nil {
		panic(err)
	}

	ZAPS = zapLogger.Sugar()
	ZAPS.Infof("zap日志 初始化成功!日志等级:%s", LogLevel)
}

// WithTimeLayout custom time format
func WithTimeLayout(timeLayout string) Option {
	return func(opt *option) {
		opt.timeLayout = timeLayout
	}
}

// WithFileRotationP write log to some file with rotation
func WithFileRotationP(file string) Option {
	dir := filepath.Dir(file)
	if err := os.MkdirAll(dir, 0766); err != nil {
		panic(err)
	}

	return func(opt *option) {
		opt.file = &lumberjack.Logger{ // concurrent-safed
			Filename:   file,           // 文件路径
			MaxSize:    LogFileMaxSize, // 单个文件最大尺寸，默认单位 M
			MaxBackups: LogFileBackup,  // 最多保留3个备份
			MaxAge:     7,              // 日志文件最多保存天数
			LocalTime:  true,           // 使用本地时间
			Compress:   true,           // 是否压缩
		}
	}
}

// NewJSONLogger return a json-encoder zap logger,
func NewJSONLogger(opts ...Option) (*zap.Logger, error) {
	level := Level(LogLevel)
	opt := &option{
		level:          level,
		fields:         make(map[string]string),
		disableConsole: false,
	}
	for _, f := range opts {
		f(opt)
	}

	timeLayout := DefaultTimeLayout
	if opt.timeLayout != "" {
		timeLayout = opt.timeLayout
	}

	// similar to zap.NewProductionEncoderConfig()
	fileConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger", // used by logger.Named(key); optional; useless
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace", // use by zap.AddStacktrace; optional; useless
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(timeLayout))
		},
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 全路径编码器
	}

	consoleConfig := zapcore.EncoderConfig{}
	if opt.level < zapcore.InfoLevel {
		config := zapcore.EncoderConfig{
			TimeKey:       "time",
			LevelKey:      "level",
			NameKey:       "logger", // used by logger.Named(key); optional; useless
			CallerKey:     "caller",
			MessageKey:    "msg",
			StacktraceKey: "stacktrace", // use by zap.AddStacktrace; optional; useless
			LineEnding:    zapcore.DefaultLineEnding,
			EncodeLevel:   zapcore.LowercaseLevelEncoder,
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format(timeLayout))
			},
			EncodeDuration: zapcore.MillisDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder, // 全路径编码器
		}
		consoleConfig = config
	} else {
		config := zapcore.EncoderConfig{
			TimeKey:       "time",
			LevelKey:      "level",
			NameKey:       "logger", // used by logger.Named(key); optional; useless
			CallerKey:     "caller",
			MessageKey:    "msg",
			StacktraceKey: "stacktrace", // use by zap.AddStacktrace; optional; useless
			LineEnding:    zapcore.DefaultLineEnding,
			EncodeLevel:   zapcore.LowercaseLevelEncoder,
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format(timeLayout))
			},
			EncodeDuration: zapcore.MillisDurationEncoder,
			//EncodeCaller:   zapcore.ShortCallerEncoder, // 全路径编码器
		}
		consoleConfig = config
	}

	fileEncoder := zapcore.NewJSONEncoder(fileConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(consoleConfig)

	// lowPriority usd by info\debug\warn
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= opt.level && lvl < zapcore.ErrorLevel
	})

	// highPriority usd by error\panic\fatal
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= opt.level && lvl >= zapcore.ErrorLevel
	})

	stdout := zapcore.Lock(os.Stdout) // lock for concurrent safe
	stderr := zapcore.Lock(os.Stderr) // lock for concurrent safe

	core := zapcore.NewTee()

	if !opt.disableConsole {
		core = zapcore.NewTee(
			zapcore.NewCore(consoleEncoder,
				zapcore.NewMultiWriteSyncer(stdout),
				lowPriority,
			),
			zapcore.NewCore(consoleEncoder,
				zapcore.NewMultiWriteSyncer(stderr),
				highPriority,
			),
		)
	}

	if opt.file != nil {
		core = zapcore.NewTee(core,
			zapcore.NewCore(fileEncoder,
				zapcore.AddSync(opt.file),
				zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
					return lvl >= opt.level
				}),
			),
		)
	}

	logger := zap.New(core,
		zap.AddCaller(),
		zap.ErrorOutput(stderr),
	)

	for key, value := range opt.fields {
		logger = logger.WithOptions(zap.Fields(zapcore.Field{
			Key:    key,
			Type:   zapcore.StringType,
			String: value,
		}))
	}
	return logger, nil
}

func Level(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return DefaultLevel
	}
}