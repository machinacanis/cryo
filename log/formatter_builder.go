package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

// DefaultFormatterBuilder 是一个日志格式构建器，通过这个构建器可以快速构建一套适配cryo的默认日志格式
type DefaultFormatterBuilder struct {
	DebugLabel        string
	InfoLabel         string
	SuccessLabel      string
	WarnLabel         string
	ErrorLabel        string
	FatalLabel        string
	PanicLabel        string
	DebugLabelColor   string
	InfoLabelColor    string
	SuccessLabelColor string
	WarnLabelColor    string
	ErrorLabelColor   string
	PanicLabelColor   string
	FatalLabelColor   string
	DebugTextColor    string
	InfoTextColor     string
	SuccessTextColor  string
	WarnTextColor     string
	ErrorTextColor    string
	FatalTextColor    string
	PanicTextColor    string

	formatFunc func(entry *logrus.Entry) ([]byte, error)
}

// GetLabel 根据日志级别获取对应的标签和颜色
func (f DefaultFormatterBuilder) GetLabel(level logrus.Level) (string, string, string) {
	switch level {
	case logrus.TraceLevel:
		return f.DebugLabel, f.DebugLabelColor, f.DebugTextColor
	case logrus.DebugLevel:
		return f.InfoLabel, f.InfoLabelColor, f.InfoTextColor
	case logrus.InfoLevel:
		return f.SuccessLabel, f.SuccessLabelColor, f.SuccessTextColor
	case logrus.WarnLevel:
		return f.WarnLabel, f.WarnLabelColor, f.WarnTextColor
	case logrus.ErrorLevel:
		return f.ErrorLabel, f.ErrorLabelColor, f.ErrorTextColor
	case logrus.FatalLevel:
		return f.FatalLabel, f.FatalLabelColor, f.FatalTextColor
	case logrus.PanicLevel:
		return f.PanicLabel, f.PanicLabelColor, f.PanicTextColor
	}
	return f.InfoLabel, f.InfoLabelColor, f.InfoTextColor
}

// Format 实现logrus.Formatter接口的Format方法
func (f DefaultFormatterBuilder) Format(entry *logrus.Entry) ([]byte, error) {
	label, labelColor, textColor := f.GetLabel(entry.Level)
	logMsg := fmt.Sprintf(
		"%s%s [%s%s%s] %s%s%s\n",
		gray,
		entry.Time.Format("01-02 15:04:05"),
		labelColor,
		label,
		gray,
		textColor,
		entry.Message,
		reset,
	)
	return []byte(logMsg), nil
}

// NewDefaultFormatterBuilder 创建一个新的默认格式构建器
func NewDefaultFormatterBuilder() DefaultFormatterBuilder {
	return DefaultFormatterBuilder{
		DebugLabel:        "🔍__DEBUG",
		InfoLabel:         "🧊___INFO",
		SuccessLabel:      "✅SUCCESS",
		WarnLabel:         "⚠️WARNING",
		ErrorLabel:        "⛔__ERROR",
		FatalLabel:        "💀__FATAL",
		PanicLabel:        "🏴‍☠️__PANIC",
		DebugLabelColor:   lightcyan,
		InfoLabelColor:    cyan,
		SuccessLabelColor: green,
		WarnLabelColor:    yellow,
		ErrorLabelColor:   red,
		PanicLabelColor:   deepred,
		FatalLabelColor:   purple,
		DebugTextColor:    lightcyan,
		InfoTextColor:     white,
		SuccessTextColor:  white,
		WarnTextColor:     yellow,
		ErrorTextColor:    red,
		FatalTextColor:    deepred,
		PanicTextColor:    purple,
		formatFunc:        nil,
	}
}
