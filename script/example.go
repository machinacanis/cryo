package main

import (
	"github.com/machinacanis/cryo"
	"github.com/machinacanis/cryo/log"
)

func main() {
	logger := log.NewLoggerBuilder().AddTextLogger(log.InfoLevel)
	bot := cryo.NewBot()
	bot.Init(logger, cryo.Config{
		EnableMessagePrintMiddleware: true,
		EnableEventDebugMiddleware:   true,
	})

	bot.OnType(cryo.PrivateMessageEventType, cryo.GroupMessageEventType).
		Handle(func(e *cryo.PrivateMessageEvent) {
			logger.Info("响应事件 " + e.EventId)
			// 自定义逻辑
		}, cryo.AsyncMiddlewareType).
		Register()

	bot.AutoConnect()
	bot.Start()
}
