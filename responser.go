package cryo

import "github.com/machinacanis/cryo/log"

// Wrapper 带泛型的事件处理函数包装器
func Wrapper[T Event](handler func(T)) func(Event) Event {
	return func(e Event) Event {
		if evt, ok := e.(T); ok {
			handler(evt)
		}
		return nil
	}
}

// OrderdWrapper 带泛型的事件处理函数包装器，支持中间件的顺序执行
func OrderdWrapper[T Event](handler func(T) T) func(Event) Event {
	return func(e Event) Event {
		if evt, ok := e.(T); ok {
			evt = handler(evt)
			return evt
		}
		return nil // 返回 nil 表示事件被截断
	}
}

// Responser 事件响应器接口
type Responser interface {
	GetId() string                               // 获取响应器的唯一标识符
	AddType(eventType ...EventType) Responser    // 添加响应器响应的事件类型
	GetType() []EventType                        // 获取响应器响应的事件类型列表
	RemoveType(eventType ...EventType) Responser // 移除响应器响应的事件类型
	Response(func(Event) Event) Responser        // 处理事件
	Register()                                   // 注册响应器
	Remove()                                     // 移除响应器注册的所有中间件
}

// UniResponser 是一个基础的事件响应器实现
type UniResponser struct {
	id          string       // 响应器的唯一标识符，是一个uuid，可以通过这个id来区分不同的响应器
	bus         *EventBus    // 事件总线，用于管理中间件和事件的分发
	eventType   []EventType  // 响应器响应的事件类型列表，这个响应器构建的中间件都会获得这些事件类型
	middlewares []Middleware // 响应器的中间件列表
}

// GetId 获取响应器的唯一标识符
func (r *UniResponser) GetId() string {
	return r.id
}

// AddType 添加响应器响应的事件类型
func (r *UniResponser) AddType(eventType ...EventType) Responser {
	r.eventType = append(r.eventType, eventType...)
	return r
}

// GetType 获取响应器响应的事件类型列表
func (r *UniResponser) GetType() []EventType {
	return r.eventType
}

// RemoveType 移除响应器响应的事件类型
func (r *UniResponser) RemoveType(eventType ...EventType) Responser {
	for _, et := range eventType {
		for i, v := range r.eventType {
			if v == et {
				r.eventType = append(r.eventType[:i], r.eventType[i+1:]...)
				break
			}
		}
	}
	return r
}

// Response 响应
func (r *UniResponser) Response(handler func(event Event) Event) Responser {
	mw := NewUniMiddleware()
	mw.AddHandler(handler)
	r.middlewares = append(r.middlewares, mw)
	return r
}

// Register 注册响应器
func (r *UniResponser) Register() {
	// 将响应器的响应事件类型注入到中间件中
	for _, et := range r.eventType {
		for _, mw := range r.middlewares {
			mw.AddType(et)
		}
	}
	// 注册响应器的中间件到异步事件总线
	// 也就是说这些中间件是完全异步互相不干扰的
	r.bus.AddAsyncMiddleware(r.middlewares...)
}

// Remove 移除响应器注册的所有中间件
func (r *UniResponser) Remove() {
	// 获取所有中间件的id
	middlewareIds := make([]string, 0)
	for _, mw := range r.middlewares {
		middlewareIds = append(middlewareIds, mw.GetId())
	}
	// 移除所有中间件
	r.bus.RemoveMiddlewareById(middlewareIds...)
}

// NewUniResponser 创建一个新的基础响应器
func NewUniResponser(bus *EventBus, eventType ...EventType) *UniResponser {
	// 创建一个新的响应器
	return &UniResponser{
		id:          newUUID(),
		bus:         bus,
		eventType:   eventType,
		middlewares: make([]Middleware, 0),
	}
}

// OnResponser 默认事件响应器实现
type OnResponser struct {
	UniResponser
	preMiddleware   Middleware // 预处理中间件列表
	postMiddleware  Middleware // 后处理中间件列表
	asyncMiddleware Middleware // 异步处理中间件列表
	syncMiddleware  Middleware // 同步处理中间件列表
}

func (r *OnResponser) AddHandler(handler EventHandler[Event], ordering MiddlewareOrdering) {
	switch ordering {
	case PreMiddlewareType:
		r.preMiddleware.AddHandler(handler)
	case PostMiddlewareType:
		r.postMiddleware.AddHandler(handler)
	case AsyncMiddlewareType:
		r.asyncMiddleware.AddHandler(handler)
	case SyncMiddlewareType:
		r.syncMiddleware.AddHandler(handler)
	default:
		r.asyncMiddleware.AddHandler(handler)
	}
}

func (r *OnResponser) Handle(handler interface{}, ordering ...MiddlewareOrdering) *OnResponser {
	var o MiddlewareOrdering
	if len(ordering) == 0 {
		o = SyncMiddlewareType
	} else {
		o = ordering[0]
	}

	switch typedHandler := handler.(type) {
	case func(event *PrivateMessageEvent):
		w := Wrapper(typedHandler)
		r.AddHandler(w, o)
	case func(event *PrivateMessageEvent) *PrivateMessageEvent:
		w := OrderdWrapper(typedHandler)
		r.AddHandler(w, o)
	case func(event *GroupMessageEvent):
		w := Wrapper(typedHandler)
		r.AddHandler(w, o)
	case func(event *GroupMessageEvent) *GroupMessageEvent:
		w := OrderdWrapper(typedHandler)
		r.AddHandler(w, o)
	case func(event *TempMessageEvent):
		w := Wrapper(typedHandler)
		r.AddHandler(w, o)
	case func(event *TempMessageEvent) *TempMessageEvent:
		w := OrderdWrapper(typedHandler)
		r.AddHandler(w, o)
	case func(event *NewFriendRequestEvent):
		w := Wrapper(typedHandler)
		r.AddHandler(w, o)
	case func(event *NewFriendRequestEvent) *NewFriendRequestEvent:
		w := OrderdWrapper(typedHandler)
		r.AddHandler(w, o)
	case func(event *NewFriendEvent):
		w := Wrapper(typedHandler)
		r.AddHandler(w, o)
	case func(event *NewFriendEvent) *NewFriendEvent:
		w := OrderdWrapper(typedHandler)
		r.AddHandler(w, o)
	case func(event *FriendRecallEvent):
		w := Wrapper(typedHandler)
		r.AddHandler(w, o)
	case func(event *FriendRecallEvent) *FriendRecallEvent:
		w := OrderdWrapper(typedHandler)
		r.AddHandler(w, o)
	case func(event *FriendRenameEvent):
		w := Wrapper(typedHandler)
		r.AddHandler(w, o)
	case func(event *FriendRenameEvent) *FriendRenameEvent:
		w := OrderdWrapper(typedHandler)
		r.AddHandler(w, o)
	case func(event *FriendPokeEvent):
		w := Wrapper(typedHandler)
		r.AddHandler(w, o)
	case func(event *FriendPokeEvent) *FriendPokeEvent:
		w := OrderdWrapper(typedHandler)
		r.AddHandler(w, o)
	case func(event *GroupMemberPermissionUpdatedEvent):
		w := Wrapper(typedHandler)
		r.AddHandler(w, o)
	case func(event *GroupMemberPermissionUpdatedEvent) *GroupMemberPermissionUpdatedEvent:
		w := OrderdWrapper(typedHandler)
		r.AddHandler(w, o)
	case func(event *GroupNameUpdatedEvent):
		w := Wrapper(typedHandler)
		r.AddHandler(w, o)
	case func(event *GroupNameUpdatedEvent) *GroupNameUpdatedEvent:
		w := OrderdWrapper(typedHandler)
		r.AddHandler(w, o)
	case func(event *GroupMuteEvent):
		w := Wrapper(typedHandler)
		r.AddHandler(w, o)
	case func(event *GroupMuteEvent) *GroupMuteEvent:
		w := OrderdWrapper(typedHandler)
		r.AddHandler(w, o)
	case func(event *GroupRecallEvent):
		w := Wrapper(typedHandler)
		r.AddHandler(w, o)
	case func(event *GroupRecallEvent) *GroupRecallEvent:
		w := OrderdWrapper(typedHandler)
		r.AddHandler(w, o)
	case func(event *GroupMemberJoinRequestEvent):
		w := Wrapper(typedHandler)
		r.AddHandler(w, o)
	case func(event *GroupMemberJoinRequestEvent) *GroupMemberJoinRequestEvent:
		w := OrderdWrapper(typedHandler)
		r.AddHandler(w, o)
	case func(event *GroupInviteEvent):
		w := Wrapper(typedHandler)
		r.AddHandler(w, o)
	case func(event *GroupInviteEvent) *GroupInviteEvent:
		w := OrderdWrapper(typedHandler)
		r.AddHandler(w, o)
	case func(event *GroupMemberDecreaseEvent):
		w := Wrapper(typedHandler)
		r.AddHandler(w, o)
	case func(event *GroupMemberDecreaseEvent) *GroupMemberDecreaseEvent:
		w := OrderdWrapper(typedHandler)
		r.AddHandler(w, o)
	case func(event *GroupDigestEvent):
		w := Wrapper(typedHandler)
		r.AddHandler(w, o)
	case func(event *GroupDigestEvent) *GroupDigestEvent:
		w := OrderdWrapper(typedHandler)
		r.AddHandler(w, o)
	case func(event *GroupReactionEvent):
		w := Wrapper(typedHandler)
		r.AddHandler(w, o)
	case func(event *GroupReactionEvent) *GroupReactionEvent:
		w := OrderdWrapper(typedHandler)
		r.AddHandler(w, o)
	case func(event *GroupMemberSpecialTitleUpdated):
		w := Wrapper(typedHandler)
		r.AddHandler(w, o)
	case func(event *GroupMemberSpecialTitleUpdated) *GroupMemberSpecialTitleUpdated:
		w := OrderdWrapper(typedHandler)
		r.AddHandler(w, o)
	case func(event *BotConnectedEvent):
		w := Wrapper(typedHandler)
		r.AddHandler(w, o)
	case func(event *BotConnectedEvent) *BotConnectedEvent:
		w := OrderdWrapper(typedHandler)
		r.AddHandler(w, o)
	case func(event *BotDisconnectedEvent):
		w := Wrapper(typedHandler)
		r.AddHandler(w, o)
	case func(event *BotDisconnectedEvent) *BotDisconnectedEvent:
		w := OrderdWrapper(typedHandler)
		r.AddHandler(w, o)
	case func(event *CustomEvent):
		w := Wrapper(typedHandler)
		r.AddHandler(w, o)
	case func(event *CustomEvent) *CustomEvent:
		w := OrderdWrapper(typedHandler)
		r.AddHandler(w, o)
	default:
		log.Error("传入了不支持的事件处理函数")
	}
	return r
}

func (r *OnResponser) Register() {
	// 将响应器的响应事件类型注入到中间件中
	for _, et := range r.eventType {
		if r.preMiddleware.GetHandlerCount() > 0 {
			r.preMiddleware.AddType(et)
		}
		if r.postMiddleware.GetHandlerCount() > 0 {
			r.postMiddleware.AddType(et)
		}
		if r.asyncMiddleware.GetHandlerCount() > 0 {
			r.asyncMiddleware.AddType(et)
		}
		if r.syncMiddleware.GetHandlerCount() > 0 {
			r.syncMiddleware.AddType(et)
		}
	}
	// 注册响应器的中间件
	if r.preMiddleware.GetHandlerCount() > 0 {
		r.bus.AddPreMiddleware(r.preMiddleware)
	}
	if r.postMiddleware.GetHandlerCount() > 0 {
		r.bus.AddPostMiddleware(r.postMiddleware)
	}
	if r.asyncMiddleware.GetHandlerCount() > 0 {
		r.bus.AddAsyncMiddleware(r.asyncMiddleware)
	}
	if r.syncMiddleware.GetHandlerCount() > 0 {
		r.bus.AddSyncMiddleware(r.syncMiddleware)
	}
}

func (r *OnResponser) Remove() {
	// 获取所有中间件的id
	middlewareIds := make([]string, 0)
	for _, mw := range []Middleware{r.preMiddleware, r.postMiddleware, r.asyncMiddleware, r.syncMiddleware} {
		middlewareIds = append(middlewareIds, mw.GetId())
	}
	// 移除所有中间件
	r.bus.RemoveMiddlewareById(middlewareIds...)
}

// NewOnResponser 创建一个新的事件响应器
func NewOnResponser(bus *EventBus, eventType ...EventType) *OnResponser {
	// 创建一个新的响应器
	return &OnResponser{
		UniResponser: UniResponser{
			id:        newUUID(),
			bus:       bus,
			eventType: eventType,
		},
		preMiddleware:   NewUniMiddleware(),
		postMiddleware:  NewUniMiddleware(),
		asyncMiddleware: NewUniMiddleware(),
		syncMiddleware:  NewUniMiddleware(),
	}
}
