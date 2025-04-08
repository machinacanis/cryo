package cryo

import (
	"github.com/machinacanis/cryo/log"
	"reflect"
)

// Handle 使用反射来实现事件处理函数的注册，反正这个方法调用频率不高，不太需要担心性能问题
func (r *OnResponser) Handle(handler interface{}, ordering ...MiddlewareOrdering) *OnResponser {
	var o MiddlewareOrdering
	if len(ordering) == 0 {
		o = AsyncMiddlewareType
	} else {
		o = ordering[0]
	}

	handlerType := reflect.TypeOf(handler)
	if handlerType.Kind() != reflect.Func {
		log.Error("传入的处理函数不是一个函数")
		return r
	}

	// 检查函数是否符合要求的模式
	if handlerType.NumIn() != 1 || !handlerType.In(0).Implements(reflect.TypeOf((*Event)(nil)).Elem()) {
		log.Error("处理函数必须接受一个实现了Event接口的参数")
		return r
	}

	// 检查是否是转换型或者消费型函数
	isTransformer := handlerType.NumOut() == 1 && handlerType.Out(0).AssignableTo(handlerType.In(0))
	isConsumer := handlerType.NumOut() == 0

	if !isTransformer && !isConsumer {
		log.Error("处理函数必须是转换型(func(T) T)或消费型(func(T))")
		return r
	}

	// 动态创建适当的包装器
	var eventWrapper func(Event) Event

	// 特殊处理 UniMessageEvent 类型的处理函数
	paramType := handlerType.In(0)
	typeName := paramType.String()

	switch {
	case typeName == "*cryo.UniMessageEvent" || typeName == "cryo.UniMessageEvent":
		// 处理 UniMessageEvent 类型参数
		if isTransformer {
			eventWrapper = r.createUniMessageOrderdWrapper(handler, r.rules)
		} else {
			eventWrapper = r.createUniMessageWrapper(handler, r.rules)
		}

	case typeName == "*cryo.UniEvent" || typeName == "cryo.UniEvent":
		// 处理 UniEvent 类型参数，适用于所有有 GetUniEvent 方法的事件
		if isTransformer {
			eventWrapper = r.createUniEventOrderdWrapper(handler, r.rules)
		} else {
			eventWrapper = r.createUniEventWrapper(handler, r.rules)
		}

	default:
		// 处理其他普通类型
		if isTransformer {
			eventWrapper = r.createRuleOrderdWrapper(handler, paramType, r.rules)
		} else {
			eventWrapper = r.createRuleWrapper(handler, paramType, r.rules)
		}
	}

	r.AddHandler(eventWrapper, o)
	return r
}

// 创建带规则的OrderdWrapper
func (r *OnResponser) createRuleOrderdWrapper(handler interface{}, eventType reflect.Type, rules []Rule[Event]) func(Event) Event {
	return func(e Event) Event {
		// 类型检查
		if !reflect.TypeOf(e).AssignableTo(eventType) {
			return e
		}

		// 规则检查
		for _, rule := range rules {
			if !rule(e) {
				return e // 如果任何规则返回false，终止处理
			}
		}

		// 所有规则通过，执行处理函数
		args := []reflect.Value{reflect.ValueOf(e)}
		result := reflect.ValueOf(handler).Call(args)
		if len(result) > 0 {
			return result[0].Interface().(Event)
		}
		return e
	}
}

// 创建带规则的Wrapper
func (r *OnResponser) createRuleWrapper(handler interface{}, eventType reflect.Type, rules []Rule[Event]) func(Event) Event {
	return func(e Event) Event {
		// 类型检查
		if !reflect.TypeOf(e).AssignableTo(eventType) {
			return e
		}

		// 规则检查
		for _, rule := range rules {
			if !rule(e) {
				return e // 如果任何规则返回false，终止处理
			}
		}

		// 所有规则通过，执行处理函数
		args := []reflect.Value{reflect.ValueOf(e)}
		reflect.ValueOf(handler).Call(args)
		return e
	}
}

// 创建针对 UniMessageEvent 的特殊 OrderdWrapper
func (r *OnResponser) createUniMessageOrderdWrapper(handler interface{}, rules []Rule[Event]) func(Event) Event {
	return func(e Event) Event {
		// 检查事件是否是三种消息类型之一
		var uniEvent *UniMessageEvent
		var isMessageEvent bool

		switch evt := e.(type) {
		case *GroupMessageEvent:
			uniEvent = evt.GetUniMessageEvent()
			isMessageEvent = true
		case *PrivateMessageEvent:
			uniEvent = evt.GetUniMessageEvent()
			isMessageEvent = true
		case *TempMessageEvent:
			uniEvent = evt.GetUniMessageEvent()
			isMessageEvent = true
		default:
			return e
		}

		if !isMessageEvent {
			return e
		}

		// 规则检查
		for _, rule := range rules {
			if !rule(e) {
				return e // 如果任何规则返回false，终止处理
			}
		}

		// 所有规则通过，执行处理函数
		args := []reflect.Value{reflect.ValueOf(uniEvent)}
		reflect.ValueOf(handler).Call(args)

		// 这里简单返回原始事件
		return e
	}
}

// 创建针对 UniMessageEvent 的特殊 Wrapper
func (r *OnResponser) createUniMessageWrapper(handler interface{}, rules []Rule[Event]) func(Event) Event {
	return func(e Event) Event {
		// 检查事件是否是三种消息类型之一
		var uniEvent *UniMessageEvent
		var isMessageEvent bool

		switch evt := e.(type) {
		case *GroupMessageEvent:
			uniEvent = evt.GetUniMessageEvent()
			isMessageEvent = true
		case *PrivateMessageEvent:
			uniEvent = evt.GetUniMessageEvent()
			isMessageEvent = true
		case *TempMessageEvent:
			uniEvent = evt.GetUniMessageEvent()
			isMessageEvent = true
		default:
			return e
		}

		if !isMessageEvent {
			return e
		}

		// 规则检查
		for _, rule := range rules {
			if !rule(e) {
				return e // 如果任何规则返回false，终止处理
			}
		}

		// 所有规则通过，执行处理函数
		args := []reflect.Value{reflect.ValueOf(uniEvent)}
		reflect.ValueOf(handler).Call(args)
		return e
	}
}

// 创建针对 UniEvent 的特殊 OrderdWrapper
func (r *OnResponser) createUniEventOrderdWrapper(handler interface{}, rules []Rule[Event]) func(Event) Event {
	return func(e Event) Event {
		// 尝试将事件转换为 UniEvent
		uniEvent, ok := getUniEventFromEvent(e)
		if !ok {
			return e
		}

		// 规则检查
		for _, rule := range rules {
			if !rule(e) {
				return e // 如果任何规则返回false，终止处理
			}
		}

		// 所有规则通过，执行处理函数
		args := []reflect.Value{reflect.ValueOf(uniEvent)}
		result := reflect.ValueOf(handler).Call(args)
		if len(result) > 0 {
			// 处理返回值，但由于是通用事件，我们不能直接修改原事件
			// 这里简单返回原始事件
			return e
		}
		return e
	}
}

// 创建针对 UniEvent 的特殊 Wrapper
func (r *OnResponser) createUniEventWrapper(handler interface{}, rules []Rule[Event]) func(Event) Event {
	return func(e Event) Event {
		// 尝试将事件转换为 UniEvent
		uniEvent, ok := getUniEventFromEvent(e)
		if !ok {
			return e
		}

		// 规则检查
		for _, rule := range rules {
			if !rule(e) {
				return e // 如果任何规则返回false，终止处理
			}
		}

		// 所有规则通过，执行处理函数
		args := []reflect.Value{reflect.ValueOf(uniEvent)}
		reflect.ValueOf(handler).Call(args)
		return e
	}
}

// getUniEventFromEvent 尝试从事件中获取 UniEvent
func getUniEventFromEvent(e Event) (*UniEvent, bool) {
	// 获取事件的反射值
	v := reflect.ValueOf(e)

	// 检查是否是指针，如果不是则无法调用方法
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return nil, false
	}

	// 尝试调用 GetUniEvent 方法
	method := v.MethodByName("GetUniEvent")
	if !method.IsValid() {
		return nil, false
	}

	// 调用方法获取 UniEvent
	results := method.Call(nil)
	if len(results) != 1 {
		return nil, false
	}

	// 尝试类型断言
	uniEvent, ok := results[0].Interface().(*UniEvent)
	return uniEvent, ok
}
