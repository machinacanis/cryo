package cryo

import "github.com/machinacanis/cryo/log"

// setDefaultMiddleware 设置默认的中间件
//
// 目前提供了以下中间件：
//
// 1. Bot连接状态打印中间件
//
// 2. 消息打印中间件
//
// 3. 事件调试中间件
func setDefaultMiddleware(bus *EventBus, logger log.CryoLogger, conf Config) {
	if conf.EnableConnectPrintMiddleware { // 是否启用连接状态打印中间件
		logger.Debug("[Cryo] 启用内置的Bot连接状态打印中间件")
		mw1 := NewUniMiddleware(BotConnectedEventType)
		mw2 := NewUniMiddleware(BotDisconnectedEventType)
		mw1.AddHandler(func(e Event) Event {
			if typedEvent, ok := e.(*BotConnectedEvent); ok {
				logger.Infof("[Cryo] %s：%s (%d) 已成功连接", typedEvent.ClientNickname, typedEvent.ClientId, typedEvent.ClientUin)
			}
			return e
		})
		mw2.AddHandler(func(e Event) Event {
			if typedEvent, ok := e.(*BotDisconnectedEvent); ok {
				logger.Infof("[Cryo] %s：%s (%d) 已断开连接", typedEvent.ClientNickname, typedEvent.ClientId, typedEvent.ClientUin)
			}
			return e
		})
		bus.AddPreMiddleware(mw1)
		bus.AddPreMiddleware(mw2)
	}

	if conf.EnableMessagePrintMiddleware { // 是否启用消息打印中间件
		logger.Debug("[Cryo] 启用内置的消息打印中间件")
		mw1 := NewUniMiddleware(PrivateMessageEventType)
		mw2 := NewUniMiddleware(GroupMessageEventType)
		mw3 := NewUniMiddleware(TempMessageEventType)
		mw1.AddHandler(func(e Event) Event {
			if typedEvent, ok := e.(*PrivateMessageEvent); ok {
				logger.Infof("[%s] [私聊] From %s(%d) - %s", typedEvent.ClientNickname, typedEvent.SenderNickname, typedEvent.SenderUin, typedEvent.MessageElements.ToString())
			}
			return e
		})
		mw2.AddHandler(func(e Event) Event {
			if typedEvent, ok := e.(*GroupMessageEvent); ok {
				logger.Infof("[%s] [群聊] [%s(%d)] From %s(%d) - %s", typedEvent.ClientNickname, typedEvent.GroupName, typedEvent.GroupUin, typedEvent.SenderNickname, typedEvent.SenderUin, typedEvent.MessageElements.ToString())
			}
			return e
		})
		mw3.AddHandler(func(e Event) Event {
			if typedEvent, ok := e.(*TempMessageEvent); ok {
				logger.Infof("[%s] [临时] From %s(%d) - %s", typedEvent.ClientNickname, typedEvent.SenderNickname, typedEvent.SenderUin, typedEvent.MessageElements.ToString())
			}
			return e
		})
		bus.AddPreMiddleware(mw1)
		bus.AddPreMiddleware(mw2)
		bus.AddPreMiddleware(mw3)
	}

	if conf.EnableEventDebugMiddleware { // 是否启用事件调试中间件
		logger.Debug("[Cryo] 启用内置的事件调试中间件")
		mw := NewUniMiddleware()
		mw.AddHandler(func(e Event) Event {
			u := e.GetUniEvent()
			logger.Debugf("[EventPublish] %s from %s(%d) with Id %s and Tags %v", u.GetEventType().ToString(), u.ClientNickname, u.ClientUin, u.EventId, u.EventTags)
			return e
		})
		bus.AddPreMiddleware(mw)
	}
}
