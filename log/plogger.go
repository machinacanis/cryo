package log

import (
	"fmt"
	"os"
	"path"
	"time"
)

// from https://github.com/ExquisiteCore/LagrangeGo-Template/blob/main/utils/log.go
//
// 基于LagrangeGo-Template的协议Logger修改而来，使其可以直接兼容本项目的Logger

var dumpspath = "dump"

// PLogger 声明一个给LagrangeGo使用的协议日志记录器
var PLogger = ProtocolLogger{}

// ProtocolLogger 协议日志记录器结构体
type ProtocolLogger struct{}

var fromProtocol = deepskyblue + "[Lagrange] "

// Info 协议日志记录器的Info方法，实际上是调用Logger的Infof方法
func (p ProtocolLogger) Info(format string, arg ...any) {
	Infof(fromProtocol+format, arg...)
}

// Warning 协议日志记录器的Warn方法，实际上是调用Logger的Warnf方法
func (p ProtocolLogger) Warning(format string, arg ...any) {
	Warnf(fromProtocol+format, arg...)
}

// Debug 协议日志记录器的Debug方法，实际上是调用Logger的Debugf方法
func (p ProtocolLogger) Debug(format string, arg ...any) {
	Debugf(fromProtocol+format, arg...)
}

// Error 协议日志记录器的Error方法，实际上是调用Logger的Errorf方法
func (p ProtocolLogger) Error(format string, arg ...any) {
	Errorf(fromProtocol+format, arg...)
}

// Dump 转储数据到文件
func (p ProtocolLogger) Dump(data []byte, format string, arg ...any) {
	message := fmt.Sprintf(format, arg...)
	if _, err := os.Stat(dumpspath); err != nil {
		err = os.MkdirAll(dumpspath, 0o755)
		if err != nil {
			Errorf("出现错误 %v. 详细信息转储失败", message)
			return
		}
	}
	dumpFile := path.Join(dumpspath, fmt.Sprintf("%v.dump", time.Now().Unix()))
	Errorf("出现错误 %v. 详细信息已转储至文件 %v 请连同日志提交给开发者处理", message, dumpFile)
	_ = os.WriteFile(dumpFile, data, 0o644)
}
