package cryo

import lgrmessage "github.com/LagrangeDev/LagrangeGo/message"

// Event cryo的事件模型接口定义
type Event interface {
	GetUniEvent() *UniEvent     // 获取事件的基础信息
	GetEventType() EventType    // 获取事件类型
	GetEventId() string         // 获取事件ID
	GetEventTag() []string      // 获取事件标签列表
	GetClient() *LagrangeClient // 获取事件的Bot客户端
	Clone() Event               // 克隆，也就是深拷贝一个事件，这个主要是用来解决事件的传递顺序问题的
}

// UniEvent 是cryo的基础事件模型，所有的事件都由这个事件组合而成
type UniEvent struct {
	payload any // 事件的负载，可以用来携带点东西，一般用不着

	// 事件的基本信息
	EventType EventType `json:"event_type,omitzero,omitempty"` // 事件类型
	EventId   string    `json:"event_id,omitzero,omitempty"`   // 事件ID，是一个uuid
	EventTags []string  `json:"event_tags,omitzero,omitempty"` // 事件标签列表
	Time      uint32    `json:"time,omitzero,omitempty"`       // 事件发生的时间戳

	// 接收到事件的Bot客户端的基本信息
	botClient   *LagrangeClient // 指向接收到事件的Bot客户端
	BotId       string          `json:"bot_id,omitzero,omitempty"`       // 机器人ID
	BotNickname string          `json:"bot_nickname,omitzero,omitempty"` // 机器人昵称
	BotUin      uint32          `json:"bot_uin,omitzero,omitempty"`      // 机器人Uin
	BotUid      string          `json:"bot_uid,omitzero,omitempty"`      // 机器人Uid
	Platform    string          `json:"platform,omitzero,omitempty"`     // 机器人平台
}

// GetUniEvent 获取事件的基础信息
func (e *UniEvent) GetUniEvent() *UniEvent {
	return e
}

// GetEventType 获取事件类型
func (e *UniEvent) GetEventType() EventType {
	return e.EventType
}

// GetEventId 获取事件ID
func (e *UniEvent) GetEventId() string {
	return e.EventId
}

// GetEventTag 获取事件标签列表
func (e *UniEvent) GetEventTag() []string {
	return e.EventTags
}

// GetClient 获取事件的Bot客户端
func (e *UniEvent) GetClient() *LagrangeClient {
	return e.botClient
}

// Clone 克隆事件
func (e *UniEvent) Clone() Event {
	// 克隆事件
	return &UniEvent{
		EventType:   e.EventType,
		EventId:     e.EventId,
		EventTags:   e.EventTags,
		Time:        e.Time,
		botClient:   e.botClient,
		BotId:       e.BotId,
		BotNickname: e.BotNickname,
		BotUin:      e.BotUin,
		BotUid:      e.BotUid,
		Platform:    e.Platform,
	}
}

// MessageEvent 是消息事件的接口定义
type MessageEvent interface {
	Event
	GetReplyDetail() (replySeq uint32, senderUin uint32, time uint32, elements []lgrmessage.IMessageElement) // 获取回复用的消息详情
	GetUniMessageEvent() *UniMessageEvent                                                                    // 获取事件的基础信息
	GetMessage() *Message                                                                                    // 获取消息元素
	GetIMessageElements() []lgrmessage.IMessageElement                                                       // 获取消息元素的LagrangeGo格式
	GetMessageId() uint32                                                                                    // 获取消息ID
	Send(args ...interface{}) (ok bool, messageId uint32)                                                    // 发送消息
	Reply(args ...interface{}) (ok bool, messageId uint32)                                                   // 回复消息
}

// UniMessageEvent 是消息事件的基础模型，其他消息事件都由这个事件组合而成
type UniMessageEvent struct {
	UniEvent              // 事件的基础信息
	MessageId      uint32 // 消息ID
	SenderUin      uint32 // 消息发送者的Uin
	SenderUid      string // 消息发送者的Uid
	SenderNickname string // 消息发送者的昵称
	SenderCardname string // 消息发送者的备注名
	IsSenderFriend bool   // 消息发送者是否是好友
	GroupUin       uint32 // 群号，如果是私聊消息则为对应的好友Uin
	GroupName      string // 群名称，如果是私聊消息则为好友名称

	MessageElements Message // 消息元素
}

func (e *UniMessageEvent) GetReplyDetail() (uint32, uint32, uint32, []lgrmessage.IMessageElement) {
	return e.MessageId, e.SenderUin, e.Time, e.MessageElements.ToIMessageElements()
}

func (e *UniMessageEvent) GetUniMessageEvent() *UniMessageEvent {
	return e
}

func (e *UniMessageEvent) GetMessage() *Message {
	return &e.MessageElements
}

func (e *UniMessageEvent) GetIMessageElements() []lgrmessage.IMessageElement {
	return e.MessageElements.ToIMessageElements()
}

func (e *UniMessageEvent) GetMessageId() uint32 {
	return e.MessageId
}

func (e *UniMessageEvent) Clone() Event {
	// 克隆事件
	return &UniMessageEvent{
		UniEvent: UniEvent{
			EventType:   e.EventType,
			EventId:     e.EventId,
			EventTags:   e.EventTags,
			Time:        e.Time,
			botClient:   e.botClient,
			BotId:       e.BotId,
			BotNickname: e.BotNickname,
			BotUin:      e.BotUin,
			BotUid:      e.BotUid,
			Platform:    e.Platform,
		},
		MessageId:       e.MessageId,
		SenderUin:       e.SenderUin,
		SenderUid:       e.SenderUid,
		SenderNickname:  e.SenderNickname,
		SenderCardname:  e.SenderCardname,
		IsSenderFriend:  e.IsSenderFriend,
		GroupUin:        e.GroupUin,
		GroupName:       e.GroupName,
		MessageElements: e.MessageElements,
	}
}

func (e *UniMessageEvent) Send(args ...interface{}) (ok bool, messageId uint32) {
	// 发送消息
	return e.botClient.Send(e, args...)
}

func (e *UniMessageEvent) Reply(args ...interface{}) (ok bool, messageId uint32) {
	// 回复消息
	return e.botClient.Send(e, args...)
}

type (
	// PrivateMessageEvent 私聊消息事件
	PrivateMessageEvent struct {
		UniMessageEvent
		InternalId uint32 // 内部ID
		ClientSeq  uint32 // 客户端序列号
		TargetUin  uint32 // 目标Uin
	}
	// GroupMessageEvent 群消息事件
	GroupMessageEvent struct {
		UniMessageEvent
		InternalId uint32 // 内部ID
	}
	// TempMessageEvent 临时消息事件
	TempMessageEvent struct {
		UniMessageEvent
	}
	// NewFriendRequestEvent 新好友请求事件
	NewFriendRequestEvent struct {
		UniEvent
		Uin      uint32
		Uid      string
		Nickname string
		Message  string
		From     string
	}
	// NewFriendEvent 新好友事件
	NewFriendEvent struct {
		UniEvent
		Uin      uint32
		Uid      string
		Nickname string
		Message  string
	}
	// FriendRecallEvent 好友撤回事件
	FriendRecallEvent struct {
		UniEvent
		Uin     uint32
		Uid     string
		Seqence uint64
		Random  uint32
	}
	// FriendRenameEvent 好友改名事件
	FriendRenameEvent struct {
		UniEvent
		IsSelf   bool
		Uin      uint32
		Uid      string
		Nickname string
	}
	// FriendPokeEvent 好友戳一戳事件
	FriendPokeEvent struct {
		UniEvent
		SenderUin uint32
		TargetUin uint32
		Suffix    string
		Action    string
	}
	// GroupMemberPermissionUpdatedEvent 群成员权限变更事件
	GroupMemberPermissionUpdatedEvent struct {
		UniEvent
		GroupUin uint32
		Uin      uint32
		Uid      string
		IsAdmin  bool
	}
	// GroupNameUpdatedEvent 群名称变更事件
	GroupNameUpdatedEvent struct {
		UniEvent
		GroupUin uint32
		NewName  string
	}
	// GroupMuteEvent 群禁言事件
	GroupMuteEvent struct {
		UniEvent
		GroupUin    uint32
		OperatorUin uint32
		OperatorUid string
		TargetUin   uint32
		TargetUid   string
		Duration    uint32
		isMuteAll   bool
	}
	// GroupRecallEvent 群撤回事件
	GroupRecallEvent struct {
		UniEvent
		GroupUin    uint32
		OperatorUin uint32
		OperatorUid string
		SenderUin   uint32
		SenderUid   string
		Seqence     uint64
		Random      uint32
	}
	// GroupMemberJoinRequestEvent 群成员入群请求事件
	GroupMemberJoinRequestEvent struct {
		UniEvent
		GroupUin       uint32
		SenderUin      uint32
		SenderUid      string
		SenderNickname string
		InviterUin     uint32
		InviterUid     string
		Answer         string
		RequestSeqence uint64
	}
	// GroupMemberIncreaseEvent 群成员增加事件
	GroupMemberIncreaseEvent struct {
		UniEvent
		GroupUin   uint32
		Uin        uint32
		Uid        string
		InviterUin uint32
		InviterUid string
		IsSelf     bool
	}
	// GroupMemberDecreaseEvent 群成员减少事件
	GroupMemberDecreaseEvent struct {
		UniEvent
		GroupUin uint32
		Uin      uint32
		Uid      string
		IsSelf   bool
		IsKicked bool
	}
	// GroupDigestEvent 群精华消息事件
	GroupDigestEvent struct {
		UniEvent
		GroupUin         uint32
		MessageId        uint32
		InternalId       uint32
		SenderUin        uint32
		SenderUid        string
		SenderNickname   string
		OperatorUin      uint32
		OperatorNickname string
		IsRemove         bool
	}
	// GroupReactionEvent 群消息表态事件
	GroupReactionEvent struct {
		UniEvent
		GroupUin  uint32
		Uin       uint32
		Uid       string
		TargetSeq uint32
		IsAdd     bool
		IsEmoji   bool
		Code      string
		Count     uint32
	}
	// GroupMemberSpecialTitleUpdated 群成员特殊头衔变更事件
	GroupMemberSpecialTitleUpdated struct {
		UniEvent
		GroupUin uint32
		Uin      uint32
		Uid      string
		NewTitle string
	}
	// GroupInviteEvent 加群邀请事件
	GroupInviteEvent struct {
		UniEvent
		GroupUin        uint32
		GroupName       string
		InviterUin      uint32
		InviterUid      string
		InviterNickname string
		RequestSeqence  uint64
	}
	// BotConnectedEvent 机器人连接事件
	BotConnectedEvent struct {
		UniEvent
		Version string
	}
	// BotDisconnectedEvent 机器人断开连接事件
	BotDisconnectedEvent struct {
		UniEvent
	}
	CustomEvent struct {
		UniEvent
		summury string      // 摘要
		payload interface{} // 负载
	}
)

func (e *PrivateMessageEvent) Clone() Event {
	// 克隆事件
	return &PrivateMessageEvent{
		UniMessageEvent: UniMessageEvent{
			UniEvent: UniEvent{
				EventType:   e.EventType,
				EventId:     e.EventId,
				EventTags:   e.EventTags,
				Time:        e.Time,
				botClient:   e.botClient,
				BotId:       e.BotId,
				BotNickname: e.BotNickname,
				BotUin:      e.BotUin,
				BotUid:      e.BotUid,
				Platform:    e.Platform,
			},
			MessageId:       e.MessageId,
			SenderUin:       e.SenderUin,
			SenderUid:       e.SenderUid,
			SenderNickname:  e.SenderNickname,
			SenderCardname:  e.SenderCardname,
			IsSenderFriend:  e.IsSenderFriend,
			GroupUin:        e.GroupUin,
			GroupName:       e.GroupName,
			MessageElements: e.MessageElements,
		},
		InternalId: e.InternalId,
		ClientSeq:  e.ClientSeq,
		TargetUin:  e.TargetUin,
	}
}

func (e *GroupMessageEvent) Clone() Event {
	// 克隆事件
	return &GroupMessageEvent{
		UniMessageEvent: UniMessageEvent{
			UniEvent: UniEvent{
				EventType:   e.EventType,
				EventId:     e.EventId,
				EventTags:   e.EventTags,
				Time:        e.Time,
				botClient:   e.botClient,
				BotId:       e.BotId,
				BotNickname: e.BotNickname,
				BotUin:      e.BotUin,
				BotUid:      e.BotUid,
				Platform:    e.Platform,
			},
			MessageId:       e.MessageId,
			SenderUin:       e.SenderUin,
			SenderUid:       e.SenderUid,
			SenderNickname:  e.SenderNickname,
			SenderCardname:  e.SenderCardname,
			IsSenderFriend:  e.IsSenderFriend,
			GroupUin:        e.GroupUin,
			GroupName:       e.GroupName,
			MessageElements: e.MessageElements,
		},
		InternalId: e.InternalId,
	}
}

func (e *TempMessageEvent) Clone() Event {
	// 克隆事件
	return &TempMessageEvent{
		UniMessageEvent: UniMessageEvent{
			UniEvent: UniEvent{
				EventType:   e.EventType,
				EventId:     e.EventId,
				EventTags:   e.EventTags,
				Time:        e.Time,
				botClient:   e.botClient,
				BotId:       e.BotId,
				BotNickname: e.BotNickname,
				BotUin:      e.BotUin,
				BotUid:      e.BotUid,
				Platform:    e.Platform,
			},
			MessageId:       e.MessageId,
			SenderUin:       e.SenderUin,
			SenderUid:       e.SenderUid,
			SenderNickname:  e.SenderNickname,
			SenderCardname:  e.SenderCardname,
			IsSenderFriend:  e.IsSenderFriend,
			GroupUin:        e.GroupUin,
			GroupName:       e.GroupName,
			MessageElements: e.MessageElements,
		},
	}
}

func (e *NewFriendRequestEvent) Clone() Event {
	// 克隆事件
	return &NewFriendRequestEvent{
		UniEvent: UniEvent{
			EventType:   e.EventType,
			EventId:     e.EventId,
			EventTags:   e.EventTags,
			Time:        e.Time,
			botClient:   e.botClient,
			BotId:       e.BotId,
			BotNickname: e.BotNickname,
			BotUin:      e.BotUin,
			BotUid:      e.BotUid,
			Platform:    e.Platform,
		},
		Uin:      e.Uin,
		Uid:      e.Uid,
		Nickname: e.Nickname,
		Message:  e.Message,
	}
}

func (e *NewFriendEvent) Clone() Event {
	// 克隆事件
	return &NewFriendEvent{
		UniEvent: UniEvent{
			EventType:   e.EventType,
			EventId:     e.EventId,
			EventTags:   e.EventTags,
			Time:        e.Time,
			botClient:   e.botClient,
			BotId:       e.BotId,
			BotNickname: e.BotNickname,
			BotUin:      e.BotUin,
			BotUid:      e.BotUid,
			Platform:    e.Platform,
		},
		Uin:      e.Uin,
		Uid:      e.Uid,
		Nickname: e.Nickname,
		Message:  e.Message,
	}
}

func (e *FriendRecallEvent) Clone() Event {
	// 克隆事件
	return &FriendRecallEvent{
		UniEvent: UniEvent{
			EventType:   e.EventType,
			EventId:     e.EventId,
			EventTags:   e.EventTags,
			Time:        e.Time,
			botClient:   e.botClient,
			BotId:       e.BotId,
			BotNickname: e.BotNickname,
			BotUin:      e.BotUin,
			BotUid:      e.BotUid,
			Platform:    e.Platform,
		},
		Uin:     e.Uin,
		Uid:     e.Uid,
		Seqence: e.Seqence,
		Random:  e.Random,
	}
}

func (e *FriendRenameEvent) Clone() Event {
	// 克隆事件
	return &FriendRenameEvent{
		UniEvent: UniEvent{
			EventType:   e.EventType,
			EventId:     e.EventId,
			EventTags:   e.EventTags,
			Time:        e.Time,
			botClient:   e.botClient,
			BotId:       e.BotId,
			BotNickname: e.BotNickname,
			BotUin:      e.BotUin,
			BotUid:      e.BotUid,
			Platform:    e.Platform,
		},
		IsSelf:   e.IsSelf,
		Uin:      e.Uin,
		Uid:      e.Uid,
		Nickname: e.Nickname,
	}
}

func (e *FriendPokeEvent) Clone() Event {
	// 克隆事件
	return &FriendPokeEvent{
		UniEvent: UniEvent{
			EventType:   e.EventType,
			EventId:     e.EventId,
			EventTags:   e.EventTags,
			Time:        e.Time,
			botClient:   e.botClient,
			BotId:       e.BotId,
			BotNickname: e.BotNickname,
			BotUin:      e.BotUin,
			BotUid:      e.BotUid,
			Platform:    e.Platform,
		},
		SenderUin: e.SenderUin,
		TargetUin: e.TargetUin,
		Suffix:    e.Suffix,
		Action:    e.Action,
	}
}

func (e *GroupMemberPermissionUpdatedEvent) Clone() Event {
	// 克隆事件
	return &GroupMemberPermissionUpdatedEvent{
		UniEvent: UniEvent{
			EventType:   e.EventType,
			EventId:     e.EventId,
			EventTags:   e.EventTags,
			Time:        e.Time,
			botClient:   e.botClient,
			BotId:       e.BotId,
			BotNickname: e.BotNickname,
			BotUin:      e.BotUin,
			BotUid:      e.BotUid,
			Platform:    e.Platform,
		},
		GroupUin: e.GroupUin,
		Uin:      e.Uin,
		Uid:      e.Uid,
		IsAdmin:  e.IsAdmin,
	}
}
func (e *GroupNameUpdatedEvent) Clone() Event {
	// 克隆事件
	return &GroupNameUpdatedEvent{
		UniEvent: UniEvent{
			EventType:   e.EventType,
			EventId:     e.EventId,
			EventTags:   e.EventTags,
			Time:        e.Time,
			botClient:   e.botClient,
			BotId:       e.BotId,
			BotNickname: e.BotNickname,
			BotUin:      e.BotUin,
			BotUid:      e.BotUid,
			Platform:    e.Platform,
		},
		GroupUin: e.GroupUin,
		NewName:  e.NewName,
	}
}
func (e *GroupMuteEvent) Clone() Event {
	// 克隆事件
	return &GroupMuteEvent{
		UniEvent: UniEvent{
			EventType:   e.EventType,
			EventId:     e.EventId,
			EventTags:   e.EventTags,
			Time:        e.Time,
			botClient:   e.botClient,
			BotId:       e.BotId,
			BotNickname: e.BotNickname,
			BotUin:      e.BotUin,
			BotUid:      e.BotUid,
			Platform:    e.Platform,
		},
		GroupUin:    e.GroupUin,
		OperatorUin: e.OperatorUin,
		OperatorUid: e.OperatorUid,
		TargetUin:   e.TargetUin,
		TargetUid:   e.TargetUid,
		Duration:    e.Duration,
		isMuteAll:   e.isMuteAll,
	}
}
func (e *GroupRecallEvent) Clone() Event {
	// 克隆事件
	return &GroupRecallEvent{
		UniEvent: UniEvent{
			EventType:   e.EventType,
			EventId:     e.EventId,
			EventTags:   e.EventTags,
			Time:        e.Time,
			botClient:   e.botClient,
			BotId:       e.BotId,
			BotNickname: e.BotNickname,
			BotUin:      e.BotUin,
			BotUid:      e.BotUid,
			Platform:    e.Platform,
		},
		GroupUin:    e.GroupUin,
		OperatorUin: e.OperatorUin,
		OperatorUid: e.OperatorUid,
		SenderUin:   e.SenderUin,
		SenderUid:   e.SenderUid,
		Seqence:     e.Seqence,
		Random:      e.Random,
	}
}
func (e *GroupMemberJoinRequestEvent) Clone() Event {
	// 克隆事件
	return &GroupMemberJoinRequestEvent{
		UniEvent: UniEvent{
			EventType:   e.EventType,
			EventId:     e.EventId,
			EventTags:   e.EventTags,
			Time:        e.Time,
			botClient:   e.botClient,
			BotId:       e.BotId,
			BotNickname: e.BotNickname,
			BotUin:      e.BotUin,
			BotUid:      e.BotUid,
			Platform:    e.Platform,
		},
		GroupUin:       e.GroupUin,
		SenderUin:      e.SenderUin,
		SenderUid:      e.SenderUid,
		SenderNickname: e.SenderNickname,
		InviterUin:     e.InviterUin,
		InviterUid:     e.InviterUid,
		RequestSeqence: e.RequestSeqence,
	}
}
func (e *GroupMemberIncreaseEvent) Clone() Event {
	// 克隆事件
	return &GroupMemberIncreaseEvent{
		UniEvent: UniEvent{
			EventType:   e.EventType,
			EventId:     e.EventId,
			EventTags:   e.EventTags,
			Time:        e.Time,
			botClient:   e.botClient,
			BotId:       e.BotId,
			BotNickname: e.BotNickname,
			BotUin:      e.BotUin,
			BotUid:      e.BotUid,
			Platform:    e.Platform,
		},
		GroupUin:   e.GroupUin,
		Uin:        e.Uin,
		Uid:        e.Uid,
		InviterUin: e.InviterUin,
		InviterUid: e.InviterUid,
		IsSelf:     e.IsSelf,
	}
}
func (e *GroupMemberDecreaseEvent) Clone() Event {
	// 克隆事件
	return &GroupMemberDecreaseEvent{
		UniEvent: UniEvent{
			EventType:   e.EventType,
			EventId:     e.EventId,
			EventTags:   e.EventTags,
			Time:        e.Time,
			botClient:   e.botClient,
			BotId:       e.BotId,
			BotNickname: e.BotNickname,
			BotUin:      e.BotUin,
			BotUid:      e.BotUid,
			Platform:    e.Platform,
		},
		GroupUin: e.GroupUin,
		Uin:      e.Uin,
		Uid:      e.Uid,
		IsSelf:   e.IsSelf,
	}
}
func (e *GroupDigestEvent) Clone() Event {
	// 克隆事件
	return &GroupDigestEvent{
		UniEvent: UniEvent{
			EventType:   e.EventType,
			EventId:     e.EventId,
			EventTags:   e.EventTags,
			Time:        e.Time,
			botClient:   e.botClient,
			BotId:       e.BotId,
			BotNickname: e.BotNickname,
			BotUin:      e.BotUin,
			BotUid:      e.BotUid,
			Platform:    e.Platform,
		},
		GroupUin:         e.GroupUin,
		MessageId:        e.MessageId,
		InternalId:       e.InternalId,
		SenderUin:        e.SenderUin,
		SenderUid:        e.SenderUid,
		SenderNickname:   e.SenderNickname,
		OperatorUin:      e.OperatorUin,
		OperatorNickname: e.OperatorNickname,
		IsRemove:         e.IsRemove,
	}
}
func (e *GroupReactionEvent) Clone() Event {
	// 克隆事件
	return &GroupReactionEvent{
		UniEvent: UniEvent{
			EventType:   e.EventType,
			EventId:     e.EventId,
			EventTags:   e.EventTags,
			Time:        e.Time,
			botClient:   e.botClient,
			BotId:       e.BotId,
			BotNickname: e.BotNickname,
			BotUin:      e.BotUin,
			BotUid:      e.BotUid,
			Platform:    e.Platform,
		},
		GroupUin:  e.GroupUin,
		Uin:       e.Uin,
		Uid:       e.Uid,
		TargetSeq: e.TargetSeq,
		IsAdd:     e.IsAdd,
		IsEmoji:   e.IsEmoji,
		Code:      e.Code,
		Count:     e.Count,
	}
}

func (e *GroupMemberSpecialTitleUpdated) Clone() Event {
	// 克隆事件
	return &GroupMemberSpecialTitleUpdated{
		UniEvent: UniEvent{
			EventType:   e.EventType,
			EventId:     e.EventId,
			EventTags:   e.EventTags,
			Time:        e.Time,
			botClient:   e.botClient,
			BotId:       e.BotId,
			BotNickname: e.BotNickname,
			BotUin:      e.BotUin,
			BotUid:      e.BotUid,
			Platform:    e.Platform,
		},
		GroupUin: e.GroupUin,
		Uin:      e.Uin,
		Uid:      e.Uid,
		NewTitle: e.NewTitle,
	}
}

func (e *GroupInviteEvent) Clone() Event {
	// 克隆事件
	return &GroupInviteEvent{
		UniEvent: UniEvent{
			EventType:   e.EventType,
			EventId:     e.EventId,
			EventTags:   e.EventTags,
			Time:        e.Time,
			botClient:   e.botClient,
			BotId:       e.BotId,
			BotNickname: e.BotNickname,
			BotUin:      e.BotUin,
			BotUid:      e.BotUid,
			Platform:    e.Platform,
		},
		GroupUin:        e.GroupUin,
		GroupName:       e.GroupName,
		InviterUin:      e.InviterUin,
		InviterUid:      e.InviterUid,
		InviterNickname: e.InviterNickname,
		RequestSeqence:  e.RequestSeqence,
	}
}

func (e *BotConnectedEvent) Clone() Event {
	// 克隆事件
	return &BotConnectedEvent{
		UniEvent: UniEvent{
			EventType:   e.EventType,
			EventId:     e.EventId,
			EventTags:   e.EventTags,
			Time:        e.Time,
			botClient:   e.botClient,
			BotId:       e.BotId,
			BotNickname: e.BotNickname,
			BotUin:      e.BotUin,
			BotUid:      e.BotUid,
			Platform:    e.Platform,
		},
		Version: e.Version,
	}
}

func (e *BotDisconnectedEvent) Clone() Event {
	// 克隆事件
	return &BotDisconnectedEvent{
		UniEvent: UniEvent{
			EventType:   e.EventType,
			EventId:     e.EventId,
			EventTags:   e.EventTags,
			Time:        e.Time,
			botClient:   e.botClient,
			BotId:       e.BotId,
			BotNickname: e.BotNickname,
			BotUin:      e.BotUin,
			BotUid:      e.BotUid,
			Platform:    e.Platform,
		},
	}
}

func (e *CustomEvent) Clone() Event {
	// 克隆事件
	return &CustomEvent{
		UniEvent: UniEvent{
			EventType:   e.EventType,
			EventId:     e.EventId,
			EventTags:   e.EventTags,
			Time:        e.Time,
			botClient:   e.botClient,
			BotId:       e.BotId,
			BotNickname: e.BotNickname,
			BotUin:      e.BotUin,
			BotUid:      e.BotUid,
			Platform:    e.Platform,
		},
		summury: e.summury,
		payload: e.payload,
	}
}
