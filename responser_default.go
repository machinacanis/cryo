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

// OnMessageToMe 创建一个新的消息事件响应器
//
// 这个响应器在 OnMessage 的基础上添加了Bot被At时的响应规则，你可以选择是否去除掉At元素，默认是去除的
func (b *Bot) OnMessageToMe(removeAt ...bool) *OnResponser {
	return NewOnResponser(b.bus, PrivateMessageEventType, GroupMessageEventType, TempMessageEventType).AddRule(ToMeRule(removeAt...)) // 使用内置的规则
}

// OnMessageToSomeOne 创建一个新的消息事件响应器
//
// 这个响应器在 OnMessage 的基础上添加了指定用户被At时的响应规则
func (b *Bot) OnMessageToSomeOne(target ...uint32) *OnResponser {
	return NewOnResponser(b.bus, PrivateMessageEventType, GroupMessageEventType, TempMessageEventType).AddRule(AtRule(target...)) // 使用内置的规则
}

// OnStartWith 创建一个新的消息事件响应器
//
// 这个响应器在 OnMessage 的基础上添加了以指定文本开头的响应规则
func (b *Bot) OnStartWith(prefix ...string) *OnResponser {
	return NewOnResponser(b.bus, PrivateMessageEventType, GroupMessageEventType, TempMessageEventType).AddRule(StartWithRule(prefix...)) // 使用内置的规则
}

// OnEndWith 创建一个新的消息事件响应器
//
// 这个响应器在 OnMessage 的基础上添加了以指定文本结尾的响应规则
func (b *Bot) OnEndWith(suffix ...string) *OnResponser {
	return NewOnResponser(b.bus, PrivateMessageEventType, GroupMessageEventType, TempMessageEventType).AddRule(EndWithRule(suffix...)) // 使用内置的规则
}

// OnFullMatch 创建一个新的消息事件响应器
//
// 这个响应器在 OnMessage 的基础上添加了文本内容完全匹配的响应规则
func (b *Bot) OnFullMatch(content ...string) *OnResponser {
	return NewOnResponser(b.bus, PrivateMessageEventType, GroupMessageEventType, TempMessageEventType).AddRule(FullMatchRule(content...)) // 使用内置的规则
}

// OnKeyWord 创建一个新的消息事件响应器
//
// 这个响应器在 OnMessage 的基础上添加了文本内容包含指定关键字的响应规则
func (b *Bot) OnKeyWord(keyword ...string) *OnResponser {
	return NewOnResponser(b.bus, PrivateMessageEventType, GroupMessageEventType, TempMessageEventType).AddRule(KeyWordRule(keyword...)) // 使用内置的规则
}

// OnAllKeyWord 创建一个新的消息事件响应器
//
// 这个响应器在 OnMessage 的基础上添加了文本内容包含所有指定关键字的响应规则
func (b *Bot) OnAllKeyWord(keyword ...string) *OnResponser {
	return NewOnResponser(b.bus, PrivateMessageEventType, GroupMessageEventType, TempMessageEventType).AddRule(AllKeyWordRule(keyword...)) // 使用内置的规则
}
