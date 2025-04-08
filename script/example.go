package main

import (
	"github.com/machinacanis/cryo"
	"github.com/machinacanis/cryo/log"
	"github.com/sirupsen/logrus"
)

func main() {
	bot := cryo.NewBot()
	bot.Init(cryo.Config{
		LogLevel:                     logrus.DebugLevel,
		EnableMessagePrintMiddleware: true,
		EnableEventDebugMiddleware:   true,
	})

	bot.OnType(cryo.PrivateMessageEventType, cryo.GroupMessageEventType).
		Handle(func(e *cryo.PrivateMessageEvent) {
			log.Info("响应事件 " + e.EventId)
			// 自定义逻辑
		}, cryo.AsyncMiddlewareType).
		Register()

	bot.AutoConnect()
	bot.Start()
}
