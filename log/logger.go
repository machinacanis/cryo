package log

import (
	"bufio"
	"github.com/sirupsen/logrus"
	"io"
	"sync"
)

// cryobot的日志工具，封装了logrus的一些常用功能
// 为了方便使用，直接提供了对应等级的日志信息函数
//
// logger支持四种日志输出方式，分别是：
// - 对终端进行文本输出（包含颜色和格式）
// - 输出到log文件（包含格式）
// - 输出到json文件
// - （暂未实现）输出到MongoDB数据库
//
// 其中后两种输出方式可以提供日志查询/检索功能，并且有对应的函数封装可供调用
// 每种日志输出方式都可以单独设置日志等级，可以通过同一个日志函数自动调用
//
// 额外的，通过实现Logger接口并替换logger变量，你可以设置使用其他的日志记录器，如zap

func init() {
	// 初始化默认的日志记录器
	logger = &CryoLogger{}
	err := logger.Init()
	if err != nil {
		panic(err)
	}
}

var logger Logger

// Logger 接口定义了日志记录器的基本功能，你可以根据自己的需求实现这个接口来替换默认的日志记录器
type Logger interface {
	Init() error
	Trace(args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
	Tracef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
}

// CryoLogger 封装了logrus的日志记录器，提供了多种日志输出方式
type CryoLogger struct {
	TextLogger     *logrus.Logger
	FileLogger     *logrus.Logger
	JsonFileLogger *logrus.Logger
	MongoLogger    *logrus.Logger
	// 添加这些变量来跟踪活跃的日志记录器
	loggersMutex  sync.RWMutex
	activeLoggers []*logrus.Logger
}

// Init 初始化日志记录器，默认使用单个终端日志记录器
func (cl *CryoLogger) Init() error {
	// 直接调用InitTextLogger()函数
	cl.InitTextLogger()
	return nil
}

// InitTextLogger 初始化终端日志，可以直接传入日志等级进行设置，如果没有传入，则默认使用InfoLevel
//
// 使用logrus的日志等级进行设置，支持以下等级：
// TraceLevel, DebugLevel, InfoLevel, WarnLevel, ErrorLevel, FatalLevel, PanicLevel
//
// 例如：
//
//	log.InitTextLogger(logrus.DebugLevel)
//
// 你也可以直接使用logrus的日志等级对应的整数值进行设置，例如：
//
//	log.InitTextLogger(1)
//
// 其中，6对应TraceLevel，5对应DebugLevel，4对应InfoLevel，3对应WarnLevel，2对应ErrorLevel，1对应FatalLevel，0对应PanicLevel
func (cl *CryoLogger) InitTextLogger(level ...logrus.Level) {
	if len(level) <= 0 { // 如果没有传入日志等级，则默认使用InfoLevel
		level = append(level, logrus.InfoLevel)
	}
	cl.TextLogger = logrus.New()
	cl.TextLogger.SetFormatter(&DefaultDarkFormatter{})
	cl.TextLogger.SetLevel(level[0])
	cl.TextLogger.SetOutput(logrus.StandardLogger().Out)

	cl.updateActiveLoggers() // 更新活跃日志记录器列表
}

// InitFileLogger 初始化文件日志记录器，需要传入一个io.Writer类型的文件对象
//
// 你可以使用os.OpenFile()函数打开一个文件，并将其传入该函数
//
// 例如：
//
//	 file, err := os.OpenFile("log.txt", os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0666)
//
//		if err != nil {
//		    log.Fatal(err)
//		}
//
//	 log.InitFileLogger(file, logrus.DebugLevel)
//
// 你可以单独设置文件日志的日志等级，如果没有传入，则默认使用InfoLevel
//
// 详细的日志等级说明请参考InitTextLogger函数
func (cl *CryoLogger) InitFileLogger(file io.Writer, level ...logrus.Level) {
	if len(level) <= 0 { // 如果没有传入日志等级，则默认使用InfoLevel
		level = append(level, logrus.InfoLevel)
	}

	bufferedFile := bufio.NewWriter(file) // 使用bufio.NewWriter创建一个带缓冲的io.Writer对象

	cl.FileLogger = logrus.New()
	cl.FileLogger.SetFormatter(&logrus.TextFormatter{})
	cl.FileLogger.SetLevel(level[0])
	cl.FileLogger.SetOutput(bufferedFile)

	cl.updateActiveLoggers() // 更新活跃日志记录器列表
}

// InitJsonFileLogger 初始化json文件日志记录器，需要传入一个io.Writer类型的文件对象
//
// 你可以使用os.OpenFile()函数打开一个文件，并将其传入该函数
//
// 例如：
//
//	file, err := os.OpenFile("log.json", os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0666)
//
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	log.InitJsonFileLogger(file, logrus.DebugLevel)
//
// 你可以单独设置json文件日志的日志等级，如果没有传入，则默认使用InfoLevel
//
// 详细的日志等级说明请参考InitTextLogger函数
func (cl *CryoLogger) InitJsonFileLogger(file io.Writer, level ...logrus.Level) {
	if len(level) <= 0 { // 如果没有传入日志等级，则默认使用InfoLevel
		level = append(level, logrus.InfoLevel)
	}

	bufferedFile := bufio.NewWriter(file) // 使用bufio.NewWriter创建一个带缓冲的io.Writer对象

	cl.JsonFileLogger = logrus.New()
	cl.JsonFileLogger.SetFormatter(&logrus.JSONFormatter{})
	cl.JsonFileLogger.SetLevel(level[0])
	cl.JsonFileLogger.SetOutput(bufferedFile)

	cl.updateActiveLoggers() // 更新活跃日志记录器列表
}

// updateActiveLoggers 重新构建活跃日志记录器的切片
// 在日志记录器被初始化或禁用时调用此函数
func (cl *CryoLogger) updateActiveLoggers() {
	// P.S 这个写法是Claude Sonnet 3.7教我的（笑
	// 用于解决高频调用日志时可能出现的性能问题
	// 但是我没有怎么测试过，也许过段时间会有完整的测试

	// 做了个简单的性能测试，答案是有点但是不多，还是尽量不要同时用一堆logger比较好，
	cl.loggersMutex.Lock()
	defer cl.loggersMutex.Unlock()
	// 重置并重建活跃日志记录器列表
	cl.activeLoggers = cl.activeLoggers[:0]
	// 添加每个非空的日志记录器
	if cl.TextLogger != nil {
		cl.activeLoggers = append(cl.activeLoggers, cl.TextLogger)
	}
	if cl.FileLogger != nil {
		cl.activeLoggers = append(cl.activeLoggers, cl.FileLogger)
	}
	if cl.JsonFileLogger != nil {
		cl.activeLoggers = append(cl.activeLoggers, cl.JsonFileLogger)
	}
	if cl.MongoLogger != nil {
		cl.activeLoggers = append(cl.activeLoggers, cl.MongoLogger)
	}
}

// Trace 在所有活跃的日志记录器上以 Trace 级别记录一条消息
func (cl *CryoLogger) Trace(args ...interface{}) {
	cl.loggersMutex.RLock()
	defer cl.loggersMutex.RUnlock()
	for _, logger := range cl.activeLoggers {
		logger.Trace(args...)
	}
}

// Debug 在所有活跃的日志记录器上以 Debug 级别记录一条消息
func (cl *CryoLogger) Debug(args ...interface{}) {
	cl.loggersMutex.RLock()
	defer cl.loggersMutex.RUnlock()
	for _, logger := range cl.activeLoggers {
		logger.Debug(args...)
	}
}

// Info 在所有活跃的日志记录器上以 Info 级别记录一条消息
func (cl *CryoLogger) Info(args ...interface{}) {
	cl.loggersMutex.RLock()
	defer cl.loggersMutex.RUnlock()
	for _, logger := range cl.activeLoggers {
		logger.Info(args...)
	}
}

// Warn 在所有活跃的日志记录器上以 Warn 级别记录一条消息
func (cl *CryoLogger) Warn(args ...interface{}) {
	cl.loggersMutex.RLock()
	defer cl.loggersMutex.RUnlock()
	for _, logger := range cl.activeLoggers {
		logger.Warn(args...)
	}
}

// Error 在所有活跃的日志记录器上以 Error 级别记录一条消息
func (cl *CryoLogger) Error(args ...interface{}) {
	cl.loggersMutex.RLock()
	defer cl.loggersMutex.RUnlock()
	for _, logger := range cl.activeLoggers {
		logger.Error(args...)
	}
}

// Fatal 在所有活跃的日志记录器上以 Fatal 级别记录一条消息
func (cl *CryoLogger) Fatal(args ...interface{}) {
	cl.loggersMutex.RLock()
	defer cl.loggersMutex.RUnlock()
	for _, logger := range cl.activeLoggers {
		logger.Fatal(args...)
	}
}

// Panic 在所有活跃的日志记录器上以 Panic 级别记录一条消息
func (cl *CryoLogger) Panic(args ...interface{}) {
	cl.loggersMutex.RLock()
	defer cl.loggersMutex.RUnlock()
	for _, logger := range cl.activeLoggers {
		logger.Panic(args...)
	}
}

// Tracef 在所有活跃的日志记录器上以 Trace 级别记录格式化的消息
func (cl *CryoLogger) Tracef(format string, args ...interface{}) {
	cl.loggersMutex.RLock()
	defer cl.loggersMutex.RUnlock()
	for _, logger := range cl.activeLoggers {
		logger.Tracef(format, args...)
	}
}

// Debugf 在所有活跃的日志记录器上以 Debug 级别记录格式化的消息
func (cl *CryoLogger) Debugf(format string, args ...interface{}) {
	cl.loggersMutex.RLock()
	defer cl.loggersMutex.RUnlock()
	for _, logger := range cl.activeLoggers {
		logger.Debugf(format, args...)
	}
}

// Infof 在所有活跃的日志记录器上以 Info 级别记录格式化的消息
func (cl *CryoLogger) Infof(format string, args ...interface{}) {
	cl.loggersMutex.RLock()
	defer cl.loggersMutex.RUnlock()
	for _, logger := range cl.activeLoggers {
		logger.Infof(format, args...)
	}
}

// Warnf 在所有活跃的日志记录器上以 Warn 级别记录格式化的消息
func (cl *CryoLogger) Warnf(format string, args ...interface{}) {
	cl.loggersMutex.RLock()
	defer cl.loggersMutex.RUnlock()
	for _, logger := range cl.activeLoggers {
		logger.Warnf(format, args...)
	}
}

// Errorf 在所有活跃的日志记录器上以 Error 级别记录格式化的消息
func (cl *CryoLogger) Errorf(format string, args ...interface{}) {
	cl.loggersMutex.RLock()
	defer cl.loggersMutex.RUnlock()
	for _, logger := range cl.activeLoggers {
		logger.Errorf(format, args...)
	}
}

// Fatalf 在所有活跃的日志记录器上以 Fatal 级别记录格式化的消息
func (cl *CryoLogger) Fatalf(format string, args ...interface{}) {
	cl.loggersMutex.RLock()
	defer cl.loggersMutex.RUnlock()
	for _, logger := range cl.activeLoggers {
		logger.Fatalf(format, args...)
	}
}

// Panicf 在所有活跃的日志记录器上以 Panic 级别记录格式化的消息
func (cl *CryoLogger) Panicf(format string, args ...interface{}) {
	cl.loggersMutex.RLock()
	defer cl.loggersMutex.RUnlock()
	for _, logger := range cl.activeLoggers {
		logger.Panicf(format, args...)
	}
}

/*

以下是一些全局快捷方式，用于更方便的调用Logger

*/

// SetLogger 设置全局日志记录器
func SetLogger(l Logger) {
	logger = l
}

// GetLogger 获取全局日志记录器
func GetLogger() Logger {
	return logger
}

// Trace 在全局日志记录器上以 Trace 级别记录一条消息
func Trace(args ...interface{}) {
	logger.Trace(args...)
}

// Debug 在全局日志记录器上以 Debug 级别记录一条消息
func Debug(args ...interface{}) {
	logger.Debug(args...)
}

// Info 在全局日志记录器上以 Info 级别记录一条消息
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Warn 在全局日志记录器上以 Warn 级别记录一条消息
func Warn(args ...interface{}) {
	logger.Warn(args...)
}

// Error 在全局日志记录器上以 Error 级别记录一条消息
func Error(args ...interface{}) {
	logger.Error(args...)
}

// Fatal 在全局日志记录器上以 Fatal 级别记录一条消息
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

// Panic 在全局日志记录器上以 Panic 级别记录一条消息
func Panic(args ...interface{}) {
	logger.Panic(args...)
}

// Tracef 在全局日志记录器上以 Trace 级别记录格式化的消息
func Tracef(format string, args ...interface{}) {
	logger.Tracef(format, args...)
}

// Debugf 在全局日志记录器上以 Debug 级别记录格式化的消息
func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

// Infof 在全局日志记录器上以 Info 级别记录格式化的消息
func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

// Warnf 在全局日志记录器上以 Warn 级别记录格式化的消息
func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

// Errorf 在全局日志记录器上以 Error 级别记录格式化的消息
func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

// Fatalf 在全局日志记录器上以 Fatal 级别记录格式化的消息
func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

// Panicf 在全局日志记录器上以 Panic 级别记录格式化的消息
func Panicf(format string, args ...interface{}) {
	logger.Panicf(format, args...)
}

// InitTextLogger 初始化终端日志记录器
func InitTextLogger(level ...logrus.Level) {
	// 首先检测logger类型
	// 如果是CryoLogger类型，则直接调用InitTextLogger()函数
	// 如果不是，那么跳过该流程
	if _, ok := logger.(*CryoLogger); ok {
		// 断言logger为CryoLogger类型
		// 然后调用InitTextLogger()函数
		// 这里的logger是全局变量，所以直接调用即可
		logger.(*CryoLogger).InitTextLogger(level...)
	} else {
		// 如果不是CryoLogger类型，那么告知用户该日志类型不支持直接设置内置的几种日志记录器
		logger.Warn("使用了自定义的日志记录器，已跳过默认的终端日志记录器初始化流程")
	}
}

// InitFileLogger 初始化文件日志记录器
func InitFileLogger(file io.Writer, level ...logrus.Level) {
	// 首先检测logger类型
	// 如果是CryoLogger类型，则直接调用InitFileLogger()函数
	// 如果不是，那么跳过该流程
	if _, ok := logger.(*CryoLogger); ok {
		// 断言logger为CryoLogger类型
		// 然后调用InitFileLogger()函数
		// 这里的logger是全局变量，所以直接调用即可
		logger.(*CryoLogger).InitFileLogger(file, level...)
	} else {
		// 如果不是CryoLogger类型，那么告知用户该日志类型不支持直接设置内置的几种日志记录器
		logger.Warn("使用了自定义的日志记录器，已跳过默认的文件日志记录器初始化流程")
	}
}

// InitJsonFileLogger 初始化json文件日志记录器
func InitJsonFileLogger(file io.Writer, level ...logrus.Level) {
	// 首先检测logger类型
	// 如果是CryoLogger类型，则直接调用InitJsonFileLogger()函数
	// 如果不是，那么跳过该流程
	if _, ok := logger.(*CryoLogger); ok {
		// 断言logger为CryoLogger类型
		// 然后调用InitJsonFileLogger()函数
		// 这里的logger是全局变量，所以直接调用即可
		logger.(*CryoLogger).InitJsonFileLogger(file, level...)
	} else {
		// 如果不是CryoLogger类型，那么告知用户该日志类型不支持直接设置内置的几种日志记录器
		logger.Warn("使用了自定义的日志记录器，已跳过默认的json文件日志记录器初始化流程")
	}
}
