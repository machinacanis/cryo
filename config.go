package cryo

import (
	"github.com/go-json-experiment/json"
	"github.com/sirupsen/logrus"
	"os"
)

var conf Config
var DefaultSignServer = "https://sign.lagrangecore.org/api/sign/30366" // 默认的签名服务器地址

// Config cryo 的配置项，通过在 Bot.Init() 中传入来控制每个Bot实例的功能
type Config struct {
	LogLevel                     logrus.Level
	LogFormat                    *logrus.Formatter
	SignServers                  []string `json:"sign_servers,omitempty,omitzero"`                    // 签名服务器列表
	EnablePluginAutoLoad         bool     `json:"enable_plugin_auto_load,omitempty,omitzero"`         // 是否启用插件自动加载（自动启用插件）
	EnableClientAutoSave         bool     `json:"enable_client_save,omitempty,omitzero"`              // 是否启用客户端信息自动保存
	EnablePrintLogo              bool     `json:"enable_print_logo,omitempty,omitzero"`               // 是否启用logo打印
	EnableConnectPrintMiddleware bool     `json:"enable_connect_print_middleware,omitempty,omitzero"` // 是否启用内置的Bot连接打印中间件
	EnableMessagePrintMiddleware bool     `json:"enable_message_print_middleware,omitempty,omitzero"` // 是否启用内置的消息打印中间件
	EnableEventDebugMiddleware   bool     `json:"enable_event_debug_middleware,omitempty,omitzero"`   // 是否启用内置的事件调试中间件
}

// ReadCryoConfig 从文件读取配置项
func ReadCryoConfig() (Config, error) {
	data, err := os.ReadFile("cryo_config.json")
	c := Config{}
	if err != nil {
		return c, err
	}

	err = json.Unmarshal(data, &c)
	if err != nil {
		return c, err
	}

	return c, nil
}

// WriteCryoConfig 写入配置项到文件
func WriteCryoConfig(config Config) error {
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	err = os.WriteFile("cryo_config.json", data, 0644)
	if err != nil {
		return err
	}

	return nil
}
