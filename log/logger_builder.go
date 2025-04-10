package log

import (
	"bufio"
	"github.com/sirupsen/logrus"
	"io"
)

// Logger 单个日志记录器
type Logger struct {
	logger    *logrus.Logger
	formatter logrus.Formatter
	level     logrus.Level
}

// LoggerBuilder 日志记录器生成器
type LoggerBuilder struct {
	loggers      []Logger
	defaultLevel CryoLogLevel
}

func (b *LoggerBuilder) Debug(args ...interface{}) {
	for _, l := range b.loggers {
		l.logger.Trace(args...)
	}
}

func (b *LoggerBuilder) Info(args ...interface{}) {
	for _, l := range b.loggers {
		l.logger.Debug(args...)
	}
}

func (b *LoggerBuilder) Success(args ...interface{}) {
	for _, l := range b.loggers {
		l.logger.Info(args...)
	}
}

func (b *LoggerBuilder) Warn(args ...interface{}) {
	for _, l := range b.loggers {
		l.logger.Warn(args...)
	}
}

func (b *LoggerBuilder) Error(args ...interface{}) {
	for _, l := range b.loggers {
		l.logger.Error(args...)
	}
}

func (b *LoggerBuilder) Fatal(args ...interface{}) {
	for _, l := range b.loggers {
		l.logger.Fatal(args...)
	}
}

func (b *LoggerBuilder) Panic(args ...interface{}) {
	for _, l := range b.loggers {
		l.logger.Panic(args...)
	}
}

func (b *LoggerBuilder) Debugf(format string, args ...interface{}) {
	for _, l := range b.loggers {
		l.logger.Tracef(format, args...)
	}
}

func (b *LoggerBuilder) Infof(format string, args ...interface{}) {
	for _, l := range b.loggers {
		l.logger.Debugf(format, args...)
	}
}

func (b *LoggerBuilder) Successf(format string, args ...interface{}) {
	for _, l := range b.loggers {
		l.logger.Infof(format, args...)
	}
}

func (b *LoggerBuilder) Warnf(format string, args ...interface{}) {
	for _, l := range b.loggers {
		l.logger.Warnf(format, args...)
	}
}

func (b *LoggerBuilder) Errorf(format string, args ...interface{}) {
	for _, l := range b.loggers {
		l.logger.Errorf(format, args...)
	}
}

func (b *LoggerBuilder) Fatalf(format string, args ...interface{}) {
	for _, l := range b.loggers {
		l.logger.Fatalf(format, args...)
	}
}

func (b *LoggerBuilder) Panicf(format string, args ...interface{}) {
	for _, l := range b.loggers {
		l.logger.Panicf(format, args...)
	}
}

func (b *LoggerBuilder) Print(args ...interface{}) {
	for _, l := range b.loggers {
		l.logger.Info(args...)
	}
}

func (b *LoggerBuilder) Printf(format string, args ...interface{}) {
	for _, l := range b.loggers {
		l.logger.Infof(format, args...)
	}
}

// NewLoggerBuilder 创建一个新的日志记录器构建器
func NewLoggerBuilder(level ...CryoLogLevel) *LoggerBuilder {
	if len(level) == 0 {
		level = append(level, InfoLevel)
	}

	return &LoggerBuilder{defaultLevel: level[0]}
}

// AddTextLogger 添加一个文本日志记录器
func (b *LoggerBuilder) AddTextLogger(level CryoLogLevel, formatter ...logrus.Formatter) *LoggerBuilder {
	if len(formatter) == 0 {
		formatter = append(formatter, NewDefaultFormatterBuilder())
	}

	logger := logrus.New()
	logger.SetFormatter(formatter[0])
	logger.SetLevel(ConvertCryoLogLevelToLogrusLevel(level)) // 设置日志级别
	logger.SetOutput(logrus.StandardLogger().Out)

	b.loggers = append(b.loggers, Logger{
		logger:    logger,
		formatter: formatter[0],
		level:     ConvertCryoLogLevelToLogrusLevel(level),
	})

	return b
}

// AddFileLogger 添加一个文件日志记录器
func (b *LoggerBuilder) AddFileLogger(level CryoLogLevel, file io.Writer, formatter ...logrus.Formatter) *LoggerBuilder {
	if len(formatter) == 0 {
		formatter = append(formatter, &logrus.TextFormatter{})
	}

	logger := logrus.New()
	logger.SetFormatter(formatter[0])
	logger.SetLevel(ConvertCryoLogLevelToLogrusLevel(level)) // 设置日志级别

	bufferedFile := bufio.NewWriter(file) // 使用bufio.NewWriter创建一个带缓冲的io.Writer对象
	logger.SetOutput(bufferedFile)

	b.loggers = append(b.loggers, Logger{
		logger:    logger,
		formatter: formatter[0],
		level:     ConvertCryoLogLevelToLogrusLevel(level),
	})

	return b
}

// AddJsonFileLogger 添加一个JSON文件日志记录器
func (b *LoggerBuilder) AddJsonFileLogger(level CryoLogLevel, file io.Writer, formatter ...logrus.Formatter) *LoggerBuilder {
	if len(formatter) == 0 {
		formatter = append(formatter, &logrus.JSONFormatter{})
	}

	logger := logrus.New()
	logger.SetFormatter(formatter[0])
	logger.SetLevel(ConvertCryoLogLevelToLogrusLevel(level)) // 设置日志级别
	bufferedFile := bufio.NewWriter(file)                    // 使用bufio.NewWriter创建一个带缓冲的io.Writer对象
	logger.SetOutput(bufferedFile)

	b.loggers = append(b.loggers, Logger{
		logger:    logger,
		formatter: formatter[0],
		level:     ConvertCryoLogLevelToLogrusLevel(level),
	})

	return b
}

// AddLogger 添加一个自定义日志记录器
func (b *LoggerBuilder) AddLogger(logger Logger) *LoggerBuilder {
	b.loggers = append(b.loggers, logger)
	return b
}
