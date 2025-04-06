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

var PLogger = ProtocolLogger{}

type ProtocolLogger struct{}

var fromProtocol = deepskyblue + "[Lagrange] "

func (p ProtocolLogger) Info(format string, arg ...any) {
	Infof(fromProtocol+format, arg...)
}

func (p ProtocolLogger) Warning(format string, arg ...any) {
	Warnf(fromProtocol+format, arg...)
}

func (p ProtocolLogger) Debug(format string, arg ...any) {
	Debugf(fromProtocol+format, arg...)
}

func (p ProtocolLogger) Error(format string, arg ...any) {
	Errorf(fromProtocol+format, arg...)
}

// Dump 输出当前日志记录器的状态
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
