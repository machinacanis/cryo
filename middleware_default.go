package cryo

import "github.com/machinacanis/cryo/log"

// setDefaultMiddleware 设置默认的中间件
func setDefaultMiddleware(bus *EventBus, conf *Config) {
	if conf.EnableEventDebugMiddleware {
		log.Debug("启用内置的事件调试中间件")
		mw := NewUniMiddleware()
		mw.AddHandler(func(e Event) Event {
			u := e.GetUniEvent()
			log.Debugf("[EventPublish] %s from %s(%d) with Id %s and Tags %v", u.GetEventType().ToString(), u.BotNickname, u.BotUin, u.BotId, u.EventTags)
			return e
		})
		bus.AddPreMiddleware(mw)
	}
}
