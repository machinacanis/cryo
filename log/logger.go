package log

// CryoLogLevel 是 cryo 的日志级别
type CryoLogLevel int

const (
	DebugLevel CryoLogLevel = iota
	InfoLevel
	SuccessLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
)

// CryoLogger 是 cryo 的日志记录器接口
//
// cryo 默认使用 logrus 作为日志记录器，但是在使用过程中，可能会碰到需要嵌入已有项目并复用已有的日志记录器的情况
//
// 通过实现这个接口，可以将 cryo 的日志记录器替换为已有的日志记录器
type CryoLogger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Success(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Successf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	Print(args ...interface{})
	Printf(format string, args ...interface{})
}
