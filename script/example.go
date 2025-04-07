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

	cryo.NewOnResponser(bot.GetBus(), cryo.PrivateMessageEventType, cryo.GroupMessageEventType).
		Handle(func(e *cryo.PrivateMessageEvent) {
			log.Info("这是一条私聊消息！")
		}, cryo.AsyncMiddlewareType).
		Handle(func(e *cryo.GroupMessageEvent) {
			log.Info("这是一条群消息！")
		}, cryo.AsyncMiddlewareType).
		Register()

	bot.AutoConnect()
	bot.Start()
}
