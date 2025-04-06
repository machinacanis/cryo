package cryo

// EventType 事件类型
type EventType uint32

const (
	UniEventType                            EventType = iota // 基础事件类型
	UniMessageEventType                                      // 消息事件类型
	PrivateMessageEventType                                  // 私聊消息事件类型
	GroupMessageEventType                                    // 群消息事件类型
	TempMessageEventType                                     // 临时消息事件类型
	NewFriendRequestEventType                                // 新好友请求事件类型
	NewFriendEventType                                       // 新好友事件类型
	FriendRecallEventType                                    // 好友撤回事件类型
	FriendRenameEventType                                    // 好友改名事件类型
	FriendPokeEventType                                      // 好友戳一戳事件类型
	GroupMemberPermissionUpdatedEventType                    // 群成员权限变更事件类型
	GroupNameUpdatedEventType                                // 群名称变更事件类型
	GroupMuteEventType                                       // 群禁言事件类型
	GroupRecallEventType                                     // 群撤回事件类型
	GroupMemberJoinRequestEventType                          // 群成员入群请求事件类型
	GroupMemberIncreaseEventType                             // 群成员增加事件类型
	GroupMemberDecreaseEventType                             // 群成员减少事件类型
	GroupDigestEventType                                     // 群精华消息事件类型
	GroupReactionEventType                                   // 群消息表情事件类型
	GroupMemberSpecialTitleUpdatedEventType                  // 群成员特殊头衔变更事件类型
	GroupInviteEventType                                     // 加群邀请事件类型
	BotConnectedEventType                                    // 机器人连接事件类型
	BotDisconnectedEventType                                 // 机器人断开连接事件类型
	CustomEventType                                          // 自定义事件类型
)

// ToString 输出事件类型的字符串表示
func (et EventType) ToString() string {
	switch et {
	case UniEventType:
		return "UniEvent"
	case UniMessageEventType:
		return "UniMessageEvent"
	case PrivateMessageEventType:
		return "PrivateMessageEvent"
	case GroupMessageEventType:
		return "GroupMessageEvent"
	case TempMessageEventType:
		return "TempMessageEvent"
	case NewFriendRequestEventType:
		return "NewFriendRequestEvent"
	case NewFriendEventType:
		return "NewFriendEvent"
	case FriendRecallEventType:
		return "FriendRecallEvent"
	case FriendRenameEventType:
		return "FriendRenameEvent"
	case FriendPokeEventType:
		return "FriendPokeEvent"
	case GroupMemberPermissionUpdatedEventType:
		return "GroupMemberPermissionUpdatedEvent"
	case GroupNameUpdatedEventType:
		return "GroupNameUpdatedEvent"
	case GroupMuteEventType:
		return "GroupMuteEvent"
	case GroupRecallEventType:
		return "GroupRecallEvent"
	case GroupMemberJoinRequestEventType:
		return "GroupMemberJoinRequestEvent"
	case GroupMemberIncreaseEventType:
		return "GroupMemberIncreaseEvent"
	case GroupMemberDecreaseEventType:
		return "GroupMemberDecreaseEvent"
	case GroupDigestEventType:
		return "GroupDigestEvent"
	case GroupReactionEventType:
		return "GroupReactionEvent"
	case GroupMemberSpecialTitleUpdatedEventType:
		return "GroupMemberSpecialTitleUpdatedEvent"
	case GroupInviteEventType:
		return "GroupInviteEvent"
	case BotConnectedEventType:
		return "BotConnectedEvent"
	case BotDisconnectedEventType:
		return "BotDisconnectedEvent"
	case CustomEventType:
		return "CustomEventType"
	default:
		return "UnknownEventType"
	}
}

// AllEventTypes 返回所有可用的事件类型
func AllEventTypes() []EventType {
	return []EventType{
		UniEventType,
		UniMessageEventType,
		PrivateMessageEventType,
		GroupMessageEventType,
		TempMessageEventType,
		NewFriendRequestEventType,
		NewFriendEventType,
		FriendRecallEventType,
		FriendRenameEventType,
		FriendPokeEventType,
		GroupMemberPermissionUpdatedEventType,
		GroupNameUpdatedEventType,
		GroupMuteEventType,
		GroupRecallEventType,
		GroupMemberJoinRequestEventType,
		GroupMemberIncreaseEventType,
		GroupMemberDecreaseEventType,
		GroupDigestEventType,
		GroupReactionEventType,
		GroupMemberSpecialTitleUpdatedEventType,
		GroupInviteEventType,
		BotConnectedEventType,
		BotDisconnectedEventType,
	}
}
