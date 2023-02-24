package log

import (
	"bytes"
	"fmt"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"tikapp/common/config"
	"time"
)

// Logger 整个项目的Logger
var Logger *zap.Logger

const (
	red    = 31
	yellow = 33
	blue   = 36
	gray   = 37
)

func Init() {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&LogFormatter{})
	logrus.Warn("logrus show test")

	if config.AppCfg.RunMode == "debug" {
		// 开发模式 日志输出到终端
		core := zapcore.NewTee(
			zapcore.NewCore(getEncoder(),
				zapcore.Lock(os.Stdout), zapcore.DebugLevel),
		)
		Logger = zap.New(core, zap.AddCaller())
	} else {
		fileLog()
	}
}

func fileLog() {
	// 调试级别
	debugPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.DebugLevel
	})
	// 日志级别
	infoPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.InfoLevel
	})
	// 警告级别
	warnPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.WarnLevel
	})
	// 错误级别
	errorPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.ErrorLevel
	})
	// panic级别
	panicPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.PanicLevel
	})
	// fatal级别
	fatalPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.FatalLevel
	})

	cores := [...]zapcore.Core{
		getEncoderCore("./debug.log", debugPriority),
		getEncoderCore("./info.log", infoPriority),
		getEncoderCore("./warn.log", warnPriority),
		getEncoderCore("./error.log", errorPriority),
		getEncoderCore("./panic.log", panicPriority),
		getEncoderCore("./fatal.log", fatalPriority),
	}

	// zap.AddCaller() 可以获取到文件名和行号
	Logger = zap.New(zapcore.NewTee(cores[:]...), zap.AddCaller())
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func getLogWriter(fileName string) zapcore.WriteSyncer {
	dir, _ := os.Getwd() // 获取项目目录
	sperator0 := os.PathSeparator
	sperator := string(sperator0)
	// 	log 存放路径
	dir = dir + sperator + "runtime" + sperator + "logs"
	if !pathExists(dir) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			logrus.Warnf("create dir %s failed", dir)
		}
	}
	lumberJackLogger := &lumberjack.Logger{
		Filename:   dir + sperator + fileName, // 日志文件路径
		MaxSize:    5,                         // 设置日志文件最大尺寸
		MaxBackups: 5,                         // 设置日志文件最多保存多少个备份
		MaxAge:     30,                        // 设置日志文件最多保存多少天
		Compress:   true,                      // 是否压缩 disabled by default
	}
	// 返回同步方式写入日志文件的zapcore.WriteSyncer
	return zapcore.AddSync(lumberJackLogger)
}

func getEncoderCore(fileName string, level zapcore.LevelEnabler) (core zapcore.Core) {
	writer := getLogWriter(fileName)
	return zapcore.NewCore(getEncoder(), writer, level)
}

func getEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(getEncoderConfig())
}

func getEncoderConfig() (config zapcore.EncoderConfig) {
	config = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 将日志级别字符串转化为小写
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder, // 执行消耗时间转化成浮点型的秒
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 以包/文件:行号 格式化调用堆栈
	}
	return config
}

// CustomTimeEncoder 自定义日志输出时间格式
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(config.LogCfg.TimeFormat))
}

// for logrus
type LogFormatter struct{}

//实现Formatter(entry *logrus.Entry) ([]byte, error)接口
func (t *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	//根据不同的level去展示颜色
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = gray
	case logrus.WarnLevel:
		levelColor = yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = red
	default:
		levelColor = blue
	}
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	//自定义日期格式
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	if entry.HasCaller() {
		//自定义文件路径
		funcVal := entry.Caller.Function
		fileVal := fmt.Sprintf("%s:%d", path.Base(entry.Caller.File), entry.Caller.Line)
		//自定义输出格式
		fmt.Fprintf(b, "[%s] \x1b[%dm[%s]\x1b[0m %s %s %s\n", timestamp, levelColor, entry.Level, fileVal, funcVal, entry.Message)
	} else {
		fmt.Fprintf(b, "[%s] \x1b[%dm[%s]\x1b[0m  %s\n", timestamp, levelColor, entry.Level, entry.Message)
	}
	return b.Bytes(), nil
}

//输出
// [2021-05-11 15:08:46] [info] [demo.go:38 main.Demo] i'm demo
