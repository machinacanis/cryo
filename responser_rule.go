package cryo

// Rule 是一个泛型规则类型，接受实现了 Event 接口的泛型参数 T
type Rule[T Event] func(event T) bool

// RuleFor 提供一个便捷的泛型方法，用于为特定事件类型创建规则
func RuleFor[T Event](ruleFunc func(T) bool) func(Event) bool {
	return func(e Event) bool {
		// 根据类型进行特殊处理
		switch typedRuleVal := any(ruleFunc).(type) {
		// 处理 UniMessageEvent 类型规则
		case func(*UniMessageEvent) bool:
			var uniMsg *UniMessageEvent
			switch evt := e.(type) {
			case *GroupMessageEvent:
				uniMsg = evt.GetUniMessageEvent()
			case *PrivateMessageEvent:
				uniMsg = evt.GetUniMessageEvent()
			case *TempMessageEvent:
				uniMsg = evt.GetUniMessageEvent()
			default:
				return false // 不是支持的消息事件类型
			}
			if uniMsg == nil {
				return false
			}
			return typedRuleVal(uniMsg)

		// 处理 UniEvent 类型规则
		case func(*UniEvent) bool:
			uniEvent, ok := getUniEventFromEvent(e)
			if !ok || uniEvent == nil {
				return false
			}
			return typedRuleVal(uniEvent)

		// 处理普通事件类型规则
		default:
			// 尝试常规类型断言
			if typedEvent, ok := e.(T); ok {
				return ruleFunc(typedEvent)
			}
			return false // 类型不匹配时规则不适用
		}
	}
}
