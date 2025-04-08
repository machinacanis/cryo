package cryo

import "github.com/machinacanis/cryo/log"

// HandleV1 Deprecated: 该函数已被使用反射实现的处理器替代，请不要使用
func (r *OnResponser) HandleV1(handler interface{}, ordering ...MiddlewareOrdering) *OnResponser {
	var o MiddlewareOrdering
	if len(ordering) == 0 {
		o = AsyncMiddlewareType
	} else {
		o = ordering[0]
	}

	switch typedHandler := handler.(type) {
	// 通过匹配传入的函数的类型来实现依赖注入
	// Well, 我承认这段代码有点丑陋
	// 但是这样可以避免使用反射
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
