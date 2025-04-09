package cryo

// OnType 可以创建一个新的事件类型响应器
//
// 示例：
//
//	 bot.OnType(cryo.PrivateMessageEventType, cryo.GroupMessageEventType).
//		 Handle(func(e *cryo.UniMessageEvent) {
//			 log.Info("响应事件 " + e.EventId)
//				// 自定义逻辑
//		 }).
//		 Register()
//
// 在没有传入事件类型的情况下，它会响应任何类型的事件
func (b *Bot) OnType(eventType ...EventType) *OnResponser {
	return NewOnResponser(b.bus, eventType...)
}

// OnMessage 创建一个新的消息事件响应器
//
// 这个响应器和 OnType 类似，但是它只会响应消息事件，即 PrivateMessageEventType 、 GroupMessageEventType 和 TempMessageEventType
//
// 你可以通过继续在 Handle 中传入 func(e *UniMessageEvent) 来统一处理消息事件，也可以分开匹配一个具体类型的事件
func (b *Bot) OnMessage() *OnResponser {
	return NewOnResponser(b.bus, PrivateMessageEventType, GroupMessageEventType, TempMessageEventType)
}
